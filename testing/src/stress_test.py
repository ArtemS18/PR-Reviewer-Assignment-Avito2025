import asyncio
import json
import random
import time
import uuid
from typing import List

import httpx

import config


def load_seed_data() -> list[dict]:
    with open(config.SEED_DATA_FILE, "r", encoding="utf-8") as f:
        data: dict = json.load(f)
    users = data.get("users", [])
    if not users:
        raise RuntimeError("No users in seed_data.json, run seed_db.py first")
    return users


async def request_create_pr(client: httpx.AsyncClient, users: list[dict]) -> httpx.Response:
    author = random.choice(users)
    pr_id = f"pr-{int(time.time() * 1e6)}"
    return await client.post(
        f"{config.BASE_URL}/pullRequest/create",
        json={
            "pull_request_id": pr_id,
            "pull_request_name": "Load Test PR",
            "author_id": author["id"],
        },
        timeout=0.3,
    )


async def request_set_user_active(client: httpx.AsyncClient, users: list[dict]) -> httpx.Response:
    user = random.choice(users)
    is_active = random.choice([True, False])
    return await client.post(
        f"{config.BASE_URL}/users/setIsActive",
        json={
            "user_id": user["id"],
            "is_active": is_active,
        },
        timeout=0.3,
    )


async def request_get_user_review(client: httpx.AsyncClient, users: list[dict]) -> httpx.Response:
    user = random.choice(users)
    return await client.get(
        f"{config.BASE_URL}/users/getReview",
        params={"user_id": user["id"]},
        timeout=0.3,
    )


async def request_get_team(client: httpx.AsyncClient, users: list[dict]) -> httpx.Response:
    user = random.choice(users)
    team_name = user.get("team") or "non-existing-team"
    return await client.get(
        f"{config.BASE_URL}/team/get",
        params={"team_name": team_name},
        timeout=0.3,
    )


async def request_add_team(client: httpx.AsyncClient) -> httpx.Response:
    team_name = f"load-test-team-{uuid.uuid4()}"
    members = [
        {
            "user_id": f"user-{uuid.uuid4()[:10]}",
            "username": "lt-user",
            "is_active": True,
        }
    ]
    return await client.post(
        f"{config.BASE_URL}/team/add",
        json={"team_name": team_name, "members": members},
        timeout=0.3,
    )

async def one_request(client: httpx.AsyncClient, users: list[dict]) -> httpx.Response:
    op = random.choice(
        [
            "create_pr",
            "set_user_active",
            "get_user_review",
            "get_team",
        ]
    )

    if op == "create_pr":
        return await request_create_pr(client, users)
    if op == "set_user_active":
        return await request_set_user_active(client, users)
    if op == "get_user_review":
        return await request_get_user_review(client, users)
    if op == "get_team":
        return await request_get_team(client, users)

    return await request_create_pr(client, users)


async def worker(
    id_: int,
    client: httpx.AsyncClient,
    interval: float,
    end_time: float,
    users: list[dict],
    latencies: list[float],
    successes: list[int],
    errors: list,
    total_counter: list[int],
):
    while time.time() < end_time:
        start = time.time()
        try:
            resp = await one_request(client, users)
            latency_ms = (time.time() - start) * 1000
            latencies.append(latency_ms)
            total_counter[0] += 1
            if 200 <= resp.status_code < 300:
                successes.append(1)
            else:
                errors.append((resp.status_code, latency_ms))
        except Exception as e:
            errors.append((str(e), None))
            total_counter[0] += 1
        sleep_for = interval - (time.time() - start)
        if sleep_for > 0:
            await asyncio.sleep(sleep_for)


async def spinner_task(total_counter: list[int], errors: list, end_time: float):
    spinner_frames = ["-", "\\", "|", "/"]
    idx = 0
    while time.time() < end_time:
        frame = spinner_frames[idx % len(spinner_frames)]
        idx += 1
        total = total_counter[0]
        err_count = len(errors)
        print(
            f"\r{frame} Running stress test... total={total} errors={err_count}",
            end="",
            flush=True,
        )
        await asyncio.sleep(0.1)
    total = total_counter[0]
    err_count = len(errors)
    print(f"\r✓ Stress test in progress... total={total} errors={err_count}")


async def run_stress_test():
    users = load_seed_data()

    duration = config.DURATION_SECONDS
    end_time = time.time() + duration

    total_rps = config.RPS
    req_per_worker_per_sec = total_rps / config.CONCURRENCY
    interval = 1.0 / req_per_worker_per_sec

    latencies: List[float] = []
    successes: List[int] = []
    errors: List = []
    total_counter = [0]

    print("\n--- Start stress test ---")

    async with httpx.AsyncClient() as client:
        worker_tasks = [
            asyncio.create_task(
                worker(
                    i,
                    client,
                    interval,
                    end_time,
                    users,
                    latencies,
                    successes,
                    errors,
                    total_counter,
                )
            )
            for i in range(config.CONCURRENCY)
        ]
        spin_task = asyncio.create_task(
            spinner_task(total_counter, errors, end_time)
        )

        await asyncio.gather(*worker_tasks)
        await spin_task

    total = len(latencies)
    success_count = len(successes)
    success_rate = (success_count / total) * 100 if total > 0 else 0

    sorted_lat = sorted(latencies)
    p50 = sorted_lat[int(0.50 * total)] if total else 0
    p95 = sorted_lat[int(0.95 * total) - 1] if total else 0
    p99 = sorted_lat[int(0.99 * total) - 1] if total else 0

    print("\n\n--- Stress test results ---")
    print(f"Total requests: {total}")
    print(f"Success: {success_count} ({success_rate:.4f}%)")
    print(f"Latency p50: {p50:.2f} ms")
    print(f"Latency p95: {p95:.2f} ms")
    print(f"Latency p99: {p99:.2f} ms")
    print(f"Errors: {len(errors)}")

    print("\nSLI check:")
    print(f"- success_rate >= 99.9%: {'OK' if success_rate >= 99.9 else 'FAIL'}")
    print(f"- p95 <= 300 ms: {'OK' if p95 <= 300 else 'FAIL'}")

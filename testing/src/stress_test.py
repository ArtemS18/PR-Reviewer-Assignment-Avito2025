import asyncio
import json
import random
import time
from typing import List
import config

import httpx



def load_seed_data() -> dict:
    with open(config.SEED_DATA_FILE, "r", encoding="utf-8") as f:
        data: dict = json.load(f)
    users = data.get("users", [])
    if not users:
        raise RuntimeError("No users in seed_data.json, run seed_db.py first")
    return users


async def one_request(client: httpx.AsyncClient, users: dict):
    author = random.choice(users)
    pr_id = f"pr-{int(time.time() * 1e6)}"

    resp = await client.post(
        f"{config.BASE_URL}/pullRequest/create",
        json={
            "pull_request_id": pr_id,
            "pull_request_name": "Load Test PR",
            "author_id": author["id"],
        },
        timeout=0.3,
    )
    return resp


async def worker(id_, client: httpx.AsyncClient, interval: float, end_time: float, users: dict, latencies: list[float], successes: list[int], errors: list):
    while time.time() < end_time:
        start = time.time()
        try:
            resp = await one_request(client, users)
            latency_ms = (time.time() - start) * 1000
            latencies.append(latency_ms)
            if 200 <= resp.status_code < 300:
                successes.append(1)
            else:
                errors.append((resp.status_code, latency_ms))
        except Exception as e:
            errors.append((str(e), None))
        sleep_for = interval - (time.time() - start)
        if sleep_for > 0:
            await asyncio.sleep(sleep_for)

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

    async with httpx.AsyncClient() as client:
        tasks = [
            asyncio.create_task(
                worker(i, client, interval, end_time, users, latencies, successes, errors)
            )
            for i in range(config.CONCURRENCY)
        ]
        await asyncio.gather(*tasks)

    total = len(latencies)
    success_count = len(successes)
    success_rate = (success_count / total) * 100 if total > 0 else 0

    sorted_lat = sorted(latencies)
    p50 = sorted_lat[int(0.50 * total)] if total else 0
    p95 = sorted_lat[int(0.95 * total) - 1] if total else 0
    p99 = sorted_lat[int(0.99 * total) - 1] if total else 0

    print("\n--- Stress test results ---")
    print(f"Total requests: {total}")
    print(f"Success: {success_count} ({success_rate:.4f}%)")
    print(f"Latency p50: {p50:.2f} ms")
    print(f"Latency p95: {p95:.2f} ms")
    print(f"Latency p99: {p99:.2f} ms")
    print(f"Errors: {len(errors)}")

    print("\nSLI check:")
    print(f"- success_rate >= 99.9%: {'OK' if success_rate >= 99.9 else 'FAIL'}")
    print(f"- p95 <= 300 ms: {'OK' if p95 <= 300 else 'FAIL'}")
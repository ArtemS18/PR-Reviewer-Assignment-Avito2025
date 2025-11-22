import asyncio
import json
import random
import config
import string

import httpx


def rand_id(prefix: str, length: int = 8) -> str:
    return f"{prefix}-" + "".join(
        random.choices(string.ascii_lowercase + string.digits, k=length)
    )


async def create_team(client: httpx.AsyncClient, idx: int):
    team_name = f"team-{idx}"
    members = []
    for i in range(config.USERS_PER_TEAM):
        user_id = rand_id(f"user{idx}_{i}")
        members.append(
            {
                "user_id": user_id,
                "username": f"user-{idx}-{i}",
                "is_active": True,
            }
        )

    resp = await client.post(
        f"{config.BASE_URL}/team/add",
        json={"team_name": team_name, "members": members},
        timeout=5.0,
    )
    resp.raise_for_status()
    return {
        "team_name": team_name,
        "members": members,
    }


async def fill_db():
    async with httpx.AsyncClient() as client:
        tasks = [create_team(client, i) for i in range(config.TEAMS_COUNT)]
        teams = await asyncio.gather(*tasks)

    users = []
    for t in teams:
        for m in t["members"]:
            users.append(
                {
                    "id": m["user_id"],
                    "name": m["username"],
                    "team": t["team_name"],
                }
            )

    data = {"teams": teams, "users": users}

    with open(config.SEED_DATA_FILE, "w", encoding="utf-8") as f:
        json.dump(data, f, ensure_ascii=False, indent=2)

    print(f"Seeded {len(teams)} teams, {len(users)} users")
    print(f"Saved data to {config.SEED_DATA_FILE}")


import json
import random
import sys
import time
from pathlib import Path

from locust import HttpUser, task, between, main

import config


def load_seed_data() -> list[dict]:
    seed_path = Path(config.SEED_DATA_FILE)
    with seed_path.open("r", encoding="utf-8") as f:
        data: dict = json.load(f)
    users = data.get("users", [])
    if not users:
        raise RuntimeError("No users in seed_data.json, run seed_db.py first")
    return users


class PRServiceUser(HttpUser):
    """
    - POST /pullRequest/create
    - POST /users/setIsActive
    - GET  /users/getReview
    - GET  /team/get
    """

    wait_time = between(0.1, 0.3)
    USERS: list[dict] = []
    def on_start(self):
        if not PRServiceUser.USERS:
            PRServiceUser.USERS = load_seed_data()
    @task(4)
    def create_pull_request(self):
        author = random.choice(self.USERS)
        pr_id = f"pr-{int(time.time() * 1e6)}"

        self.client.post(
            "/pullRequest/create",
            json={
                "pull_request_id": pr_id,
                "pull_request_name": "Load Test PR",
                "author_id": author["id"],
            },
            name="pullRequest/create",
            timeout=0.3,
        )

    @task(2)
    def set_user_active(self):
        user = random.choice(self.USERS)
        is_active = random.choice([True, False])

        self.client.post(
            "/users/setIsActive",
            json={
                "user_id": user["id"],
                "is_active": is_active,
            },
            name="users/setIsActive",
            timeout=0.3,
        )

    @task(2)
    def get_user_review(self):
        user = random.choice(self.USERS)
        self.client.get(
            "/users/getReview",
            params={"user_id": user["id"]},
            name="users/getReview",
            timeout=0.3,
        )

    @task(1)
    def get_team(self):
        user = random.choice(self.USERS)
        team_name = user.get("team") or "non-existing-team"

        self.client.get(
            "/team/get",
            params={"team_name": team_name},
            name="team/get",
            timeout=0.3,
        )


def run_stress_test():
    PRServiceUser.USERS = load_seed_data()
    sys.argv = [
        "locust",
        "-f", "src/stress_test.py",
        f"--host={config.BASE_URL}",
        "--headless", "-u", "10", "-r", "2", "-t", "1m",
    ]
    main.main()

if __name__ == "__main__":
    run_stress_test()
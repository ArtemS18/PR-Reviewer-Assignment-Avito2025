import json
import random
import string
import uuid

import psycopg2

import config


def rand_suffix(length: int = 8) -> str:
    return "".join(
        random.choices(string.ascii_lowercase + string.digits, k=length)
    )
def truncate_db():
    conn = psycopg2.connect(
        host=config.DB_HOST,
        port=config.DB_PORT,
        dbname=config.DB_NAME,
        user=config.DB_USER,
        password=config.DB_PASS,
    )
    conn.autocommit = False
    cur = conn.cursor()
    cur.execute("TRUNCATE reviewers, pull_requests, users, teams RESTART IDENTITY CASCADE;")
    conn.commit()
    cur.close()
    conn.close()

def fill_db():
    conn = psycopg2.connect(
        host=config.DB_HOST,
        port=config.DB_PORT,
        dbname=config.DB_NAME,
        user=config.DB_USER,
        password=config.DB_PASS,
    )
    conn.autocommit = False
    cur = conn.cursor()

    cur.execute("TRUNCATE reviewers, pull_requests, users, teams RESTART IDENTITY CASCADE;")

    teams = []
    users = []

    for ti in range(config.TEAMS_COUNT):
        team_id = f"team-{uuid.uuid4()}"
        team_name = f"team-{ti}"

        cur.execute(
            """
            INSERT INTO teams (id, name)
            VALUES (%s, %s)
            ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
            RETURNING id
            """,
            (team_id, team_name),
        )
        row = cur.fetchone()
        if row:
            team_id = row[0]

        team_record = {"team_name": team_name, "members": []}

        for ui in range(config.USERS_PER_TEAM):
            user_id = f"user-{rand_suffix()}"
            username = f"user-{ti}-{ui}"
            cur.execute(
                """
                INSERT INTO users (id, name, is_active, team_id)
                VALUES (%s, %s, %s, %s)
                ON CONFLICT (id) DO NOTHING
                """,
                (user_id, username, True, team_id),
            )

            user_info = {
                "id": user_id,
                "name": username,
                "team": team_name,
            }
            users.append(user_info)
            team_record["members"].append(
                {
                    "user_id": user_id,
                    "username": username,
                    "is_active": True,
                }
            )

        teams.append(team_record)

    conn.commit()
    cur.close()
    conn.close()

    data = {"teams": teams, "users": users}

    with open(config.SEED_DATA_FILE, "w", encoding="utf-8") as f:
        json.dump(data, f, ensure_ascii=False, indent=2)

    print(f"Seeded {len(teams)} teams, {len(users)} users")
    print(f"Saved data to {config.SEED_DATA_FILE}")


if __name__ == "__main__":
    fill_db()

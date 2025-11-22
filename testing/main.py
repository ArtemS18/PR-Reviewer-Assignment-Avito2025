import asyncio
from src import seed_db, stress_test
import sys

async def main():
    if sys.argv[1] == "test":
        await stress_test.run_stress_test()
    elif sys.argv[1] == "fill":
        await seed_db.fill_db()
    else:
        raise RuntimeError("unexpect command")


if __name__ == "__main__":
    asyncio.run(main())

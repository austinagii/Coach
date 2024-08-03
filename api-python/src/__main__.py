import json
import pymongo
from pymongo import MongoClient
from openai import OpenAI
from .user import User
from .goal import Goal, Milestone
from .schedule import ScheduledActivity, DailySchedule

def main():
    user = User(
        name="Joe Shmoe",
        summary="Joe Shmoe is a budding runner who wants to buy a home",
        goals=[
            Goal(
                id=1,
                title="Run a marathon", 
                description="Run a marathon in less than a year",
                milestones = [
                    Milestone(
                        title="Buy running shoes", 
                        description="Save enough money to buy a pair of running shoes",
                        targetDate="2024-10-11"
                    ),
                    Milestone(
                        title="Run 1 mile",
                        description="Build up towards running a mile in 10 minutes",
                        targetDate="2025-03-30"
                    ),
                    Milestone(
                        title="Run the marathon",
                        description="Run the Boston marathon",
                        targetDate="2025-10-30"
                    )
                ]
            )
        ],
        schedule=DailySchedule(
            routine=[
                ScheduledActivity(
                    title="Stretch",
                    description="Loosen up your muscles in preparation to run",
                    startTime="06:00",
                    endTime="06:15"
                ),
                ScheduledActivity(
                    title="Complete Interval Training",
                    description="Run intervals of 2 mins running 1 min walking",
                    startTime="06:15",
                    endTime="06:45"
                ),
                ScheduledActivity(
                    title="Meal Prep",
                    description="Prepare meals for the day to minimize spending",
                    startTime="07:30",
                    endTime="09:00"
                ),
                ScheduledActivity(
                    title="Work A Job",
                    description="Go to work to earn money to buy your home",
                    startTime="10:00",
                    endTime="18:00"
                ),
                ScheduledActivity(
                    title="Go To Bed",
                    description="Go to bed and get some rest ahead of tomorrows run",
                    startTime="21:00",
                    endTime="05:00"
                ),
            ]
        )
    )

    try:
        db_uri = "mongodb://host.docker.internal:27017"
        client = MongoClient(db_uri)
        database = client["aisu"]
        collection = database["users"]
        
        user_dict = user.dict()
        result = collection.insert_one(user_dict)
        print(f"User inserted with id: {result.inserted_id}")

        client.close()
    except Exception as e:
        raise Exception("The following error occurred", e)
        
    
def promptModel():
    client = OpenAI()
    completion = client.chat.completions.create(
        model="gpt-4o-mini",
        messages=[
            {"role": "system", "content": "You are a poetic assistant, skilled in explaining complex programming concepts with creative flair."},
            {"role": "user", "content": "Compose a poem that explains the concept of recursion in programming."}
        ]
    )

    print(completion.choices[0].message)

promptModel()

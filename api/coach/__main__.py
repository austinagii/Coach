import json
import pymongo
from pymongo import MongoClient
from openai import OpenAI
from fastapi import FastAPI

from .user import User
from .goal import Goal, Milestone
from .schedule import ScheduledActivity, DailySchedule

app = FastAPI()

@app.post("/users")
def get_user():
     

promptModel()

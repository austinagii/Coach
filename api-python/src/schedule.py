from typing import List
from pydantic import BaseModel, Field

class ScheduledActivity(BaseModel):
    """A singular timeboxed activity that a user wants to complete"""
    title: str 
    description: str
    startTime: str
    endTime: str

class DailySchedule(BaseModel):
    """A daily collection of timeboxed activities that a user wants to complete"""
    routine: List[ScheduledActivity]


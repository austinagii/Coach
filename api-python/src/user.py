from typing import List, Optional
from pydantic import BaseModel, Field

from .goal import Goal
from .schedule import DailySchedule

class User(BaseModel):
    """A new or existing user"""
    id: Optional[str] = Field(default=None)
    name: str 
    summary: Optional[str] = Field(default=None)
    goals: Optional[List[Goal]] = Field(default=None)
    schedule: Optional[DailySchedule] = Field(default=None)


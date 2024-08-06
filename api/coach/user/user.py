from typing import List, Optional
from pydantic.dataclasses import dataclass

from .goal import Goal
from .schedule import DailySchedule

@dataclass
class User:
    """A new or existing user"""
    name: str 
    id: Optional[str] = None 
    summary: Optional[str] = None 
    goals: Optional[List[Goal]] = None 
    schedule: Optional[DailySchedule] = None 


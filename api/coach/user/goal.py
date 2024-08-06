from typing import List, Optional
from pydantic import BaseModel, Field

class Milestone(BaseModel):
    "A milestone to track a user's progress towards completing a goal"
    title: str
    description: str
    targetDate: str
    isComplete: bool = Field(default=False)
    completionDate: Optional[bool] = Field(default=None)
    isDeleted: bool = Field(default=False)

class Goal(BaseModel):
    """A personal goal defined by a user"""
    id: Optional[int]
    title: str
    description: str
    milestones: Optional[List[Milestone]] = Field(default=None)
    isComplete: bool = Field(default=False)
    isDeleted: bool = Field(default=False)
    completionDate: Optional[str] = Field(default=None)


import time
from pymongo import MongoClient
from coach.common.data import MongoFactory

DATABASE_NAME = "aisu"
COLLECTION_NAME = "assistant"

class AssistantRepository:
    def __init__(self, 
                 task: Task,
                 User: user,
                 client: MongoClient,
                 created_at,
                 updated_at):
        """Initialize the assistant repository

        Args:
            client (MongoClient): The database client to be used
        """
        self.id = None
        # Handle errors if database or collection could not be found
        self.collection = mongoClient[DATABASE_NAME][COLLECTION_NAME]
        self.user = None
        self.chat = None
        self.created_at = None
        self.updated_at = None

    def save(assistant: Assistant):
        assistant.date_updated = time.time_ns()


        self.collecion.insert_one()


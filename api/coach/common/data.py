import os
import structlog
from pymongo import MongoClient
from pymongo.errors import PyMongoError
from typing import Optional
from threading import Lock


log = structlog.get_logger()


class MongoClientFactory:
    _mongoClient: Optional[MongoClient] = None
    _lock = Lock()

    @classmethod
    def get_mongo_client(cls) -> MongoClient:
        if cls._mongoClient is None:
            with cls._lock:
                if cls._mongoClient is None:
                    mongoUri = os.getenv("MONGODB_URI")
                    if not mongoUri:
                        log.warning("MONGODB_URI environment variable is not set")
                        raise RuntimeError("MONGODB_URI environment variable is not set")
                    try:
                        cls._mongoClient = MongoClient(mongoUri)
                        log.info("MongoDB client initialized successfully")
                    except PyMongoError as e:
                        log.error(f"An error occurred while initializing the MongoDB client: {e}")
                        log.error("Shutting down...")
                        os._exit(-1)

        return cls._mongoClient

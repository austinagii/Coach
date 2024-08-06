from pymongo import MongoClient
from dataclasses import asdict
import structlog

from coach.common.data import MongoClientFactory
from coach.user import User


log = structlog.get_logger()


class UserRepository:
    def __init__(self, mongoClient: MongoClient):
        """Initializes an instance of a UserRepository.

        Args:
            mongoClient (MongoClient): The mongo client that will be used to interact with the database
        """
        self.collection = mongoClient['aisu']['users']

    def save(self, user: User) -> User:
        """Saves a new or existing user to the database.

        Args:
            user (User): The user to be saved

        Raises:
            MongoError: If an error occurs (TODO: Use accurate error
        """
        user_dict = asdict(user)
        if (user_id := user_dict.get("id")) is None or user_id == "": # User is a new user.
            try:
                result = self.collection.insert_one(user_dict)
            except Exception as e:
                log.error("An error occurred while inserting the user into the database", error=e)
                raise RuntimeException("The following error occurred", e)

            log.info(f"User inserted with id: {result.inserted_id}")
            user_dict["id"] = str(result.inserted_id)
            return User(**user_dict)
        else:
            pass
            # try:
            #     user_dict = user.dict()
            #     user_dict["_id", ObjectId(user.id)]
            #     del user_dict["id"]
            # except (TypeError, InvalidId) as e:
            # except InvalidObjectIdError as e:
            #     return UserNotFoundError("The given user does not exists")
            # log.warning("An error occurred while converting the user id to a mongo object id", user_id=user.id, error=e)
            # user = self.__update__(user)


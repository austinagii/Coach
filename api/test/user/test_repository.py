import pytest
from unittest.mock import MagicMock

from coach.user import User, UserRepository
from coach.common.data import MongoClientFactory

class TestUserRepository:

    @pytest.mark.unit
    def test_pass(self):
         # Create a mock MongoDB client and setup database and collection mocks
        mongoClient = MagicMock()
        database = MagicMock()
        mongoClient.__getitem__.return_value = database

        collection = MagicMock()
        database.__getitem__.return_value = collection

        # Create a mock result object for insert_one
        insert_result = MagicMock()
        inserted_id = "notavalidmongoobjectid"
        insert_result.inserted_id = inserted_id
        collection.insert_one.return_value = insert_result

        # Initialize the repository with the mocked client
        repository = UserRepository(mongoClient)
        user = User(name="Test")

        # Call the save method
        createdUser = repository.save(user)

        # Assert that the user's attributes are as expected
        assert createdUser.name == user.name 
        assert createdUser.id == inserted_id

    @pytest.mark.integration
    def test_pass_integrated(self):
        client = MongoClientFactory.get_mongo_client()
        repository = UserRepository(client)

        user = User(name="Test")

        createdUser = repository.save(user)

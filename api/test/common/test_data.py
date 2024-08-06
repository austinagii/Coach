import pytest
from pymongo import MongoClient

from coach.common.data import MongoClientFactory


class TestMongoClientFactory:

    @pytest.mark.unit
    def test_successfully_creates_mongo_client(self):
        client = MongoClientFactory.get_mongo_client()
        assert isinstance(client, MongoClient), "The mongo client factory did not return an instance of a mongo client"

    @pytest.mark.integration
    def test_mongo_client_is_connected(self):
        client = MongoClientFactory.get_mongo_client()
        client.admin.command("ping")

    @pytest.mark.unit
    def test_returns_mongo_client_singleton(self):
        expected = MongoClientFactory.get_mongo_client()
        actual = MongoClientFactory.get_mongo_client()
        assert actual is expected, "The mongo client factory returned multiple instances of a mongo client"

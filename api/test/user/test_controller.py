import pytest

from coach.user import controller

@pytest.mark.unit
def test_throws_error(self) -> None:
    output = controller.create_user(req)
    expected = '{"error": "name", "description": "name must be between 1 and 50 characters inclusive"}'

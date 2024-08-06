import json
import pytest
import time 

from coach.chat import Sender, Message, Chat

class TestSender:
    @pytest.mark.unit
    def test_sender_is_serializable_to_json(self):
        expected = '{"sender": "user"}'
        actual = json.dumps({"sender": Sender.USER})
        assert actual == expected

    @pytest.mark.unit
    def test_sender_is_deserialzable_from_json(self):
        actual = json.loads('{"sender": "user"}')
        assert Sender(actual["sender"]) is Sender.USER

    @pytest.mark.unit
    def test_value_error_is_raised_for_invalid_sender(self):
        actual = json.loads('{"sender": "users"}')
        with pytest.raises(ValueError):
            Sender(actual["sender"])


class TestMessage:
    @pytest.mark.unit
    def test_default_message_creation_time_is_correct(self):
        message = Message(sender=Sender.USER, content="Hello")
        current_time = time.time_ns() // 1_000_000
        # assert message creation time is within the last second 
        assert (message.created_at >= (current_time - 1_000) 
                and message.created_at <= current_time)

    @pytest.mark.unit
    def test_message_is_correctly_converted_to_a_string(self):
        message = Message(sender=Sender.USER, content="Hello")
        expected = "user: Hello"
        assert str(message) == expected


class TestChat:
    @pytest.mark.unit
    def test_value_error_is_raised_if_message_limit_is_less_than_one(self):
        with pytest.raises(ValueError):
            Chat(message_limit=0)

        with pytest.raises(ValueError):
            Chat(message_limit=-1)

    @pytest.mark.unit
    def test_rotation(self):
        chat = Chat(message_limit=1)
        
        expected = Message(sender=Sender.ASSISTANT, content="Hey there!")
        chat.append(expected)

        actual = chat.append(Message(sender=Sender.USER, content="Hello!"))

        assert len(chat.messages) == 1
        assert expected is actual

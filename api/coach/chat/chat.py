import time
from collections import deque
from dataclasses import dataclass, field
from typing import Optional
from enum import StrEnum


# The default number of messages that a chat can contain.
DEFAULT_MESSAGE_LIMIT: int = 20

# The initiator of a message.
Sender = StrEnum("Sender", ["USER", "ASSISTANT"])

@dataclass(frozen=True)
class Message:
    """A one way communication between a user and an assistant"""

    sender: Sender
    content: str
    created_at: int = field(default_factory=lambda: time.time_ns() // 1_000_000) # unix time in milliseconds

    def __str__(self):
        return f"{str(self.sender)}: {self.content}"

@dataclass
class Chat:
    """A collection of messages exchanged between a user and a assistant"""
    
    message_limit: int = field(default=DEFAULT_MESSAGE_LIMIT)
    messages: deque = field(default_factory=deque)

    def __post_init__(self) -> None:
        if self.message_limit < 1:
            raise ValueError("A chat's message limit must be greater than or equal to 1")

    def append(self, message: Message) -> None:
        """Adds a new message to the chat returning the oldest message if the limit is exceeded.

        If the chat's message limit is reached, the oldest message will be removed from the chat 
        and returned to the caller.

        Args:
            message (Message): The message to be added to the chat

        Returns:
            Optional[Message]: The oldest message in the chat if the limit is reached

        Examples:
            >>> chat = Chat(10)  # Create a chat with a message limit of 10.
            >>> message = Message(sender=Sender.USER, content="Hey there!")
            >>> chat.append(message)
        """
        assert len(self.messages) <= self.message_limit, f"Chat size exceeds limit, size={len(self.messages)}, limit={self.message_limit}"

        self.messages.append(message)
        overflowMessage = None
        if len(self.messages) > self.message_limit:
            overflowMessage = self.messages.popleft()
        return overflowMessage

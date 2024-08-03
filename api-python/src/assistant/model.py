import functools
import logging
import time
from typing import Callable, Optional
import uuid

from openai import OpenAI
from pydantic import BaseModel, Field

log = structlog.getLogger()


class ModelPrompt(BaseModel):
    message: str


class OpenaiLanguageModel:
    """A wrapper for an OpenAI large language model"""
    def __init__(self, client: OpenAI, model: str) -> None:
        self._client = client
        self._model = model

    def __call__(self, systemMessage: str, userMessage: str | Object) -> str:
        """Prompts the model using the given messages and returns the response

        Args:
            systemMessage (str):
            userMessage (str):
        Returns
            str:

        Examples:
            >>> gpt4o = OpenaiLanguageModel(client=client, model="gpt-4o")
            >>> gpt4o("You're an assistant", "Teach me algebra")
            Sure, I can teach you algebra. They key to algebra....
        """
        exchangeId = uuid.uuid4()
        log.info(f"Prompting {self.model}...", prompt_id=exchangeId)

        exchange = ModelExchange(exchangeId, systemMessage, userMessage)
        # TODO: Add millisecond component for consistency.
        exchange.promptedAt = time.strftime("%Y/%m/%d %H/%M/%s")
        completion = self.client.chat.completions.create(
            model=self.model,
            messages=[
                {"role": "system", "content": systemMessage},
                {"role": "user", "content": userMessage}
            ]
        )

        exchange.response = completion.choices[0].message
        exchange.respondedAt = time.strftime("%Y/%m/%d %H/%M/%s")

        asyncio.create_task(auditExchange(exchange))
        
        return exchangeId, exchange.response



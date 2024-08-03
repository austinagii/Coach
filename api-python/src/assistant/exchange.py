from .model import ModelPrompt

class ModelExchange(BaseModel):
    id: str
    systemMessage: str
    userMessage: ModelPrompt
    assistantMessage: ModelResponse
    promptedAt: str
    respondedAt: str

    def __init__(self, 
                 systemMessage: str, 
                 prompt: str,
                 response: str) -> None:
        """Initializes a ModelExchange instance.

        Args:
            systemMessage (str):
            prompt (str):
            promptedAt (str):
            response (str):
            respondedAt (str):
        """
        self.id = uuid.uuid4()
        self.systemMessage = systemMessage
        self.prompt = prompt



# The generic description of an assistant excluding any details about it's current task.
ASSISTANT_DESCRIPTION: str = ""

# Maps objectives to tailored user prompts that an assistant will use to start a chat.
# e.g. If a user requests a new assistant with the objective of goal creation, the message
# 'What's goal do you want to set?' will be used to prompt the user to start defining their goal.
CHAT_PROMPT_BY_OBJECTIVE: Dict[str] = {}

# ErrNotIntialized indicates that the creation of a new Assistant was attempted before loading the
# static data required.
ERR_NOT_INITIALIZED = errors.New("The static data required to create a new assistant has not been loaded")


class Assistant:
    def __init__(self, 
                 task: Task,
                 user: User):
        self.id = None
        self.task = None
        self.user = None
        self.chat = None
        self.client = None
        self.model_exchange_repo = None
        self.created_at = None
        self.updated_at = None





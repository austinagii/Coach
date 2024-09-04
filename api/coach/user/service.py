class UserService:
    def __init__(self, repository: UserRepository):
        self.repository = repository

    def save(user: User) -> User:
        """Persists the given user to the database"""

        try:
            saved_user = self.repository.save(user)
        except MongoError as e:
            pass
        except Exception as e:
            pass

        return saved_user

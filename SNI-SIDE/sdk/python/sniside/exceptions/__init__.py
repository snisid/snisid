class SNISIDError(Exception):
    def __init__(self, status_code: int, message: str):
        self.status_code = status_code
        super().__init__(f"[{status_code}] {message}")


class AuthenticationError(SNISIDError):
    def __init__(self, message: str):
        super().__init__(401, message)


class PermissionError(SNISIDError):
    def __init__(self, message: str):
        super().__init__(403, message)


class NotFoundError(SNISIDError):
    def __init__(self, message: str):
        super().__init__(404, message)


class RateLimitError(SNISIDError):
    def __init__(self, message: str):
        super().__init__(429, message)


class ServerError(SNISIDError):
    def __init__(self, message: str):
        super().__init__(500, message)

class NIREResolver:
    def __init__(self, db_url: str = "sqlite:///:memory:"):
        self.db_url = db_url

    def validate_nire(self, nire: str) -> bool:
        if not nire or len(nire) != 15:
            return False
        return nire[:2].isalpha() and nire[2:].isdigit()

    def format_nire(self, year: int, dept: str, serial: int) -> str:
        return f"{dept}{year:04d}{serial:08d}"

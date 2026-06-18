import hashlib, uuid, re
from datetime import datetime


NIU_CHARS = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"


def generate_niu(seed: str = None) -> str:
    if seed:
        h = hashlib.sha256(str(seed).encode()).hexdigest().upper()
        niu = "".join(c for c in h if c in NIU_CHARS)[:10]
        if len(niu) < 10:
            niu = niu + "X" * (10 - len(niu))
        return niu
    return str(uuid.uuid4()).replace("-", "").upper()[:10]


def generate_nui(prefix: str = "NIU") -> str:
    return prefix + str(uuid.uuid4()).replace("-", "").upper()[:7]


def normalize_phone(phone: str, country_code: str = "509") -> str:
    digits = re.sub(r"\D", "", str(phone))
    if len(digits) == 8:
        return f"+{country_code}{digits}"
    if digits.startswith(country_code) and len(digits) == 11:
        return f"+{digits}"
    if digits.startswith("+"):
        return digits
    return f"+{country_code}{digits[-8:]}" if len(digits) >= 8 else phone


def normalize_plate(plate: str) -> str:
    return re.sub(r"\s+", "-", str(plate).upper().strip())[:15]


def normalize_name(name: str) -> str:
    name = str(name).strip().upper()
    name = re.sub(r"\s+", " ", name)
    replacements = {
        "Á": "A", "À": "A", "Â": "A", "Ã": "A", "Ä": "A",
        "É": "E", "È": "E", "Ê": "E", "Ë": "E",
        "Í": "I", "Ì": "I", "Î": "I", "Ï": "I",
        "Ó": "O", "Ò": "O", "Ô": "O", "Õ": "O", "Ö": "O",
        "Ú": "U", "Ù": "U", "Û": "U", "Ü": "U",
        "Ç": "C", "Ñ": "N",
    }
    for a, b in replacements.items():
        name = name.replace(a, b)
    return name


def parse_date(value: str, formats: list = None) -> str:
    if not formats:
        formats = ["%d/%m/%Y", "%Y-%m-%d", "%m/%d/%Y", "%d-%m-%Y", "%Y%m%d", "%d %B %Y"]
    for fmt in formats:
        try:
            return datetime.strptime(str(value)[:25].strip(), fmt).isoformat()
        except (ValueError, IndexError):
            continue
    return None

import re
from datetime import datetime


def validate(record: dict, validations: list) -> list:
    errors = []
    for v in validations:
        field = v.get("field")
        rule = v.get("rule")
        value = record.get(field)
        msg = v.get("message", f"Validation échouée pour '{field}': {rule}")

        if rule == "required":
            if value is None or (isinstance(value, str) and value.strip() == ""):
                errors.append(msg)

        elif rule == "niu_format":
            if value and not re.match(r"^[A-Z0-9]{10}$", str(value)):
                errors.append(msg)

        elif rule == "email":
            if value and not re.match(r"^[^@\s]+@[^@\s]+\.[^@\s]+$", str(value)):
                errors.append(msg)

        elif rule == "phone":
            if value:
                digits = re.sub(r"\D", "", str(value))
                if len(digits) < 8:
                    errors.append(msg)

        elif rule == "plate":
            if value and not re.match(r"^[A-Z0-9-]{2,15}$", str(value).upper()):
                errors.append(msg)

        elif rule == "date":
            if value:
                fmt = v.get("format", "%Y-%m-%d")
                try:
                    datetime.strptime(str(value)[:25].strip(), fmt)
                except ValueError:
                    errors.append(msg)

        elif rule == "not_future":
            if value:
                try:
                    dt = datetime.fromisoformat(str(value)) if "T" in str(value) else datetime.strptime(str(value)[:10], "%Y-%m-%d")
                    if dt > datetime.utcnow():
                        errors.append(msg)
                except (ValueError, IndexError):
                    pass

        elif rule == "min_length":
            min_len = v.get("min", 1)
            if value and len(str(value)) < min_len:
                errors.append(msg)

        elif rule == "max_length":
            max_len = v.get("max", 255)
            if value and len(str(value)) > max_len:
                errors.append(msg)

        elif rule == "regex":
            pattern = v.get("pattern", r".*")
            if value and not re.match(pattern, str(value)):
                errors.append(msg)

        elif rule == "enum":
            allowed = v.get("values", [])
            if value and str(value) not in allowed:
                errors.append(msg)

        elif rule == "range":
            min_v = v.get("min", float("-inf"))
            max_v = v.get("max", float("inf"))
            try:
                num = float(value)
                if num < min_v or num > max_v:
                    errors.append(msg)
            except (TypeError, ValueError):
                errors.append(msg)

    return errors

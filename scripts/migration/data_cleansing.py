from __future__ import annotations

import re
import unicodedata
from datetime import datetime
from typing import Dict, List, Optional, Set, Tuple

from dateutil.parser import parse as dateutil_parse

from config import CleansingRules


class NameNormalizer:
    HAITIAN_ABBREVIATIONS = {
        "JN": "JEAN",
        "JNE": "JEANNE",
        "JE": "JEAN",
        "ST": "SAINT",
        "STE": "SANTE",
        "FT": "FRANTZ",
        "ME": "MARIE",
        "MT": "MARTIN",
        "LA": "LA",
        "DU": "DU",
        "DE": "DE",
        "DES": "DES",
        "DEL": "DEL",
    }

    HAITIAN_CREOLE_MAP = {
        "BATIS": "BAPTISTE",
        "JAN": "JEAN",
        "JANBATIS": "JEAN-BAPTISTE",
        "OKAP": "CAP-HAITIEN",
        "P-AU-P": "PORT-AU-PRINCE",
        "PORTAUPRINCE": "PORT-AU-PRINCE",
        "W": "OU",
        "K": "C",
        "C": "K",
    }

    BOGUS_VALUES: Set[str] = {
        "N/A", "N/A ", "NA", "NULL", "TEST", "INCONNU", "UNKNOWN",
        "???", "***", "---", "...", " ", "NONE", "VIDE", "SANS",
    }

    @classmethod
    def normalize_last_name(cls, name: Optional[str]) -> str:
        if not name or not name.strip():
            return "UNKNOWN"
        name = cls._clean_text(name)
        name = cls._expand_abbreviations(name)
        name = cls._apply_phonetic_rules(name)
        name = name.upper().strip()
        if not name or name in cls.BOGUS_VALUES:
            return "UNKNOWN"
        return name

    @classmethod
    def normalize_first_name(cls, name: Optional[str]) -> str:
        if not name or not name.strip():
            return "UNKNOWN"
        name = cls._clean_text(name)
        parts = name.split()
        normalized_parts = []
        for part in parts:
            p = part.strip()
            if p.upper() in cls.BOGUS_VALUES:
                continue
            if p.upper() in cls.HAITIAN_ABBREVIATIONS:
                p = cls.HAITIAN_ABBREVIATIONS[p.upper()]
            p = p.capitalize()
            normalized_parts.append(p)
        result = " ".join(normalized_parts).strip()
        return result if result else "UNKNOWN"

    @classmethod
    def normalize_full_name(cls, first_name: Optional[str], last_name: Optional[str]) -> Tuple[str, str]:
        return cls.normalize_first_name(first_name), cls.normalize_last_name(last_name)

    @classmethod
    def _clean_text(cls, text: str) -> str:
        normalized = unicodedata.normalize("NFD", text)
        ascii_text = "".join(c for c in normalized if unicodedata.category(c) != "Mn")
        ascii_text = ascii_text.strip().upper()
        ascii_text = re.sub(r"[^A-Z\-\'\s]", "", ascii_text)
        ascii_text = re.sub(r"\s+", " ", ascii_text)
        for bogus in cls.BOGUS_VALUES:
            ascii_text = re.sub(rf"\b{re.escape(bogus)}\b", "", ascii_text)
        return ascii_text.strip()

    @classmethod
    def _expand_abbreviations(cls, text: str) -> str:
        words = text.split()
        expanded = []
        for w in words:
            w_upper = w.upper()
            if w_upper in cls.HAITIAN_ABBREVIATIONS:
                expanded.append(cls.HAITIAN_ABBREVIATIONS[w_upper])
            else:
                expanded.append(w)
        return " ".join(expanded)

    @classmethod
    def _apply_phonetic_rules(cls, text: str) -> str:
        s = text.upper()
        s = re.sub(r"\bHYPPOLITE\b", "HIPPOLYTE", s)
        s = re.sub(r"\bCHARLEMAGNE\b", "CHARLEMAGNE", s)
        s = re.sub(r"\bCHRISTOPHE\b", "CHRISTOPHE", s)
        s = re.sub(r"EAU", "O", s)
        s = re.sub(r"AUX", "O", s)
        s = re.sub(r"PH", "F", s)
        s = re.sub(r"QU", "K", s)
        s = s.replace("-", " ").strip()
        s = re.sub(r"\s+", "-", s)
        return s


class DateNormalizer:
    KNOWN_FORMATS = [
        "%d/%m/%Y",
        "%d-%m-%Y",
        "%Y-%m-%d",
        "%m/%d/%Y",
        "%m-%d-%Y",
        "%Y/%m/%d",
        "%d.%m.%Y",
        "%Y.%m.%d",
        "%d %m %Y",
        "%Y %m %d",
        "%d/%m/%y",
        "%d-%m-%y",
        "%m/%d/%y",
        "%Y%m%d",
        "%d%m%Y",
    ]

    FRENCH_MONTHS = {
        "JANVIER": 1, "FEVRIER": 2, "FÉVRIER": 2, "MARS": 3, "AVRIL": 4,
        "MAI": 5, "JUIN": 6, "JUILLET": 7, "AOUT": 8, "AOÛT": 8,
        "SEPTEMBRE": 9, "OCTOBRE": 10, "NOVEMBRE": 11, "DECEMBRE": 12, "DÉCEMBRE": 12,
    }

    HAITIAN_CREOLE_MONTHS = {
        "JANVYE": 1, "FEVRIYE": 2, "MAS": 3, "AVRIL": 4,
        "ME": 5, "JEN": 6, "JIYE": 7, "DAO": 8,
        "SEPTANM": 9, "OKTOB": 10, "NOVANM": 11, "DESANM": 12,
    }

    @classmethod
    def normalize(cls, date_str: Optional[str]) -> Tuple[Optional[str], Optional[str]]:
        if not date_str or not date_str.strip():
            return None, "Missing date value"

        date_str = date_str.strip()

        try:
            parsed = dateutil_parse(date_str, dayfirst=True, fuzzy=False)
            return cls._validate_and_format(parsed)
        except (ValueError, OverflowError):
            pass

        for fmt in cls.KNOWN_FORMATS:
            try:
                parsed = datetime.strptime(date_str, fmt)
                return cls._validate_and_format(parsed)
            except (ValueError, OverflowError):
                continue

        result = cls._try_french_text_date(date_str)
        if result[0]:
            return result

        result = cls._try_creole_text_date(date_str)
        if result[0]:
            return result

        return None, f"Unrecognized date format: {date_str[:50]}"

    @classmethod
    def _validate_and_format(cls, dt: datetime) -> Tuple[Optional[str], Optional[str]]:
        if dt.year < 1900 or dt.year > 2026:
            return None, f"Year out of valid range (1900-2026): {dt.year}"
        return dt.strftime("%Y-%m-%d"), None

    @classmethod
    def _try_french_text_date(cls, date_str: str) -> Tuple[Optional[str], Optional[str]]:
        date_upper = date_str.upper().strip()
        for month_name, month_num in cls.FRENCH_MONTHS.items():
            if month_name in date_upper:
                parts = re.split(r"[\s/\\-]+", date_upper)
                day = None
                year = None
                for p in parts:
                    p_clean = p.strip()
                    if p_clean.isdigit():
                        val = int(p_clean)
                        if val > 31:
                            year = val
                        elif day is None:
                            day = val
                if day is not None and year is not None:
                    try:
                        dt = datetime(year, month_num, day)
                        return cls._validate_and_format(dt)
                    except (ValueError, OverflowError):
                        return None, f"Invalid French date: {date_str}"
        return None, None

    @classmethod
    def _try_creole_text_date(cls, date_str: str) -> Tuple[Optional[str], Optional[str]]:
        date_upper = date_str.upper().strip()
        for month_name, month_num in cls.HAITIAN_CREOLE_MONTHS.items():
            if month_name in date_upper:
                parts = re.split(r"[\s/\\-]+", date_upper)
                day = None
                year = None
                for p in parts:
                    p_clean = p.strip()
                    if p_clean.isdigit():
                        val = int(p_clean)
                        if val > 31:
                            year = val
                        elif day is None:
                            day = val
                if day is not None and year is not None:
                    try:
                        dt = datetime(year, month_num, day)
                        return cls._validate_and_format(dt)
                    except (ValueError, OverflowError):
                        return None, f"Invalid Creole date: {date_str}"
        return None, None

    @classmethod
    def get_age(cls, date_str: str, reference_date: Optional[datetime] = None) -> Optional[int]:
        normalized, err = cls.normalize(date_str)
        if err or not normalized:
            return None
        dt = datetime.strptime(normalized, "%Y-%m-%d")
        ref = reference_date or datetime.now()
        age = ref.year - dt.year
        if (ref.month, ref.day) < (dt.month, dt.day):
            age -= 1
        return age

    @classmethod
    def is_valid_birth_date(cls, date_str: str) -> bool:
        _, err = cls.normalize(date_str)
        return err is None


class PhoneNormalizer:
    HAITI_COUNTRY_CODE = "509"
    VALID_PREFIXES = {"3", "4", "5", "6", "7", "8", "2"}
    HAITIAN_OPERATORS = {
        "3": "Digicel",
        "4": "Digicel",
        "5": "Natcom",
        "6": "Natcom",
        "7": "Starlink",
        "8": "Access Haiti",
        "2": "Teleco",
    }

    @classmethod
    def normalize(cls, phone: Optional[str]) -> Tuple[Optional[str], Optional[str]]:
        if not phone or not phone.strip():
            return None, "Missing phone number"

        cleaned = re.sub(r"[^\d+]", "", phone.strip())

        if cleaned.startswith("+"):
            if cleaned.startswith(f"+{cls.HAITI_COUNTRY_CODE}"):
                cleaned = cleaned[1:]
            else:
                return None, f"Non-Haitian country code in {phone}"

        if cleaned.startswith("00"):
            cleaned = cleaned[2:]

        if cleaned.startswith(cls.HAITI_COUNTRY_CODE) and len(cleaned) > 8:
            local = cleaned[len(cls.HAITI_COUNTRY_CODE):]
            if len(local) == 8:
                cleaned = local

        if len(cleaned) == 8 and cleaned.isdigit():
            prefix = cleaned[0]
            if prefix in cls.VALID_PREFIXES:
                normalized = f"+{cls.HAITI_COUNTRY_CODE}{cleaned}"
                return normalized, None
            return None, f"Invalid Haitian prefix '{prefix}' in {phone}"

        if len(cleaned) == 10 and cleaned.isdigit() and cleaned.startswith("509"):
            local = cleaned[3:]
            prefix = local[0]
            if prefix in cls.VALID_PREFIXES:
                normalized = f"+{cls.HAITI_COUNTRY_CODE}{local}"
                return normalized, None
            return None, f"Invalid Haitian prefix in full code: {phone}"

        return None, f"Unrecognized phone format: {phone}"

    @classmethod
    def get_operator(cls, phone: str) -> Optional[str]:
        normalized, err = cls.normalize(phone)
        if err or not normalized:
            return None
        local = normalized.replace(f"+{cls.HAITI_COUNTRY_CODE}", "")
        prefix = local[0] if local else ""
        return cls.HAITIAN_OPERATORS.get(prefix)

    @classmethod
    def format_display(cls, phone: str) -> Optional[str]:
        normalized, err = cls.normalize(phone)
        if err or not normalized:
            return None
        local = normalized.replace(f"+{cls.HAITI_COUNTRY_CODE}", "")
        return f"+{cls.HAITI_COUNTRY_CODE} {local[:3]}-{local[3:5]}-{local[5:]}"


class AddressNormalizer:
    DEPARTMENT_ALIASES = {
        "OUEST": "OUEST", "L'OUEST": "OUEST", "WÈS": "OUEST",
        "SUD": "SUD", "SID": "SUD",
        "NORD": "NORD", "NÒ": "NORD",
        "ARTIBONITE": "ARTIBONITE", "L'ARTIBONITE": "ARTIBONITE", "ATIBONIT": "ARTIBONITE",
        "CENTRE": "CENTRE", "SANT": "CENTRE",
        "SUD-EST": "SUD-EST", "SUDEST": "SUD-EST", "SIDÈS": "SUD-EST",
        "NORD-OUEST": "NORD-OUEST", "NORDOUEST": "NORD-OUEST", "NÒDWÈS": "NORD-OUEST",
        "NORD-EST": "NORD-EST", "NORDEST": "NORD-EST", "NÒDÈS": "NORD-EST",
        "GRANDE-ANSE": "GRANDE-ANSE", "GRANDE ANSE": "GRANDE-ANSE", "GRANDANS": "GRANDE-ANSE",
        "NIPES": "NIPPES", "NIPPES": "NIPPES",
    }

    COMMUNE_ALIASES = {
        "PORT-AU-PRINCE": "PORT-AU-PRINCE", "PÖTOPRENS": "PORT-AU-PRINCE",
        "P-AU-P": "PORT-AU-PRINCE", "PAP": "PORT-AU-PRINCE",
        "CAP-HAITIEN": "CAP-HAITIEN", "OKAP": "CAP-HAITIEN", "LE CAP": "CAP-HAITIEN",
        "JACMEL": "JACMEL", "YAKMEL": "JACMEL",
        "GONAIVES": "GONAÏVES", "GONAÏVES": "GONAÏVES", "GONAYIV": "GONAÏVES",
        "CAYES": "LES CAYES", "LES CAYES": "LES CAYES", "OKAY": "LES CAYES",
        "PETION-VILLE": "PÉTION-VILLE", "PETIONVILLE": "PÉTION-VILLE", "PÉTION-VILLE": "PÉTION-VILLE",
        "DELMAS": "DELMAS", "DÈLMA": "DELMAS",
        "CARREFOUR": "CARREFOUR", "KAREFOU": "CARREFOUR",
    }

    ADDRESS_KEYWORDS = {
        "RUE": "RUE", "RV": "RUE", "RU": "RUE",
        "BLVD": "BOULEVARD", "BOULEVARD": "BOULEVARD", "BVL": "BOULEVARD",
        "AVE": "AVENUE", "AV": "AVENUE", "AVENUE": "AVENUE",
        "IMP": "IMPASSE", "IMPASSE": "IMPASSE",
        "CITE": "CITÉ", "CITÉ": "CITÉ", "SITE": "CITÉ",
        "RD": "ROUTE", "ROUTE": "ROUTE", "RT": "ROUTE",
        "CHEMIN": "CHEMIN", "CH": "CHEMIN",
        "PLACE": "PLACE", "PL": "PLACE",
        "BAT": "BÂTIMENT", "BATIMENT": "BÂTIMENT", "BT": "BÂTIMENT",
    }

    @classmethod
    def normalize(cls, address: Optional[str]) -> Tuple[Optional[str], Optional[str]]:
        if not address or not address.strip():
            return None, "Missing address"

        cleaned = cls._clean_text(address)
        cleaned = cls._normalize_keywords(cleaned)
        cleaned = cls._normalize_department(cleaned)
        cleaned = cls._normalize_commune(cleaned)
        cleaned = cls._format_address(cleaned)

        if len(cleaned) < 3:
            return "UNKNOWN", None

        return cleaned, None

    @classmethod
    def _clean_text(cls, text: str) -> str:
        normalized = unicodedata.normalize("NFD", text)
        ascii_text = "".join(c for c in normalized if unicodedata.category(c) != "Mn")
        ascii_text = ascii_text.strip().upper()
        ascii_text = re.sub(r"[^A-Z0-9\-\'\s,/]", " ", ascii_text)
        ascii_text = re.sub(r"\s+", " ", ascii_text)
        return ascii_text.strip()

    @classmethod
    def _normalize_keywords(cls, text: str) -> str:
        words = text.split()
        normalized = []
        for w in words:
            w_upper = w.upper().strip(",. ")
            if w_upper in cls.ADDRESS_KEYWORDS:
                normalized.append(cls.ADDRESS_KEYWORDS[w_upper])
            elif w_upper.endswith("S") and w_upper not in ["NO", "NUM"] and len(w_upper) <= 3:
                pass
            else:
                normalized.append(w.capitalize())
        return " ".join(normalized)

    @classmethod
    def _normalize_department(cls, text: str) -> str:
        for alias, dept in cls.DEPARTMENT_ALIASES.items():
            text = re.sub(rf"\b{re.escape(alias)}\b", dept, text, flags=re.IGNORECASE)
        return text

    @classmethod
    def _normalize_commune(cls, text: str) -> str:
        for alias, commune in cls.COMMUNE_ALIASES.items():
            text = re.sub(rf"\b{re.escape(alias)}\b", commune, text, flags=re.IGNORECASE)
        return text

    @classmethod
    def _format_address(cls, text: str) -> str:
        words = text.split()
        formatted = []
        for w in words:
            if w.isdigit():
                formatted.append(w)
            elif len(w) <= 2 and w.upper() in ["NO", "BP", "CP", "APT", "NUM"]:
                formatted.append(w.upper())
            elif w.isupper() and len(w) > 2:
                formatted.append(w.capitalize())
            else:
                formatted.append(w)
        return " ".join(formatted)


class IDNormalizer:
    CIN_PATTERN = re.compile(r"^(\d{1,3})-?(\d{1,3})-?(\d{1,3})-?(\d{1,2})$")
    NIF_PATTERN = re.compile(r"^(\d{3})-?(\d{3})-?(\d{3})-?(\d{1})$")
    PASSPORT_PATTERN = re.compile(r"^([A-Z]{2})\d{6,8}$", re.IGNORECASE)

    @classmethod
    def normalize_cin(cls, cin: Optional[str]) -> Tuple[Optional[str], Optional[str]]:
        if not cin or not cin.strip():
            return None, "Missing CIN"
        cleaned = re.sub(r"[^\d-]", "", cin.strip()).upper()
        match = cls.CIN_PATTERN.match(cleaned)
        if match:
            parts = match.groups()
            normalized = f"{parts[0]}-{parts[1]}-{parts[2]}-{parts[3]}"
            return normalized, None
        if cleaned.replace("-", "").isdigit() and len(cleaned.replace("-", "")) >= 8:
            digits = cleaned.replace("-", "")
            normalized = f"{digits[:3]}-{digits[3:6]}-{digits[6:9]}-{digits[9:11]}" if len(digits) >= 11 else digits
            return normalized, None
        return None, f"Invalid CIN format: {cin}"

    @classmethod
    def normalize_nif(cls, nif: Optional[str]) -> Tuple[Optional[str], Optional[str]]:
        if not nif or not nif.strip():
            return None, "Missing NIF"
        cleaned = re.sub(r"[^\d-]", "", nif.strip())
        if len(cleaned) == 10 and cleaned.isdigit():
            return f"{cleaned[:3]}-{cleaned[3:6]}-{cleaned[6:9]}-{cleaned[9]}", None
        match = cls.NIF_PATTERN.match(cleaned)
        if match:
            parts = match.groups()
            return f"{parts[0]}-{parts[1]}-{parts[2]}-{parts[3]}", None
        return None, f"Invalid NIF format: {nif}"

    @classmethod
    def normalize_passport(cls, passport: Optional[str]) -> Tuple[Optional[str], Optional[str]]:
        if not passport or not passport.strip():
            return None, "Missing passport"
        cleaned = passport.strip().upper().replace(" ", "")
        if len(cleaned) >= 8 and cleaned[:2].isalpha() and cleaned[2:].isdigit():
            return cleaned, None
        return None, f"Invalid passport format: {passport}"

    @classmethod
    def is_valid_cin(cls, cin: str) -> bool:
        _, err = cls.normalize_cin(cin)
        return err is None


class DataCleansingEngine:
    def __init__(self, config: Optional[CleansingRules] = None):
        from config import CleansingRules
        self.config = config or CleansingRules()

    def cleanse(self, record: Dict) -> Tuple[Dict, List[str]]:
        errors = []
        cleaned = dict(record)

        if self.config.normalize_names:
            cleaned, name_errors = self._cleanse_names(cleaned, record)
            errors.extend(name_errors)

        if self.config.normalize_dates:
            cleaned, date_errors = self._cleanse_dates(cleaned, record)
            errors.extend(date_errors)

        if self.config.normalize_phones:
            cleaned, phone_errors = self._cleanse_phones(cleaned, record)
            errors.extend(phone_errors)

        if self.config.normalize_addresses:
            cleaned, addr_errors = self._cleanse_addresses(cleaned, record)
            errors.extend(addr_errors)

        cleaned = self._cleanse_ids(cleaned, record)

        if self.config.strip_special_chars:
            cleaned = self._strip_special_chars(cleaned)

        self._ensure_required_fields(cleaned, errors)

        return cleaned, errors

    def _cleanse_names(self, cleaned: Dict, original: Dict) -> Tuple[Dict, List[str]]:
        errors = []
        for field in ["nom", "last_name", "surname", "lastname"]:
            if field in cleaned:
                cleaned[field] = NameNormalizer.normalize_last_name(cleaned[field])
                break
        else:
            for key in cleaned:
                if key.lower() in ["nom", "last_name", "surname", "lastname"]:
                    cleaned[key] = NameNormalizer.normalize_last_name(cleaned[key])
                    break

        for field in ["prenom", "first_name", "given_name", "firstname"]:
            if field in cleaned:
                cleaned[field] = NameNormalizer.normalize_first_name(cleaned[field])
                break
        else:
            for key in cleaned:
                if key.lower() in ["prenom", "first_name", "given_name", "firstname"]:
                    cleaned[key] = NameNormalizer.normalize_first_name(cleaned[key])
                    break

        if cleaned.get("last_name") == "UNKNOWN" or cleaned.get("first_name") == "UNKNOWN":
            errors.append("Name normalization produced UNKNOWN")
        return cleaned, errors

    def _cleanse_dates(self, cleaned: Dict, original: Dict) -> Tuple[Dict, List[str]]:
        errors = []
        date_fields = ["birth_date", "date_naissance", "date_naiss", "dob", "date_of_birth",
                       "created_date", "date_creation", "updated_date", "date_modification"]
        for field in date_fields:
            if field in cleaned and cleaned[field]:
                normalized, err = DateNormalizer.normalize(cleaned[field])
                if normalized:
                    cleaned[field] = normalized
                elif err:
                    errors.append(f"{field}: {err}")
        return cleaned, errors

    def _cleanse_phones(self, cleaned: Dict, original: Dict) -> Tuple[Dict, List[str]]:
        errors = []
        phone_fields = ["phone", "telephone", "tel", "mobile", "phone_number", "cell", "telephone_mobile"]
        for field in phone_fields:
            if field in cleaned and cleaned[field]:
                normalized, err = PhoneNormalizer.normalize(cleaned[field])
                if normalized:
                    cleaned[field] = normalized
                elif err:
                    errors.append(f"{field}: {err}")
        return cleaned, errors

    def _cleanse_addresses(self, cleaned: Dict, original: Dict) -> Tuple[Dict, List[str]]:
        errors = []
        addr_fields = ["address", "adresse", "birth_place", "lieu_naissance",
                       "residence", "domicile", "street", "rue"]
        for field in addr_fields:
            if field in cleaned and cleaned[field]:
                normalized, err = AddressNormalizer.normalize(cleaned[field])
                if normalized:
                    cleaned[field] = normalized
                elif err:
                    errors.append(f"{field}: {err}")
        return cleaned, errors

    def _cleanse_ids(self, cleaned: Dict, original: Dict) -> Dict:
        for field in ["cin", "nif", "passport"]:
            if field in cleaned and cleaned[field]:
                if field == "cin":
                    normalized, _ = IDNormalizer.normalize_cin(cleaned[field])
                elif field == "nif":
                    normalized, _ = IDNormalizer.normalize_nif(cleaned[field])
                elif field == "passport":
                    normalized, _ = IDNormalizer.normalize_passport(cleaned[field])
                else:
                    normalized = None
                if normalized:
                    cleaned[field] = normalized
        return cleaned

    def _strip_special_chars(self, cleaned: Dict) -> Dict:
        for key, value in cleaned.items():
            if isinstance(value, str):
                cleaned[key] = re.sub(r"[\x00-\x08\x0B\x0C\x0E-\x1F]", "", value)
        return cleaned

    def _ensure_required_fields(self, cleaned: Dict, errors: List[str]) -> None:
        has_last = any(v.lower() == "last_name" or v == "nom" for v in cleaned)
        has_first = any(v.lower() == "first_name" or v == "prenom" for v in cleaned)
        if not has_last:
            if "last_name" not in cleaned and "nom" not in cleaned:
                errors.append("Missing required field: last_name/nom")
        if not has_first:
            if "first_name" not in cleaned and "prenom" not in cleaned:
                errors.append("Missing required field: first_name/prenom")

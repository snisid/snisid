"""Tests for DataCleansingEngine and normalizer components."""

import unittest
from datetime import datetime

from data_cleansing import (
    DataCleansingEngine,
    NameNormalizer,
    DateNormalizer,
    PhoneNormalizer,
    AddressNormalizer,
    IDNormalizer,
)
from config import CleansingRules


class TestNameNormalization(unittest.TestCase):
    def test_normalize_last_name_with_accents(self):
        result = NameNormalizer.normalize_last_name("Méçan")
        self.assertEqual(result, "MECAN")

    def test_normalize_last_name_uppercase(self):
        result = NameNormalizer.normalize_last_name("DUPONT")
        self.assertEqual(result, "DUPONT")

    def test_normalize_last_name_mixed_case(self):
        result = NameNormalizer.normalize_last_name("DuPont")
        self.assertEqual(result, "DUPONT")

    def test_normalize_last_name_with_special_chars(self):
        result = NameNormalizer.normalize_last_name("O'Brien")
        self.assertEqual(result, "OBRIEN")

    def test_normalize_last_name_empty(self):
        result = NameNormalizer.normalize_last_name("")
        self.assertEqual(result, "UNKNOWN")

    def test_normalize_last_name_none(self):
        result = NameNormalizer.normalize_last_name(None)
        self.assertEqual(result, "UNKNOWN")

    def test_normalize_last_name_bogus(self):
        result = NameNormalizer.normalize_last_name("N/A")
        self.assertEqual(result, "UNKNOWN")

    def test_normalize_first_name_with_abbreviation(self):
        result = NameNormalizer.normalize_first_name("JN")
        self.assertEqual(result, "Jean")

    def test_normalize_first_name_multiple_parts(self):
        result = NameNormalizer.normalize_first_name("jean baptiste")
        self.assertEqual(result, "Jean Baptiste")

    def test_normalize_first_name_empty(self):
        result = NameNormalizer.normalize_first_name("")
        self.assertEqual(result, "UNKNOWN")

    def test_normalize_last_name_haitian_creole(self):
        result = NameNormalizer.normalize_last_name("BATIS")
        self.assertEqual(result, "BAPTISTE")

    def test_normalize_full_name(self):
        fn, ln = NameNormalizer.normalize_full_name("jean", "dupont")
        self.assertEqual(fn, "Jean")
        self.assertEqual(ln, "DUPONT")


class TestDateNormalization(unittest.TestCase):
    def test_date_slash_format_dmy(self):
        result, err = DateNormalizer.normalize("15/03/2020")
        self.assertIsNone(err)
        self.assertEqual(result, "2020-03-15")

    def test_date_dash_format_ymd(self):
        result, err = DateNormalizer.normalize("2020-03-15")
        self.assertIsNone(err)
        self.assertEqual(result, "2020-03-15")

    def test_date_dot_format(self):
        result, err = DateNormalizer.normalize("15.03.2020")
        self.assertIsNone(err)
        self.assertEqual(result, "2020-03-15")

    def test_date_compact_format(self):
        result, err = DateNormalizer.normalize("20200315")
        self.assertIsNone(err)
        self.assertEqual(result, "2020-03-15")

    def test_date_french_text(self):
        result, err = DateNormalizer.normalize("15 MARS 2020")
        self.assertIsNone(err)
        self.assertEqual(result, "2020-03-15")

    def test_date_creole_text(self):
        result, err = DateNormalizer.normalize("15 MAS 2020")
        self.assertIsNone(err)
        self.assertEqual(result, "2020-03-15")

    def test_date_empty(self):
        result, err = DateNormalizer.normalize("")
        self.assertIsNone(result)
        self.assertIsNotNone(err)

    def test_date_none(self):
        result, err = DateNormalizer.normalize(None)
        self.assertIsNone(result)
        self.assertIsNotNone(err)

    def test_date_invalid(self):
        result, err = DateNormalizer.normalize("not-a-date")
        self.assertIsNone(result)
        self.assertIsNotNone(err)

    def test_date_year_out_of_range(self):
        result, err = DateNormalizer.normalize("01/01/1800")
        self.assertIsNone(result)
        self.assertIsNotNone(err)


class TestPhoneNormalization(unittest.TestCase):
    def test_phone_local_format(self):
        result, err = PhoneNormalizer.normalize("36305123")
        self.assertIsNone(err)
        self.assertEqual(result, "+50936305123")

    def test_phone_haiti_country_code(self):
        result, err = PhoneNormalizer.normalize("+50936305123")
        self.assertIsNone(err)
        self.assertEqual(result, "+50936305123")

    def test_phone_double_zero_prefix(self):
        result, err = PhoneNormalizer.normalize("0050936305123")
        self.assertIsNone(err)
        self.assertEqual(result, "+50936305123")

    def test_phone_ten_digit_format(self):
        result, err = PhoneNormalizer.normalize("50936305123")
        self.assertIsNone(err)
        self.assertEqual(result, "+50936305123")

    def test_phone_invalid_prefix(self):
        result, err = PhoneNormalizer.normalize("96305123")
        self.assertIsNotNone(err)

    def test_phone_non_haitian_code(self):
        result, err = PhoneNormalizer.normalize("+33612345678")
        self.assertIsNotNone(err)

    def test_phone_empty(self):
        result, err = PhoneNormalizer.normalize("")
        self.assertIsNone(result)
        self.assertIsNotNone(err)

    def test_phone_none(self):
        result, err = PhoneNormalizer.normalize(None)
        self.assertIsNone(result)
        self.assertIsNotNone(err)

    def test_phone_display_format(self):
        result = PhoneNormalizer.format_display("36305123")
        self.assertEqual(result, "+509 363-05-123")


class TestAddressNormalization(unittest.TestCase):
    def test_address_simple(self):
        result, err = AddressNormalizer.normalize("Rue 25 Delmas")
        self.assertIsNone(err)
        self.assertIn("Delmas", result)

    def test_address_with_department_alias(self):
        result, err = AddressNormalizer.normalize("12 Rue Centre, OUEST")
        self.assertIsNone(err)
        self.assertIn("OUEST", result)

    def test_address_with_commune_alias(self):
        result, err = AddressNormalizer.normalize("P-AU-P")
        self.assertIsNone(err)
        self.assertIn("Port-Au-Prince", result)

    def test_address_with_creole_commune(self):
        result, err = AddressNormalizer.normalize("Okap")
        self.assertIsNone(err)
        self.assertIn("Cap-Haitien", result)

    def test_address_empty(self):
        result, err = AddressNormalizer.normalize("")
        self.assertIsNone(result)
        self.assertIsNotNone(err)

    def test_address_none(self):
        result, err = AddressNormalizer.normalize(None)
        self.assertIsNone(result)
        self.assertIsNotNone(err)


class TestDataCleansingEngine(unittest.TestCase):
    def setUp(self):
        self.engine = DataCleansingEngine()

    def test_cleanse_basic_record(self):
        record = {
            "nom": "Dupont",
            "prenom": "Jean",
            "birth_date": "15/03/1990",
            "phone": "36305123",
            "address": "Rue 25 Delmas",
            "cin": "001-123-456-78",
        }
        cleaned, errors = self.engine.cleanse(record)
        self.assertEqual(len(errors), 0)
        self.assertEqual(cleaned["nom"], "DUPONT")
        self.assertEqual(cleaned["birth_date"], "1990-03-15")
        self.assertEqual(cleaned["phone"], "+50936305123")

    def test_cleanse_empty_record(self):
        record = {}
        cleaned, errors = self.engine.cleanse(record)
        self.assertIsInstance(cleaned, dict)

    def test_cleanse_with_disabled_rules(self):
        config = CleansingRules(
            normalize_names=False,
            normalize_dates=False,
            normalize_phones=False,
            normalize_addresses=False,
        )
        engine = DataCleansingEngine(config)
        record = {"nom": "Dupont", "birth_date": "15/03/1990", "phone": "36305123"}
        cleaned, errors = engine.cleanse(record)
        self.assertEqual(cleaned["nom"], "Dupont")
        self.assertEqual(cleaned["birth_date"], "15/03/1990")

    def test_cleanse_with_none_values(self):
        record = {"nom": None, "prenom": None}
        cleaned, errors = self.engine.cleanse(record)
        self.assertEqual(cleaned["nom"], "UNKNOWN")

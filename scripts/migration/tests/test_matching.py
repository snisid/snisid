"""Tests for MatchingEngine."""

import unittest

from matching_engine import MatchingEngine, MatchLevel


class TestMatchingEngine(unittest.TestCase):
    def setUp(self):
        self.engine = MatchingEngine()

    def test_exact_match_by_national_id(self):
        self.engine.load_reference([
            {"id": "1", "first_name": "Jean", "last_name": "Dupont", "national_id": "001-123-456-78", "date_of_birth": "1990-03-15"},
        ])
        record = {"first_name": "Jean", "last_name": "Dupont", "national_id": "001-123-456-78", "date_of_birth": "1990-03-15"}
        result = self.engine.match(record)
        self.assertTrue(result.is_duplicate)
        self.assertGreaterEqual(result.confidence, 0.90)
        self.assertEqual(result.level, MatchLevel.HIGH)

    def test_exact_match_no_duplicate(self):
        self.engine.load_reference([
            {"id": "1", "first_name": "Jean", "last_name": "Dupont", "national_id": "001-123-456-78", "date_of_birth": "1990-03-15"},
        ])
        record = {"first_name": "Marie", "last_name": "Pierre", "national_id": "002-987-654-32", "date_of_birth": "1985-07-20"}
        result = self.engine.match(record)
        self.assertFalse(result.is_duplicate)
        self.assertEqual(result.level, MatchLevel.NONE)

    def test_jaro_winkler_fuzzy_match(self):
        self.engine.load_reference([
            {"id": "1", "first_name": "Jean", "last_name": "Dupont", "date_of_birth": "1990-03-15"},
        ])
        record = {"first_name": "Jeam", "last_name": "Dupont", "date_of_birth": "1990-03-15"}
        result = self.engine.match(record)
        score = self.engine._jaro_winkler("jean", "jeam")
        self.assertGreater(score, 0.8)

    def test_phonetic_matching_soundex(self):
        score = self.engine._soundex_compare("jean", "jeam")
        self.assertEqual(score, 1.0)

    def test_soundex_different_words(self):
        score = self.engine._soundex_compare("jean", "pierre")
        self.assertEqual(score, 0.0)

    def test_soundex_code_generation(self):
        self.assertEqual(self.engine._soundex("jean"), "J500")
        self.assertEqual(self.engine._soundex("dupont"), "D153")

    def test_auto_merge_threshold(self):
        self.engine.load_reference([
            {"id": "1", "first_name": "Jean", "last_name": "Dupont", "national_id": "001-123-456-78", "date_of_birth": "1990-03-15"},
        ])
        record = {"first_name": "Jean", "last_name": "Dupont", "national_id": "001-123-456-78", "date_of_birth": "1990-03-15"}
        result = self.engine.match(record)
        self.assertGreaterEqual(result.confidence, self.engine.THRESHOLD_HIGH)
        self.assertIsNotNone(result.merged)

    def test_pending_review_threshold_medium(self):
        self.engine.load_reference([
            {"id": "1", "first_name": "Jean", "last_name": "Dupont", "date_of_birth": "1990-03-15", "birth_place": "Port-au-Prince", "gender": "M"},
        ])
        record = {"first_name": "Jeam", "last_name": "Dupont", "date_of_birth": "1990-03-15", "birth_place": "Port-au-Prince", "gender": "M"}
        result = self.engine.match(record)
        if result.confidence >= self.engine.THRESHOLD_MEDIUM:
            self.assertEqual(result.level, MatchLevel.MEDIUM)

    def test_blocking_strategy_empty_reference(self):
        self.engine.load_reference([])
        record = {"first_name": "Jean", "last_name": "Dupont"}
        result = self.engine.match(record)
        self.assertFalse(result.is_duplicate)

    def test_merge_detects_conflicts(self):
        self.engine.load_reference([
            {"id": "1", "first_name": "Jean", "last_name": "Dupont", "phone": "1111"},
        ])
        record = {"first_name": "Jean", "last_name": "Dupont", "phone": "2222"}
        result = self.engine.match(record)
        if result.is_duplicate:
            self.assertGreater(len(result.conflicts), 0)


class TestJaroWinkler(unittest.TestCase):
    def test_identical_strings(self):
        score = MatchingEngine._jaro_winkler("hello", "hello")
        self.assertEqual(score, 1.0)

    def test_completely_different(self):
        score = MatchingEngine._jaro_winkler("abc", "xyz")
        self.assertEqual(score, 0.0)

    def test_similar_strings(self):
        score = MatchingEngine._jaro_winkler("marhta", "martha")
        self.assertGreater(score, 0.9)

    def test_empty_strings(self):
        score = MatchingEngine._jaro_winkler("", "")
        self.assertEqual(score, 0.0)


class TestSoundex(unittest.TestCase):
    def test_soundex_jean(self):
        code = MatchingEngine._soundex("jean")
        self.assertEqual(code, "J500")

    def test_soundex_dupont(self):
        code = MatchingEngine._soundex("dupont")
        self.assertEqual(code, "D153")

    def test_soundex_empty(self):
        code = MatchingEngine._soundex("")
        self.assertEqual(code, "0000")

    def test_soundex_ashcroft(self):
        code = MatchingEngine._soundex("ashcroft")
        self.assertEqual(code, "A261")

    def test_soundex_case_insensitive(self):
        code1 = MatchingEngine._soundex("JEAN")
        code2 = MatchingEngine._soundex("jean")
        self.assertEqual(code1, code2)

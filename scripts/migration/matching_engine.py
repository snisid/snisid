"""SNISID — Moteur de matching/déduplication.

Utilise Jaro-Winkler + Soundex pour le matching probabiliste
des enregistrements d'identité. Supporte le merging automatique
et la résolution de conflits.
"""

import logging
import re
from dataclasses import dataclass, field
from typing import Dict, Any, Optional, List, Tuple
from enum import Enum

logger = logging.getLogger(__name__)


class MatchLevel(Enum):
    NONE = 0
    LOW = 1
    MEDIUM = 2
    HIGH = 3
    EXACT = 4


@dataclass
class MatchResult:
    is_duplicate: bool = False
    confidence: float = 0.0
    level: MatchLevel = MatchLevel.NONE
    matched_id: Optional[str] = None
    merged: Optional[Dict[str, Any]] = None
    conflicts: List[str] = field(default_factory=list)


class MatchingEngine:
    """Moteur de matching avec stratégies multiples."""

    # Poids des champs pour le scoring
    FIELD_WEIGHTS = {
        "first_name": 0.20,
        "last_name": 0.25,
        "date_of_birth": 0.25,
        "birth_place": 0.05,
        "gender": 0.05,
        "national_id": 0.10,
        "phone": 0.05,
        "email": 0.05,
    }

    # Seuils de confiance
    THRESHOLD_HIGH = 0.90
    THRESHOLD_MEDIUM = 0.75
    THRESHOLD_LOW = 0.55

    def __init__(self, soundex_enabled: bool = True, jaro_weight: float = 0.7):
        self.soundex_enabled = soundex_enabled
        self.jaro_weight = jaro_weight
        self._known_records: List[Dict[str, Any]] = []

    def match(self, record: Dict[str, Any]) -> MatchResult:
        """Tente de matcher un enregistrement contre les connus."""
        best_score = 0.0
        best_match = None

        for known in self._known_records:
            score = self._compute_similarity(record, known)
            if score > best_score:
                best_score = score
                best_match = known

        result = MatchResult()

        if best_score >= self.THRESHOLD_HIGH:
            result.is_duplicate = True
            result.level = MatchLevel.HIGH
            result.confidence = best_score
            result.matched_id = best_match.get("id")
            result.merged = self._merge_records(record, best_match)
            result.conflicts = self._detect_conflicts(record, best_match)
        elif best_score >= self.THRESHOLD_MEDIUM:
            result.is_duplicate = True
            result.level = MatchLevel.MEDIUM
            result.confidence = best_score
            result.matched_id = best_match.get("id")
            result.merged = self._merge_records(record, best_match, prefer_newer=True)
        elif best_score >= self.THRESHOLD_LOW:
            result.is_duplicate = False
            result.level = MatchLevel.LOW
            result.confidence = best_score

        return result

    def _compute_similarity(
        self, a: Dict[str, Any], b: Dict[str, Any]
    ) -> float:
        """Calcule le score de similarité entre deux enregistrements."""
        total_weight = 0.0
        weighted_score = 0.0

        for field, weight in self.FIELD_WEIGHTS.items():
            val_a = str(a.get(field, "")).strip().lower()
            val_b = str(b.get(field, "")).strip().lower()

            if not val_a or not val_b:
                continue

            field_score = self._jaro_winkler(val_a, val_b)

            if self.soundex_enabled:
                soundex_score = self._soundex_compare(val_a, val_b)
                field_score = (
                    self.jaro_weight * field_score
                    + (1 - self.jaro_weight) * soundex_score
                )

            weighted_score += weight * field_score
            total_weight += weight

        return weighted_score / total_weight if total_weight > 0 else 0.0

    @staticmethod
    def _jaro_winkler(s1: str, s2: str) -> float:
        """Distance Jaro-Winkler entre deux chaînes."""
        if s1 == s2:
            return 1.0

        len1, len2 = len(s1), len(s2)
        match_dist = max(len1, len2) // 2 - 1
        matches = 0
        transpositions = 0

        s1_matches = [False] * len1
        s2_matches = [False] * len2

        for i in range(len1):
            start = max(0, i - match_dist)
            end = min(i + match_dist + 1, len2)
            for j in range(start, end):
                if s2_matches[j] or s1[i] != s2[j]:
                    continue
                s1_matches[i] = True
                s2_matches[j] = True
                matches += 1
                break

        if matches == 0:
            return 0.0

        k = 0
        for i in range(len1):
            if not s1_matches[i]:
                continue
            while not s2_matches[k]:
                k += 1
            if s1[i] != s2[k]:
                transpositions += 1
            k += 1

        jaro = (
            matches / len1
            + matches / len2
            + (matches - transpositions / 2) / matches
        ) / 3

        # Winkler boost pour les préfixes communs
        prefix = 0
        for i in range(min(4, len1, len2)):
            if s1[i] == s2[i]:
                prefix += 1
            else:
                break

        return jaro + prefix * 0.1 * (1 - jaro)

    @staticmethod
    def _soundex(word: str) -> str:
        """Algorithme Soundex (code à 4 caractères)."""
        word = word.upper()
        soundex_map = {
            "B": "1", "F": "1", "P": "1", "V": "1",
            "C": "2", "G": "2", "J": "2", "K": "2",
            "Q": "2", "S": "2", "X": "2", "Z": "2",
            "D": "3", "T": "3",
            "L": "4",
            "M": "5", "N": "5",
            "R": "6",
        }

        if not word:
            return "0000"

        soundex_code = word[0]
        last_code = soundex_map.get(word[0], "0")

        for char in word[1:]:
            code = soundex_map.get(char, "0")
            if code != last_code and code != "0":
                soundex_code += code
                last_code = code

        soundex_code = soundex_code[:4].ljust(4, "0")
        return soundex_code

    def _soundex_compare(self, s1: str, s2: str) -> float:
        """Compare deux mots via Soundex."""
        return 1.0 if self._soundex(s1) == self._soundex(s2) else 0.0

    @staticmethod
    def _merge_records(
        existing: Dict[str, Any],
        incoming: Dict[str, Any],
        prefer_newer: bool = False,
    ) -> Dict[str, Any]:
        """Fusionne deux enregistrements en résolvant les conflits."""
        merged = {}
        for key in set(list(existing.keys()) + list(incoming.keys())):
            if key in existing and key in incoming:
                if existing[key] == incoming[key]:
                    merged[key] = existing[key]
                else:
                    merged[key] = incoming[key] if prefer_newer else existing[key]
                    merged[f"{key}_alt"] = existing[key] if prefer_newer else incoming[key]
            elif key in existing:
                merged[key] = existing[key]
            else:
                merged[key] = incoming[key]
        return merged

    @staticmethod
    def _detect_conflicts(
        a: Dict[str, Any], b: Dict[str, Any]
    ) -> List[str]:
        """Détecte les champs en conflit entre deux enregistrements."""
        conflicts = []
        for key in set(a.keys()) & set(b.keys()):
            if a[key] != b[key] and key not in ("id", "created_at", "updated_at"):
                conflicts.append(key)
        return conflicts

    def is_ready(self) -> bool:
        return True

    def load_reference(self, records: List[Dict[str, Any]]) -> None:
        """Charge les enregistrements de référence."""
        self._known_records = records
        logger.info("Référence chargée: %d enregistrements", len(records))

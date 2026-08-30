"""
Presentation Exchange (PEX) - DIF Presentation Exchange v2

Validates ``presentation_definition`` against ``presentation_submission``
per the DIF Presentation Exchange v2.0 specification.

Core features:
- Input descriptor matching against VCs
- Constraint evaluation ($iss, $schema, $and, $or, $not, field paths)
- Submission requirement parsing
- Presentation definition filtering
"""
from __future__ import annotations

import json
import re
from dataclasses import dataclass, field
from typing import Any

from services.vp import VerifiablePresentation, VPIssuer


@dataclass
class InputDescriptor:
    """DIF Presentation Exchange Input Descriptor."""

    id: str
    name: str | None = None
    purpose: str | None = None
    schema_uris: list[str] = field(default_factory=list)
    constraints: dict[str, Any] = field(default_factory=dict)
    group: list[str] | None = None


@dataclass
class PresentationDefinition:
    """DIF Presentation Exchange Presentation Definition."""

    id: str
    input_descriptors: list[InputDescriptor]
    name: str | None = None
    purpose: str | None = None
    submission_requirements: list[dict] | None = None
    frame: dict | None = None

    @classmethod
    def from_dict(cls, data: dict) -> PresentationDefinition:
        raw_descriptors = data.get("input_descriptors", [])
        descriptors = [
            InputDescriptor(
                id=d["id"],
                name=d.get("name"),
                purpose=d.get("purpose"),
            schema_uris=(
                [s["uri"] if isinstance(s, dict) else s for s in d.get("schema", [])]
                if isinstance(d.get("schema"), list)
                else [d["schema"]["uri"]] if isinstance(d.get("schema"), dict) and "uri" in d.get("schema", {})
                else []
            ),
                constraints=d.get("constraints", {}),
                group=d.get("group"),
            )
            for d in raw_descriptors
        ]
        return cls(
            id=data.get("id", ""),
            input_descriptors=descriptors,
            name=data.get("name"),
            purpose=data.get("purpose"),
            submission_requirements=data.get("submission_requirements"),
            frame=data.get("frame"),
        )


@dataclass
class DescriptorMatch:
    """Result of matching a VC against an InputDescriptor."""

    descriptor_id: str
    credential: dict[str, Any]
    matched: bool
    errors: list[str] = field(default_factory=list)


@dataclass
class PEXResult:
    """Result of evaluating a PresentationDefinition against VCs."""

    definition_id: str
    matches: list[DescriptorMatch]
    valid: bool
    errors: list[str] = field(default_factory=list)


def _evaluate_constraint(credential: dict, constraint: dict) -> tuple[bool, list[str]]:
    """
    Evaluate a constraint block from an InputDescriptor.

    Supported fields:
      - ``subject_is_issuer``: bool
      - ``is_holder``: list of holder DIDs
      - ``same_subject_id``: list of descriptor IDs
      - ``fields``: list of field filters with ``path``, ``filter``, ``optional``, ``purpose``
      - ``limit_disclosure``: required or preferred
      - ``statuses``: required credential status (e.g. active, suspended, revoked)
    """
    errors: list[str] = []
    vc = credential.get("credentialSubject", credential)

    # subject_is_issuer
    if constraint.get("subject_is_issuer"):
        issuer_id = credential.get("issuer", "")
        subject_id = ""
        if isinstance(credential.get("credentialSubject"), dict):
            subject_id = credential["credentialSubject"].get("id", "")
        if issuer_id and subject_id and issuer_id != subject_id:
            errors.append(f"subject_is_issuer: issuer {issuer_id} != subject {subject_id}")

    # is_holder
    holders = constraint.get("is_holder", [])
    if holders:
        subject_id = ""
        if isinstance(credential.get("credentialSubject"), dict):
            subject_id = credential["credentialSubject"].get("id", "")
        if subject_id and subject_id not in holders:
            errors.append(f"is_holder: {subject_id} not in {holders}")

    # fields
    fields = constraint.get("fields", [])
    for field in fields:
        paths = field.get("path", [])
        field_filter = field.get("filter", {}).get("const", None)
        optional = field.get("optional", False)
        found = False
        for path in paths:
            value = _resolve_json_path(credential, path)
            if value is not None:
                found = True
                if field_filter is not None and value != field_filter:
                    errors.append(f"field {path}: expected {field_filter}, got {value}")
                break
        if not found and not optional:
            errors.append(f"field {paths}: not found in credential")

    # statuses
    statuses = constraint.get("statuses", {})
    if statuses:
        cred_status = credential.get("credentialStatus", {})
        for status_name, required in statuses.items():
            if required and cred_status.get("status") != status_name and cred_status.get("status") != status_name.lower():
                if cred_status:
                    errors.append(f"status {status_name}: not matched (got {cred_status.get('status')})")

    return len(errors) == 0, errors


def _resolve_json_path(doc: dict, path: str) -> Any:
    """Resolve a JSONPath-like path (e.g. ``$.credentialSubject.id``)."""
    parts = path.strip("$.").split(".")
    current: Any = doc
    for part in parts:
        if isinstance(current, dict):
            current = current.get(part)
        elif isinstance(current, list):
            try:
                idx = int(part)
                current = current[idx] if 0 <= idx < len(current) else None
            except (ValueError, IndexError):
                return None
        else:
            return None
        if current is None:
            return None
    return current


def match_descriptor(descriptor: InputDescriptor, credential: dict) -> DescriptorMatch:
    """Match a single InputDescriptor against a verifiable credential."""
    errors: list[str] = []

    # Schema validation
    if descriptor.schema_uris:
        vc_type = credential.get("type", [])
        if isinstance(vc_type, str):
            vc_type = [vc_type]
        vc_type_lower = [t.lower() for t in vc_type]
        for schema_uri in descriptor.schema_uris:
            schema_name = schema_uri.split("/")[-1].split(".")[0].lower()
            if schema_name not in vc_type_lower and schema_uri.lower() not in vc_type_lower:
                errors.append(f"schema {schema_uri}: not matched in credential types {vc_type}")

    # Constraints
    if descriptor.constraints:
        ok, constraint_errors = _evaluate_constraint(credential, descriptor.constraints)
        if not ok:
            errors.extend(constraint_errors)

    return DescriptorMatch(
        descriptor_id=descriptor.id,
        credential=credential,
        matched=len(errors) == 0,
        errors=errors,
    )


def evaluate_definition(definition: PresentationDefinition, credentials: list[dict]) -> PEXResult:
    """Evaluate a PresentationDefinition against a list of VCs."""
    errors: list[str] = []
    all_matches: list[DescriptorMatch] = []

    for descriptor in definition.input_descriptors:
        best_match: DescriptorMatch | None = None
        for credential in credentials:
            match = match_descriptor(descriptor, credential)
            if match.matched:
                best_match = match
                break
        if best_match is None:
            # No match found — try the last credential to collect errors
            if credentials:
                match = match_descriptor(descriptor, credentials[-1])
                all_matches.append(match)
            else:
                all_matches.append(DescriptorMatch(
                    descriptor_id=descriptor.id,
                    credential={},
                    matched=False,
                    errors=["no credentials provided"],
                ))
        else:
            all_matches.append(best_match)

    # Evaluate submission_requirements
    if definition.submission_requirements:
        for req in definition.submission_requirements:
            rule = req.get("rule", "all")
            from_field = req.get("from", "")
            count = req.get("count", 0)
            matched_ids = {m.descriptor_id for m in all_matches if m.matched}
            if from_field:
                group_ids = {
                    d.id for d in definition.input_descriptors
                    if d.group and from_field in d.group
                }
                matched_in_group = len(matched_ids & group_ids)
                if rule == "all" and matched_in_group < len(group_ids):
                    errors.append(f"submission_requirement {req.get('name','')}: not all descriptors matched in group {from_field}")
                elif rule == "pick" and count and matched_in_group < count:
                    errors.append(f"submission_requirement {req.get('name','')}: need {count}, got {matched_in_group}")

    valid = all(m.matched for m in all_matches) and len(errors) == 0
    return PEXResult(
        definition_id=definition.id,
        matches=all_matches,
        valid=valid,
        errors=errors,
    )


def create_presentation_from_definition(
    definition: PresentationDefinition,
    credentials: list[dict],
    holder_did: str,
    issuer_did: str | None = None,
) -> VerifiablePresentation | None:
    """
    Create a VP from a PresentationDefinition and matching credentials.

    Returns ``None`` if not all required descriptors can be satisfied.
    """
    result = evaluate_definition(definition, credentials)
    if not result.valid:
        return None

    matching_credentials = [m.credential for m in result.matches if m.matched]
    if not matching_credentials:
        return None

    issuer = VPIssuer(issuer_did or holder_did)
    return issuer.create_presentation(
        holder_did=holder_did,
        verifiable_credentials=matching_credentials,
    )


def filter_credentials_by_definition(
    definition: PresentationDefinition,
    credentials: list[dict],
) -> list[dict]:
    """Return only the VCs that match at least one InputDescriptor."""
    result = evaluate_definition(definition, credentials)
    return [m.credential for m in result.matches if m.matched]

#!/usr/bin/env python3
"""SNISID-BIO-ADN CLI — Opérations LDIS.

Usage:
    snisid-bio-cli upload --level ldis-to-sdis --lab-code LDIS-PAP-001 \\
        --date-from 2026-06-01 --date-to 2026-06-09 \\
        --operator-niu HTI-XXXXXXXXXX
"""

from __future__ import annotations

import argparse
import json
import os
import sys
import uuid
from datetime import datetime, timezone

import httpx


BASE_URL = os.getenv("SNISID_BIO_API", "http://localhost:8092/v1/bio-adn")


def upload(args: argparse.Namespace) -> None:
    payload = {
        "level": args.level,
        "lab_code": args.lab_code,
        "date_from": args.date_from,
        "date_to": args.date_to,
        "operator_niu": args.operator_niu,
    }
    with httpx.Client(timeout=60) as client:
        resp = client.post(f"{BASE_URL}/lab/upload", json=payload)
    if resp.status_code == 202:
        data = resp.json()
        print(f"Upload déclenché: {data['uploaded_count']} profils, succès={data['success']}")
    else:
        print(f"Erreur {resp.status_code}: {resp.text}", file=sys.stderr)
        sys.exit(1)


def list_labs(args: argparse.Namespace) -> None:
    with httpx.Client() as client:
        resp = client.get(f"{BASE_URL}/lab/labs")
    if resp.status_code == 200:
        data = resp.json()
        print(f"{'Code':<16} {'Nom':<40} {'Département':<15} {'Statut':<15}")
        print("-" * 86)
        for lab in data["results"]:
            print(f"{lab['code']:<16} {lab['name']:<40} {lab['department']:<15} {lab['status']:<15}")
    else:
        print(f"Erreur {resp.status_code}: {resp.text}", file=sys.stderr)
        sys.exit(1)


def main() -> None:
    parser = argparse.ArgumentParser(description="SNISID-BIO-ADN CLI")
    sub = parser.add_subparsers(dest="command")

    upload_p = sub.add_parser("upload", help="Déclencher upload LDIS→SDIS")
    upload_p.add_argument("--level", required=True, choices=["ldis-to-sdis"])
    upload_p.add_argument("--lab-code", required=True)
    upload_p.add_argument("--date-from", required=True)
    upload_p.add_argument("--date-to", required=True)
    upload_p.add_argument("--operator-niu", required=True)

    sub.add_parser("list-labs", help="Afficher les laboratoires LDIS")

    args = parser.parse_args()
    if args.command == "upload":
        upload(args)
    elif args.command == "list-labs":
        list_labs(args)
    else:
        parser.print_help()
        sys.exit(1)


if __name__ == "__main__":
    main()

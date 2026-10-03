#!/usr/bin/env python3
"""Generate vercel.json (Vercel multi-services config) from the repository layout.

Detection rules:
  - FastAPI service : root has main.py + requirements.txt and imports fastapi
  - Go service      : root has go.mod AND a main package (cmd/ or *.go at root)
  - Vite frontend   : root has package.json with vite dependency
  - Container svc   : root has Dockerfile but none of the above applies on Vercel
All other services/* dirs are emitted as "container" placeholders pending review.
"""
import json
import os
import re
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

# Directories that are NOT deployable (docs, drafts, libs, tooling)
EXCLUDE = {
    "snisid-afis-ht/internal_legacy_draft",
}


def detect_kind(root_abs):
    files = set(os.listdir(root_abs))
    # Vite / node frontend
    if "package.json" in files:
        try:
            pkg = json.load(open(os.path.join(root_abs, "package.json")))
        except Exception:
            pkg = {}
        deps = {**pkg.get("dependencies", {}), **pkg.get("devDependencies", {})}
        if "vite" in deps:
            return ("vite", {"framework": "vite"})
        if "next" in deps:
            return ("next", {"framework": "nextjs"})
        return ("node", {"runtime": "node"})
    # Python FastAPI
    if "requirements.txt" in files or "pyproject.toml" in files:
        entry = None
        for cand in ("main.py", "app/main.py", "api/main.py"):
            if os.path.exists(os.path.join(root_abs, cand)):
                entry = cand
                break
        is_fastapi = False
        pyfiles = [f for f in files if f.endswith(".py")]
        for f in pyfiles[:20]:
            try:
                txt = open(os.path.join(root_abs, f), encoding="utf8", errors="ignore").read(4000)
                if "from fastapi import" in txt or "FastAPI(" in txt:
                    is_fastapi = True
                    break
            except Exception:
                pass
        if not is_fastapi and entry:
            try:
                txt = open(os.path.join(root_abs, entry), encoding="utf8", errors="ignore").read()
                is_fastapi = "fastapi" in txt.lower()
            except Exception:
                pass
        if is_fastapi:
            cfg = {"framework": "fastapi"}
            if entry and entry != "main.py":
                cfg["entrypoint"] = entry.replace("/", ".")[:-3]
            return ("fastapi", cfg)
        return ("python", {"runtime": "python"})
    # Go
    if "go.mod" in files:
        has_main = False
        for dirpath, dirnames, filenames in os.walk(root_abs):
            if "/." in dirpath or "vendor" in dirpath:
                continue
            for fn in filenames:
                if fn.endswith(".go"):
                    try:
                        head = open(os.path.join(dirpath, fn), encoding="utf8", errors="ignore").read(500)
                        if re.search(r"^package main\b", head, re.M):
                            has_main = True
                            break
                    except Exception:
                        pass
            if has_main:
                break
        if has_main:
            return ("go", {"runtime": "go"})
    if "Dockerfile" in files:
        return ("container", {"runtime": "container"})
    return (None, None)


def name_for(rel):
    n = rel.replace("/", "-").replace(" ", "-").strip("-").lower()
    n = re.sub(r"[^a-z0-9._-]", "-", n)
    n = re.sub(r"-{2,}", "-", n)
    return n[:100]


def main():
    services = {}
    report = []

    candidates = []
    # top-level app roots
    for rel in [
        ".", "backend", "frontend", "scripts/migration",
        "National-Executive-Operations/api", "National-Executive-Operations/frontend",
        "National-Sustainability-Billing", "snisid-ui", "snisid mcp",
        "SNI-SIDE/etl", "SNI-SIDE/sdk/python", "SNI-SIDE/services",
        "SNI-SIDE/services/event_processor", "SNI-SIDE/services/graphrag_engine",
        "SNI-SIDE/web",
    ]:
        candidates.append(rel)
    # services/*
    sdir = os.path.join(ROOT, "services")
    if os.path.isdir(sdir):
        for d in sorted(os.listdir(sdir)):
            p = os.path.join(sdir, d)
            if os.path.isdir(p) and not d.startswith("."):
                candidates.append(f"services/{d}")
    # src/*
    for sub in ["src/frontend", "src/services", "src/web"]:
        p = os.path.join(ROOT, sub)
        if os.path.isdir(p):
            for d in sorted(os.listdir(p)):
                if os.path.isdir(os.path.join(p, d)):
                    candidates.append(f"{sub}/{d}")

    for rel in candidates:
        if rel in EXCLUDE:
            continue
        abs_path = os.path.join(ROOT, rel)
        if not os.path.isdir(abs_path):
            report.append(f"MISSING  {rel}")
            continue
        kind, cfg = detect_kind(abs_path)
        if kind is None:
            report.append(f"SKIP     {rel} (no deployable marker)")
            continue
        nm = name_for(rel) if rel != "." else "app"
        if nm in services:
            nm = nm + "-" + str(len(services))
        services[nm] = {"root": rel, **cfg}
        report.append(f"{kind.upper():9s} {nm:40s} <- {rel}")

    # --- Merge existing vercel.json if present (preserves manual edits:
    # bindings, public rewrites, per-service overrides such as entrypoint). ---
    vpath = os.path.join(ROOT, "vercel.json")
    existing = {}
    if os.path.exists(vpath):
        try:
            existing = json.load(open(vpath))
        except Exception:
            existing = {}

    for nm, cfg in existing.get("services", {}).items():
        if nm not in services:
            report.append(f"KEPT      {nm} (manual service from existing vercel.json)")
            services[nm] = cfg
        else:
            # manual overrides win for keys the generator cannot know about
            for k in ("bindings", "entrypoint", "buildCommand", "devCommand"):
                if k in cfg:
                    services[nm][k] = cfg[k]
            # allow explicit runtime/framework override
            for k in ("runtime", "framework"):
                if k in cfg and cfg[k] != services[nm].get(k):
                    services[nm][k] = cfg[k]

    rewrites = existing.get("rewrites") or [
        {"source": "/intelligence/(.*)", "destination": {"service": "sni-side-services"}},
        {"source": "/neo/api/(.*)", "destination": {"service": "national-executive-operations-api"}},
        {"source": "/billing/(.*)", "destination": {"service": "national-sustainability-billing"}},
        {"source": "/migration/(.*)", "destination": {"service": "scripts-migration"}},
        {"source": "/mcp/(.*)", "destination": {"service": "snisid-mcp"}},
        {"source": "/identity/(.*)", "destination": {"service": "src-services-identity-service"}},
        {"source": "/registry/(.*)", "destination": {"service": "src-services-citizen-registry-service"}},
        {"source": "/consent/(.*)", "destination": {"service": "src-services-consent-service"}},
        {"source": "/api/(.*)", "destination": {"service": "backend"}},
        {"source": "/(.*)", "destination": {"service": "frontend"}},
    ]

    out = {
        "$schema": "https://openapi.vercel.sh/vercel.json",
        "services": services,
        "rewrites": rewrites,
    }
    with open(vpath, "w") as f:
        json.dump(out, f, indent=2)
        f.write("\n")
    print("\n".join(report))
    print(f"\nTOTAL SERVICES: {len(services)}")


if __name__ == "__main__":
    sys.exit(main())

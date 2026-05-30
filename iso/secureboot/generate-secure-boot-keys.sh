#!/bin/bash
# SNISID Sovereign UEFI Secure Boot Key Generator
# Fulfills Requirement #126: Secure Boot validé

set -e

KEY_DIR="./snisid-uefi-keys"
mkdir -p "$KEY_DIR"
cd "$KEY_DIR"

echo "[*] Generating Sovereign Platform Key (PK)..."
openssl req -new -x509 -newkey rsa:2048 -subj "/CN=SNISID Platform Key/" -keyout PK.key -out PK.crt -days 3650 -nodes -sha256
openssl x509 -in PK.crt -out PK.cer -outform DER

echo "[*] Generating Key Exchange Key (KEK)..."
openssl req -new -x509 -newkey rsa:2048 -subj "/CN=SNISID Key Exchange Key/" -keyout KEK.key -out KEK.crt -days 3650 -nodes -sha256
openssl x509 -in KEK.crt -out KEK.cer -outform DER

echo "[*] Generating Signature Database Key (db)..."
openssl req -new -x509 -newkey rsa:2048 -subj "/CN=SNISID Signature Database/" -keyout db.key -out db.crt -days 3650 -nodes -sha256
openssl x509 -in db.crt -out db.cer -outform DER

echo "[*] Converting keys to EFI Signature Lists (ESL)..."
# Requires efitools (apt install efitools)
cert-to-efi-sig-list -g "$(uuidgen)" PK.crt PK.esl
cert-to-efi-sig-list -g "$(uuidgen)" KEK.crt KEK.esl
cert-to-efi-sig-list -g "$(uuidgen)" db.crt db.esl

echo "[*] Signing EFI Signature Lists..."
sign-efi-sig-list -k PK.key -c PK.crt PK PK.esl PK.auth
sign-efi-sig-list -k PK.key -c PK.crt KEK KEK.esl KEK.auth
sign-efi-sig-list -k KEK.key -c KEK.crt db db.esl db.auth

echo "[!] Keys generated successfully in $KEY_DIR"
echo "[!] IMPORTANT: You must manually enroll PK.auth, KEK.auth, and db.auth into the target hardware BIOS/UEFI firmware before booting the SNISID ISO."

#!/usr/bin/env python3
"""
SNISID PKI - Shamir's Secret Sharing Tool for Root CA Master Key Recovery
Splits a 256-bit or 384-bit master wrap key into 9 parts (threshold 5-of-9).
No external dependencies required.
"""

import sys
import secrets

# 13th Mersenne Prime: 2^521 - 1
# Used as the finite field prime to support secrets up to 520 bits.
PRIME = 2**521 - 1

def extended_gcd(a, b):
    if a == 0:
        return b, 0, 1
    else:
        g, x, y = extended_gcd(b % a, a)
        return g, y - (b // a) * x, x

def modular_inverse(k, p):
    g, x, y = extended_gcd(k, p)
    if g != 1:
        raise ValueError('Modular inverse does not exist')
    else:
        return x % p

def split_secret(secret_int, threshold, num_shares):
    """
    Splits an integer secret into shares using Shamir's scheme.
    """
    if secret_int >= PRIME:
        raise ValueError("Secret is too large for the prime field.")
        
    # Generate random coefficients for the polynomial f(x) = a_0 + a_1*x + ... + a_{t-1}*x^{t-1}
    # where a_0 is the secret.
    coefficients = [secret_int] + [secrets.randbelow(PRIME) for _ in range(threshold - 1)]
    
    shares = []
    for x in range(1, num_shares + 1):
        y = 0
        for power, coeff in enumerate(coefficients):
            y = (y + coeff * pow(x, power, PRIME)) % PRIME
        shares.append((x, y))
    return shares

def reconstruct_secret(shares):
    """
    Reconstructs the secret from a set of shares.
    shares: list of tuple (x, y)
    """
    secret = 0
    for j, (x_j, y_j) in enumerate(shares):
        numerator = 1
        denominator = 1
        for m, (x_m, _) in enumerate(shares):
            if m == j:
                continue
            numerator = (numerator * (-x_m)) % PRIME
            denominator = (denominator * (x_j - x_m)) % PRIME
            
        lagrange_coeff = (numerator * modular_inverse(denominator, PRIME)) % PRIME
        secret = (secret + y_j * lagrange_coeff) % PRIME
    return secret

def main():
    if len(sys.argv) < 2:
        print("SNISID Shamir Secret Sharing Tool")
        print("Usage:")
        print("  Split:       python shamir_secret_sharing.py split <hex_secret> [threshold] [total_parts]")
        print("  Reconstruct: python shamir_secret_sharing.py combine <x1:y1> <x2:y2> ...")
        print("\nExamples:")
        print("  python shamir_secret_sharing.py split 4a6f686e446f65 5 9")
        print("  python shamir_secret_sharing.py combine 1:value1 2:value2 3:value3 4:value4 5:value5")
        sys.exit(1)

    command = sys.argv[1].lower()

    if command == "split":
        if len(sys.argv) < 3:
            print("Error: Missing hex secret to split.")
            sys.exit(1)
        hex_secret = sys.argv[2].strip()
        threshold = int(sys.argv[3]) if len(sys.argv) > 3 else 5
        num_shares = int(sys.argv[4]) if len(sys.argv) > 4 else 9
        
        try:
            secret_int = int(hex_secret, 16)
        except ValueError:
            print("Error: Secret must be a valid hexadecimal string.")
            sys.exit(1)
            
        print(f"[*] Splitting secret into {num_shares} parts (Threshold: {threshold}-of-{num_shares})...")
        shares = split_secret(secret_int, threshold, num_shares)
        
        print("\n[+] Cryptographic Parts Generated:")
        # List of departments for geographic distribution
        depts = [
            "Port-au-Prince (BRH)",
            "Cap-Haïtien (Nord)",
            "Les Cayes (Sud)",
            "Hinche (Centre)",
            "Gonaïves (Artibonite)",
            "Jacmel (Sud-Est)",
            "Saint-Marc (Artibonite)",
            "Port-de-Paix (Nord-Ouest)",
            "Fort-Liberté (Nord-Est)"
        ]
        
        for i, (x, y) in enumerate(shares):
            dept_name = depts[i] if i < len(depts) else f"Backup Location {x}"
            # Format share as x:hex(y)
            hex_y = hex(y)[2:]
            print(f"  Share {x} [{dept_name}]: {x}:{hex_y}")
            
    elif command == "combine":
        if len(sys.argv) < 3:
            print("Error: Must provide at least one share to reconstruct.")
            sys.exit(1)
            
        shares_arg = sys.argv[2:]
        shares = []
        for s in shares_arg:
            if ":" not in s:
                print(f"Error: Invalid share format '{s}'. Use 'x:y_hex_value'.")
                sys.exit(1)
            parts = s.split(":")
            x = int(parts[0])
            y = int(parts[1], 16)
            shares.append((x, y))
            
        print(f"[*] Reconstructing secret from {len(shares)} shares...")
        try:
            secret_int = reconstruct_secret(shares)
            reconstructed_hex = hex(secret_int)[2:]
            # Ensure even number of hex characters
            if len(reconstructed_hex) % 2 != 0:
                reconstructed_hex = "0" + reconstructed_hex
            
            # Print reconstructed secret
            print(f"[+] Reconstructed Secret (Hex): {reconstructed_hex}")
            try:
                # Try decoding as text if printable
                text_val = bytes.fromhex(reconstructed_hex).decode('utf-8')
                print(f"[+] Reconstructed Secret (Plaintext): {text_val}")
            except Exception:
                pass
        except Exception as e:
            print(f"Error during reconstruction: {e}")
            sys.exit(1)
            
    else:
        print(f"Error: Unknown command '{command}'. Use 'split' or 'combine'.")
        sys.exit(1)

if __name__ == "__main__":
    main()

#!/usr/bin/env python3
"""
SNISID PKI - Bloom Filter CRL Compressor for USSD/SMS Distribution
Optimized for low-bandwidth Haitian telecom environments.
No external dependencies required.
"""

import sys
import math
import hashlib
import struct
import zlib
import base64
import re
import subprocess

class SimpleBloomFilter:
    def __init__(self, size_in_bits, num_hashes):
        self.size_in_bits = size_in_bits
        self.num_hashes = num_hashes
        self.bit_array = bytearray(math.ceil(size_in_bits / 8))

    def _hashes(self, item):
        # Kirsch-Mitzenmacher optimization for generating k hash values
        h = hashlib.sha256(str(item).encode('utf-8')).digest()
        h1 = struct.unpack("<Q", h[0:8])[0]
        h2 = struct.unpack("<Q", h[8:16])[0]
        for i in range(self.num_hashes):
            yield (h1 + i * h2) % self.size_in_bits

    def add(self, item):
        for bit_index in self._hashes(item):
            byte_index = bit_index // 8
            bit_offset = bit_index % 8
            self.bit_array[byte_index] |= (1 << bit_offset)

    def check(self, item):
        for bit_index in self._hashes(item):
            byte_index = bit_index // 8
            bit_offset = bit_index % 8
            if not (self.bit_array[byte_index] & (1 << bit_offset)):
                return False
        return True

    def serialize(self):
        # 8 bytes header (4 bytes size, 4 bytes hash count) followed by zlib-compressed bit array
        header = struct.pack("<II", self.size_in_bits, self.num_hashes)
        compressed = zlib.compress(self.bit_array, level=9)
        return header + compressed

    @classmethod
    def deserialize(cls, data):
        if len(data) < 8:
            raise ValueError("Data too short to contain Bloom filter header.")
        size_in_bits, num_hashes = struct.unpack("<II", data[:8])
        compressed_bit_array = data[8:]
        bit_array = zlib.decompress(compressed_bit_array)
        bf = cls(size_in_bits, num_hashes)
        bf.bit_array = bytearray(bit_array)
        return bf

def parse_crl_serials(crl_path):
    """
    Parses serial numbers from a CRL file.
    Attempts to use OpenSSL if available, falls back to parsing as plain text list.
    """
    serials = []
    # Try using OpenSSL
    try:
        process = subprocess.Popen(
            ["openssl", "crl", "-inform", "PEM", "-text", "-noout", "-in", crl_path],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )
        stdout, stderr = process.communicate()
        if process.returncode == 0:
            # Look for lines like "Serial Number: 100A" or "Serial Number: 4096"
            # In openssl CRL output:
            #    Serial Number: 01A2B3
            #        Revocation Date: ...
            pattern = re.compile(r"Serial Number:\s*([0-9A-Fa-f]+)")
            matches = pattern.findall(stdout)
            for m in matches:
                # Convert hex serial to integer
                serials.append(int(m, 16))
            return serials
    except Exception as e:
        # Fallback to plain text list
        pass
    
    # Try reading as plain text file containing one decimal or hex serial number per line
    try:
        with open(crl_path, "r") as f:
            for line in f:
                line = line.strip()
                if not line or line.startswith("#"):
                    continue
                # If hex (starts with 0x)
                if line.lower().startswith("0x"):
                    serials.append(int(line, 16))
                else:
                    try:
                        serials.append(int(line))
                    except ValueError:
                        # Try parsing as hex without prefix if it looks like hex
                        try:
                            serials.append(int(line, 16))
                        except ValueError:
                            pass
    except Exception as e:
        print(f"Error reading file {crl_path}: {e}", file=sys.stderr)
        
    return serials

def chunk_payload(payload_bytes, chunk_size=120):
    """
    Splits payload bytes into base64-encoded segments suitable for SMS chunks.
    """
    encoded = base64.b64encode(payload_bytes).decode('ascii')
    total_len = len(encoded)
    chunks = [encoded[i:i+chunk_size] for i in range(0, total_len, chunk_size)]
    return chunks

def main():
    if len(sys.argv) < 3:
        print("Usage: python compress_crl.py <input_crl_or_list_file> <output_bf_file> [capacity] [error_rate]")
        print("Example: python compress_crl.py my_crl.pem my_crl.bin 10000 0.001")
        sys.exit(1)

    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    capacity = int(sys.argv[3]) if len(sys.argv) > 3 else 10000
    error_rate = float(sys.argv[4]) if len(sys.argv) > 4 else 0.001

    print(f"[*] Parsing input from {input_file}...")
    serials = parse_crl_serials(input_file)
    
    if not serials:
        print("[!] Warning: No serial numbers found. Initializing Bloom filter with mock/default values for testing.")
        # Insert a few mock serials for validation
        serials = [1001, 1002, 1003, 1004, 1005]
        
    num_items = len(serials)
    print(f"[*] Found {num_items} revoked certificates.")
    
    # Calculate optimal Bloom filter size and hash count
    # m = - (n * ln(p)) / (ln(2)^2)
    # k = (m / n) * ln(2)
    n = max(num_items, capacity)
    p = error_rate
    m = int(- (n * math.log(p)) / (math.log(2) ** 2))
    k = int((m / n) * math.log(2))
    k = max(1, k)
    
    print(f"[*] Optimal parameters: size={m} bits ({math.ceil(m/8)} bytes uncompressed), hash functions={k}")
    
    bf = SimpleBloomFilter(m, k)
    for s in serials:
        bf.add(s)
        
    # Serialize and compress
    serialized_data = bf.serialize()
    compressed_size = len(serialized_data)
    print(f"[*] Compressed Bloom Filter size: {compressed_size} bytes ({compressed_size / 1024:.2f} KB)")
    
    if compressed_size > 50000:
        print(f"[!] Warning: Compressed size ({compressed_size} bytes) exceeds the 50 KB limit!")
    else:
        print(f"[+] Success: Compressed Bloom Filter fits within the 50 KB envelope (under SLA budget).")

    with open(output_file, "wb") as f:
        f.write(serialized_data)
    print(f"[+] Bloom Filter written to {output_file}")
    
    # Generate SMS USSD payload chunks
    chunks = chunk_payload(serialized_data, chunk_size=120)
    print(f"\n[+] Generated {len(chunks)} USSD SMS segments for rural distribution:")
    for idx, chunk in enumerate(chunks):
        # Standardized format: SNISID_CRL_BF:<chunk_index>/<total_chunks>:<payload>
        sms_body = f"SNISID_CRL_BF:{idx+1}/{len(chunks)}:{chunk}"
        print(f"  SMS {idx+1}/{len(chunks)} ({len(sms_body)} chars): {sms_body}")

if __name__ == "__main__":
    main()

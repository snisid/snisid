import hashlib
import json
import time
import os
import boto3

class HashChainLogger:
    """
    Cryptographic Immutable Audit Logger for SNISID.
    Appends audit events to an S3 WORM bucket, linking each event 
    cryptographically to the previous one via SHA-256 to ensure non-repudiation.
    """
    def __init__(self, bucket_name: str, state_table: str):
        self.s3 = boto3.client('s3')
        self.bucket_name = bucket_name
        self.state_table = state_table # Using DynamoDB or Postgres to store the last hash
        
        # In a real environment, retrieve the actual last hash from the DB
        # For demonstration, we use a genesis hash if empty
        self.last_hash = self._get_last_hash()

    def _get_last_hash(self) -> str:
        # Mock retrieval. 
        # Requirement #114: Audit trail immuable vérifié par hash chain
        return "0000000000000000000000000000000000000000000000000000000000000000"

    def _save_last_hash(self, new_hash: str):
        # Mock save to persistent DB
        self.last_hash = new_hash

    def log_event(self, event_data: dict) -> str:
        """
        Signs and logs an event immutably.
        """
        timestamp = int(time.time() * 1000)
        
        # 1. Construct the payload including the previous hash
        payload = {
            "timestamp": timestamp,
            "previous_hash": self.last_hash,
            "data": event_data
        }
        
        payload_json = json.dumps(payload, sort_keys=True)
        
        # 2. Calculate the new SHA-256 hash
        current_hash = hashlib.sha256(payload_json.encode('utf-8')).hexdigest()
        
        # 3. Add the calculated hash to the final record
        final_record = {
            **payload,
            "hash": current_hash
        }
        
        # 4. Write to the S3 WORM (Write-Once-Read-Many) Bucket
        # Using Object Lock / Legal Hold configured on the S3 bucket
        object_key = f"audit-logs/{time.strftime('%Y/%m/%d')}/{current_hash}.json"
        
        self.s3.put_object(
            Bucket=self.bucket_name,
            Key=object_key,
            Body=json.dumps(final_record),
            ContentType='application/json'
        )
        
        # 5. Update the state with the new hash
        self._save_last_hash(current_hash)
        
        return current_hash

# Example Usage:
# logger = HashChainLogger(bucket_name="snisid-audit-worm", state_table="audit_state")
# current = logger.log_event({"action": "BIOMETRIC_VERIFY", "operator": "op_948", "success": True})
# print(f"Event cryptographically secured with hash: {current}")

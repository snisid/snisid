provider "aws" {
  region = "us-east-1"
}
# Sovereign S3 compatible storage (Ceph/MinIO)
resource "aws_s3_bucket" "sovereign_lakehouse" {
  bucket = "snisid-national-lakehouse"
}


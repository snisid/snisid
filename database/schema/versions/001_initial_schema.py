"""initial_schema

Revision ID: 001_initial
Revises: 
Create Date: 2026-05-23 00:00:00.000000

"""
from alembic import op
import sqlalchemy as sa

revision = '001_initial'
down_revision = None
branch_labels = None
depends_on = None

def upgrade() -> None:
    # Enable pgcrypto for PII encryption (Requirement #80)
    op.execute('CREATE EXTENSION IF NOT EXISTS pgcrypto;')

    # Create partitioning base table for Audit Logs
    op.execute('''
        CREATE TABLE audit_logs (
            event_id UUID NOT NULL,
            timestamp TIMESTAMPTZ NOT NULL,
            actor_id VARCHAR(255) NOT NULL,
            action VARCHAR(255) NOT NULL,
            resource_type VARCHAR(255),
            resource_id VARCHAR(255),
            ip_address INET,
            PRIMARY KEY (event_id, timestamp)
        ) PARTITION BY RANGE (timestamp);
    ''')

    # Create current year partition
    op.execute('''
        CREATE TABLE audit_logs_2026 PARTITION OF audit_logs
        FOR VALUES FROM ('2026-01-01') TO ('2027-01-01');
    ''')

    # Citizens Table
    op.execute('''
        CREATE TABLE citizens (
            id UUID PRIMARY KEY,
            -- Symmetric encryption using Vault-managed key injected at query time
            national_id_number_enc BYTEA NOT NULL,
            first_name VARCHAR(255) NOT NULL,
            last_name VARCHAR(255) NOT NULL,
            date_of_birth_enc BYTEA NOT NULL,
            created_at TIMESTAMPTZ DEFAULT NOW(),
            updated_at TIMESTAMPTZ DEFAULT NOW()
        );
    ''')

    # Identities Table (CQRS Event Sourcing Store)
    op.execute('''
        CREATE TABLE events_store (
            event_id UUID PRIMARY KEY,
            aggregate_id UUID NOT NULL,
            aggregate_type VARCHAR(50) NOT NULL,
            event_type VARCHAR(100) NOT NULL,
            data JSONB NOT NULL,
            version INT NOT NULL,
            timestamp TIMESTAMPTZ DEFAULT NOW()
        );
        CREATE UNIQUE INDEX idx_aggregate_version ON events_store(aggregate_id, version);
    ''')

    # Biometric Templates Table
    op.execute('''
        CREATE TABLE biometric_templates (
            template_id UUID PRIMARY KEY,
            citizen_id UUID NOT NULL REFERENCES citizens(id),
            modality VARCHAR(50) NOT NULL,
            -- Vault Transit Encrypted payload
            template_data_vault_enc TEXT NOT NULL,
            quality_score FLOAT NOT NULL,
            created_at TIMESTAMPTZ DEFAULT NOW()
        );
    ''')

def downgrade() -> None:
    op.execute('DROP TABLE IF EXISTS biometric_templates;')
    op.execute('DROP TABLE IF EXISTS events_store;')
    op.execute('DROP TABLE IF EXISTS citizens;')
    op.execute('DROP TABLE IF EXISTS audit_logs;')
    op.execute('DROP EXTENSION IF EXISTS pgcrypto;')

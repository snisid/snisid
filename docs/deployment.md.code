# SNISID Deployment Guide

## Prerequisites

- **Windows Server 2016+** or **Ubuntu 20.04+**
- 2 GB RAM minimum, 4 GB recommended
- 500 MB disk space for binaries, plus data storage
- HTTPS certificate (self-signed for testing, CA-signed for production)
- Database: PostgreSQL 13+ or SQL Server 2017+

## Installation

### Windows (MSI)
```powershell
msiexec /i snisid-1.0.0.msi /qn
```

### Windows (NSIS)
```powershell
.\SNISID-1.0.0-Setup.exe /S
```

### Windows Server (Automated)
```powershell
.\deploy\windows-server\deploy.ps1 -CertThumbprint "<thumbprint>"
```

### Linux
```bash
sudo dpkg -i snisid_1.0.0_amd64.deb
# or
sudo rpm -ivh snisid-1.0.0-1.x86_64.rpm
```

## Configuration

Edit `C:\Program Files\SNISID\config.json` (Windows) or `/etc/snisid/config.json` (Linux):

```json
{
  "database": {
    "host": "localhost",
    "port": 5432,
    "name": "snisid",
    "user": "snisid",
    "password": "<secret>"
  },
  "api": {
    "port": 443,
    "tls_cert": "/etc/snisid/cert.pem",
    "tls_key": "/etc/snisid/key.pem"
  },
  "sync": {
    "inbox": "/var/snisid/sync/inbox",
    "archive": "/var/snisid/sync/archive"
  }
}
```

## Verification

1. Check services are running:
   - Windows: `Get-Service snisid-core, snisid-api`
   - Linux: `systemctl status snisid-core snisid-api`

2. Hit the health endpoint:
   ```powershell
   curl -k https://localhost/api/v1/health
   ```

3. Review logs:
   - Windows: `C:\ProgramData\SNISID\logs\`
   - Linux: `journalctl -u snisid-core`

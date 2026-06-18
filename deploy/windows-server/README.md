# SNISID Windows Server Deployment

## Prerequisites
- Windows Server 2016+ or Windows 10/11 Pro
- PowerShell 5.1+ (Admin)
- Build artifacts in `../../build/`

## Steps
```powershell
# Deploy with defaults
.\deploy.ps1

# Deploy with a specific certificate for HTTPS
.\deploy.ps1 -CertThumbprint "THUMBPRINT_HERE"

# Skip IIS installation (if already configured)
.\deploy.ps1 -SkipIIS
```

## What it does
1. Installs IIS with required features (ASP.NET, management tools)
2. Creates an IIS application pool and site
3. Copies API binaries to C:\inetpub\wwwroot\SNISID
4. Configures HTTPS binding with provided certificate
5. Registers Go backend services and starts them
6. Opens firewall port 443

## Verification
- Browse to `https://localhost/health` to verify the API is running
- Run `Get-Service snisid-core, snisid-sync` to confirm services are running

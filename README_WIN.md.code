# SNISID Windows Deployment Guide

This guide describes how to install and manage the SNISID platform on a Windows workstation for local or offline (air-gapped) operations.

## 🚀 One-Click Installation

1.  **Download** the repository/installer package.
2.  **Right-click** `install.ps1` and select **Run with PowerShell**.
3.  The installer will:
    - Verify and enable **WSL2**.
    - Verify **Docker Desktop**.
    - Install **k3d**.
    - Bootstrap the local **Kubernetes** cluster.
    - Create a **Desktop Shortcut** to the SNISID Dashboard.

## 🛠️ Platform Management

Use the `SNISIDManager.ps1` script to manage the platform lifecycle:
```powershell
# Check status
.\SNISIDManager.ps1 status

# Stop the platform
.\SNISIDManager.ps1 stop

# Start the platform
.\SNISIDManager.ps1 start
```

## 📶 Offline SOC Mode (Air-Gapped)

To deploy on a machine without internet access:
1.  On an online machine, run `.\scripts\export_offline.ps1` to generate `offline_images.tar`.
2.  Copy the entire `SNISID` folder and `offline_images.tar` to the offline machine.
3.  Run `.\install.ps1`. The installer will automatically detect and import the offline images into the local k3d cluster.

## 🔧 Prerequisites
- **Windows 10/11 Pro/Enterprise**
- **Docker Desktop** installed and configured for WSL2 backend.
- **PowerShell 5.1+** with execution policy set to `Bypass`.

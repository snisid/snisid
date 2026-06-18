@echo off
wpeinit
title SNISID Emergency Operations

echo ============================================
echo      SNISID - Offline Emergency Interface
echo ============================================
echo.
echo Loading SNISID environment...

set SNISID_ROOT=X:\SNISID
set PATH=%SNISID_ROOT%;%PATH%

if exist "%SNISID_ROOT%\snisid-offline.cmd" (
    call "%SNISID_ROOT%\snisid-offline.cmd"
) else (
    echo.
    echo SNISID tools available at %SNISID_ROOT%
    echo Run snisid-viewer.exe to launch the offline viewer.
    echo.
    X:
    cd \SNISID
    cmd /k
)

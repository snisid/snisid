!define PRODUCT_NAME "SNISID"
!define PRODUCT_VERSION "1.0.0"
!define PRODUCT_PUBLISHER "SNISID Systems"

Name "${PRODUCT_NAME} ${PRODUCT_VERSION}"
OutFile "SNISID-${PRODUCT_VERSION}-Setup.exe"
InstallDir "$PROGRAMFILES64\SNISID"
RequestExecutionLevel admin

Section "Core Services (required)" SEC_CORE
  SectionIn RO
  SetOutPath "$INSTDIR\bin"
  File /r "..\build\bin\*.exe"
  File /r "..\build\bin\*.dll"

  SetOutPath "$INSTDIR\services"
  File /r "..\build\services\*.*"

  WriteUninstaller "$INSTDIR\uninstall.exe"

  nsExec::ExecToLog '"$INSTDIR\bin\snisid-svc.exe" install'
  nsExec::ExecToLog 'sc start snisid-core'
SectionEnd

Section "API Service" SEC_API
  SetOutPath "$INSTDIR\api"
  File /r "..\build\api\*.*"

  nsExec::ExecToLog '"$INSTDIR\bin\snisid-api.exe" install'
  nsExec::ExecToLog 'sc start snisid-api'
SectionEnd

Section "Frontend" SEC_FRONTEND
  SetOutPath "$INSTDIR\frontend"
  File /r "..\build\frontend\*.*"
SectionEnd

Section "Offline Tools" SEC_TOOLS
  SetOutPath "$INSTDIR\tools"
  File /r "..\build\tools\*.*"
SectionEnd

Section - "Registry & PATH"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}" \
    "DisplayName" "${PRODUCT_NAME} ${PRODUCT_VERSION}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}" \
    "UninstallString" "$INSTDIR\uninstall.exe"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}" \
    "DisplayIcon" "$INSTDIR\bin\snisid.exe,0"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}" \
    "Publisher" "${PRODUCT_PUBLISHER}"
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}" \
    "NoModify" 1
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}" \
    "NoRepair" 1

  EnVar::AddValue "PATH" "$INSTDIR\bin"
SectionEnd

Section "Start Menu Shortcuts" SEC_SHORTCUTS
  CreateDirectory "$SMPROGRAMS\SNISID"
  CreateShortCut "$SMPROGRAMS\SNISID\SNISID Console.lnk" "$INSTDIR\bin\snisid.exe" "" "$INSTDIR\bin\snisid.exe" 0
  CreateShortCut "$SMPROGRAMS\SNISID\Uninstall.lnk" "$INSTDIR\uninstall.exe" "" "$INSTDIR\uninstall.exe" 0
  CreateShortCut "$DESKTOP\SNISID.lnk" "$INSTDIR\bin\snisid.exe" "" "$INSTDIR\bin\snisid.exe" 0
SectionEnd

Section "Uninstall"
  nsExec::ExecToLog 'sc stop snisid-core'
  nsExec::ExecToLog 'sc delete snisid-core'
  nsExec::ExecToLog 'sc stop snisid-api'
  nsExec::ExecToLog 'sc delete snisid-api'

  RMDir /r "$INSTDIR"
  RMDir /r "$SMPROGRAMS\SNISID"
  Delete "$DESKTOP\SNISID.lnk"

  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}"
  EnVar::DeleteValue "PATH" "$INSTDIR\bin"
SectionEnd

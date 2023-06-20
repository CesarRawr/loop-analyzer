@echo off
cd /d "%~dp0"
if NOT "%cd%"=="%cd: =%" (
  echo El directorio actual contine espacios en el path.
  echo Este comando debe estar en un path que no contenga espacios. 
  rundll32.exe cmdext.dll,MessageBeepStub
  pause
  echo.
  goto :EOF
)

if {%1} EQU {[adm]} goto :data
REG QUERY HKU\S-1-5-19\Environment >NUL 2>&1 && goto :data

set command="""%~f0""" [adm] %*
setlocal enabledelayedexpansion
set "command=!command:'=''!"

powershell -NoProfile Start-Process -FilePath '%COMSPEC%' ^
-ArgumentList '/c """!command!"""' -Verb RunAs 2>NUL
goto :EOF

:data
setlocal enabledelayedexpansion
if {%1} EQU {[adm]} (
 set adm=%1
 shift
) ELSE (set adm=)

:cuerpo

mkdir "C:\Program Files\Analyzer"
icacls "C:\Program Files\Analyzer" /grant SYSTEM:F
copy "bin\analyzer.exe" "C:\Program Files\Analyzer"
copy "bin\service.exe" "C:\Program Files\Analyzer"
icacls "C:\Program Files\Analyzer\analyzer.exe" /grant SYSTEM:F
icacls "C:\Program Files\Analyzer\service.exe" /grant SYSTEM:F
"C:\Program Files\Analyzer\service.exe" -service install

:fin
pause
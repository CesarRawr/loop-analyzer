mkdir "C:\Program Files\Analyzer"
"C:\Windows\System32\icacls.exe" "C:\Program Files\Analyzer" /grant SYSTEM:F
copy "E:\Analyzer\bin\analyzer.exe" "C:\Program Files\Analyzer"
copy "D:\Analyzer\bin\analyzer.exe" "C:\Program Files\Analyzer"
copy "J:\Analyzer\bin\analyzer.exe" "C:\Program Files\Analyzer"
copy "F:\Analyzer\bin\analyzer.exe" "C:\Program Files\Analyzer"
copy "E:\Analyzer\bin\service.exe" "C:\Program Files\Analyzer"
copy "D:\Analyzer\bin\service.exe" "C:\Program Files\Analyzer"
copy "J:\Analyzer\bin\service.exe" "C:\Program Files\Analyzer"
copy "F:\Analyzer\bin\service.exe" "C:\Program Files\Analyzer"
"C:\Windows\System32\icacls.exe" "C:\Program Files\Analyzer\analyzer.exe" /grant SYSTEM:F
"C:\Windows\System32\icacls.exe" "C:\Program Files\Analyzer\service.exe" /grant SYSTEM:F
"C:\Program Files\Analyzer\service.exe" -service install
REG ADD HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Analyzer\Parameters /v Application /t REG_SZ /d "C:\Program Files\Analyzer\service.exe"
timeout /t 1200
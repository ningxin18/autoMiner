@echo off
SET autoMiner=autoMiner.exe
if exist %autoMiner% (
    echo Create autoMiner
) else (
    go build .\autoMiner.go
)

start usedStart.bat

SET userCert=C:\%homepath%\AppData\Roaming\usechain\user.crt
SET userDataFile=C:\%homepath%\AppData\Roaming\usechain\userData.json
if exist %userCert% (
    echo You already have user.crt, you can register and miner directly, please continue...
) else (
    echo Get certificate
    if exist %userDataFile% (
       choice  /c YN /m "You already have userData.json, Do you want to modify? Please input yes(Y) or no(N)" /d N /t 5
       if %errorlevel%==1 echo laksjdlf
       if %errorlevel%==2 echo You choosed no
    ) else (
        used.exe verify
    )
)

for /F %%i in ('used.exe verify --info=userData.json --photo="use.jpg;use2.jpg') do ( set commitid=%%i)
echo commitid=%commitid%


echo start autoMiner
autoMiner.exe
pause
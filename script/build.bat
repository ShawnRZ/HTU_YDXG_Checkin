SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o main
.\tools\zip.exe -q -r HTU_YDXG_Checkin.zip main mail.html configs
del main
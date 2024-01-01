rm -r ./dist
GOARCH=amd64 GOOS=linux go build -o dist/linux-x86_64 .
GOARCH=arm64 GOOS=linux go build -o dist/linux-aarch64 .
GOARCH=amd64 GOOS=windows go build -o dist/windows-amd64.exe .
rm -r ./dist
GOARCH=arm64 GOOS=linux go build -o dist/linux-aarch64 .
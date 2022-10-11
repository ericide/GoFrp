mkdir build
#export GOARCH=arm64
env GOOS=linux GOARCH=arm64 go build -o build/mainarm main.go
env GOOS=linux GOARCH=amd64 go build -o build/mainamd main.go

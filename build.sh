mkdir build
#export GOARCH=arm64
env GOOS=linux GOARCH=arm64 go build -o build/frp_linux_arm64 main.go
env GOOS=linux GOARCH=arm go build -o build/frp_linux_armv7 main.go
env GOOS=linux GOARCH=amd64 go build -o build/frp_linux_amd64 main.go
env GOOS=darwin GOARCH=arm64 go build -o build/frp_darwin_arm64 main.go
env GOOS=darwin GOARCH=amd64 go build -o build/frp_darwin_amd64 main.go
env GOOS=windows GOARCH=arm64 go build -o build/frp_windows_arm64 main.go
env GOOS=windows GOARCH=amd64 go build -o build/frp_windows_amd64 main.go
# Build for Windows
go build main.go -o "build/banwords_win.exe"

# Build for Linux
$Env:GOOS = "linux"; $Env:GOARCH = "amd64"
go build main.go -o "build/banwords_amd64"
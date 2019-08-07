# disable CGO for cross-compiling
$Env:CGO_ENABLED="0"

# compile for windows
# note this only compiles for amd64
go build -o release/windows/amd64/drone-gc github.com/drone/drone-gc

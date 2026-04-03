$ErrorActionPreference = "Stop"

$projectRoot = Split-Path -Parent $PSScriptRoot
$swagPath = Join-Path $env:USERPROFILE "go\\bin\\swag.exe"

if (-not (Test-Path $swagPath)) {
    throw "未找到 swag.exe，请先执行：go install github.com/swaggo/swag/cmd/swag@latest"
}

Push-Location $projectRoot
try {
    & $swagPath init --parseInternal -g ./cmd/api/main.go
}
finally {
    Pop-Location
}

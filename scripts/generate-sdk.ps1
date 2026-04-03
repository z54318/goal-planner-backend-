$ErrorActionPreference = "Stop"

$projectRoot = Split-Path -Parent $PSScriptRoot
$sdkOutput = Join-Path $projectRoot "sdk\\typescript-axios"
$generatorCli = Join-Path $projectRoot "node_modules\\.bin\\openapi-generator-cli.cmd"

Push-Location $projectRoot
try {
    & (Join-Path $PSScriptRoot "generate-openapi.ps1")

    if (-not (Test-Path $generatorCli)) {
        throw "未找到 openapi-generator-cli，请先执行：npm install"
    }

    if (Test-Path $sdkOutput) {
        Remove-Item -Recurse -Force $sdkOutput
    }

    & $generatorCli generate `
        -i ./docs/swagger.json `
        -g typescript-axios `
        -o ./sdk/typescript-axios `
        "--additional-properties=npmName=@goal-planner/backend-sdk,npmVersion=1.0.0,withSeparateModelsAndApi=true,apiPackage=api,modelPackage=models,useSingleRequestParameter=true"
}
finally {
    Pop-Location
}

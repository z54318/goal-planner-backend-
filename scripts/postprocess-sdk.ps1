$ErrorActionPreference = "Stop"

$projectRoot = Split-Path -Parent $PSScriptRoot
$sdkRoot = Join-Path $projectRoot "sdk\\typescript-axios"

function Join-Chars {
    param([int[]]$Codes)
    return (-join ($Codes | ForEach-Object { [char]$_ }))
}

function Add-EnumComment {
    param(
        [string]$FilePath,
        [string]$EnumName,
        [string]$CommentBlock
    )

    if (-not (Test-Path $FilePath)) {
        return
    }

    $utf8 = New-Object System.Text.UTF8Encoding($false)
    $content = [System.IO.File]::ReadAllText($FilePath, $utf8)
    if ($content.Contains($CommentBlock)) {
        return
    }

    $target = "export enum $EnumName"
    if (-not $content.Contains($target)) {
        return
    }

    $updated = $content.Replace($target, "$CommentBlock`r`n$target")
    [System.IO.File]::WriteAllText($FilePath, $updated, $utf8)
}

$goalStatusTitle = Join-Chars @(30446, 26631, 29366, 24577)
$draftDesc = Join-Chars @(26410, 25191, 34892)
$activeDesc = Join-Chars @(25191, 34892, 20013)
$completedDesc = Join-Chars @(24050, 23436, 25104)
$archivedDesc = Join-Chars @(24050, 24402, 26723)

$goalStatusComment = @(
    "/**",
    " * $goalStatusTitle",
    " * - draft: $draftDesc",
    " * - active: $activeDesc",
    " * - completed: $completedDesc",
    " * - archived: $archivedDesc",
    " */"
) -join "`r`n"

Add-EnumComment `
    -FilePath (Join-Path $sdkRoot "models\\goal-goal-status.ts") `
    -EnumName "GoalGoalStatus" `
    -CommentBlock $goalStatusComment

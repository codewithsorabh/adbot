# ==============================================================================
# ADBOT Windows Installer (PowerShell)
# Repository: https://github.com/codewithsorabh/adbot
# ==============================================================================

$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'

$Repo = "codewithsorabh/adbot"
$GitHubApi = "https://api.github.com/repos/$Repo"
$GitHubUrl = "https://github.com/$Repo"
$BinName = "adbot.exe"

$InstallDir = Join-Path $env:LOCALAPPDATA "Programs\adbot"

function Write-Step ($message) {
    Write-Host "==> $message" -ForegroundColor Cyan
}

function Write-Info ($message) {
    Write-Host "[INFO] $message" -ForegroundColor Green
}

function Write-Fatal ($message) {
    Write-Host "[ERROR] $message" -ForegroundColor Red
    Exit 1
}

# 1. Architecture Detection
Write-Step "Detecting architecture..."
$Arch = "amd64"
if ([System.Environment]::Is64BitOperatingSystem -eq $false) {
    Write-Fatal "Only 64-bit Windows is supported."
}
Write-Info "Platform: windows/$Arch"

# 2. Determine Release Version
Write-Step "Determining latest ADBOT release..."
$Tag = $env:ADBOT_VERSION
if (-not $Tag) { $Tag = $env:VERSION }

if (-not $Tag) {
    try {
        $ReleaseInfo = Invoke-RestMethod -Uri "$GitHubApi/releases/latest" -Headers @{ "Accept" = "application/vnd.github.v3+json" }
        $Tag = $ReleaseInfo.tag_name
    } catch {
        try {
            $RawVersion = (Invoke-RestMethod -Uri "https://raw.githubusercontent.com/$Repo/main/VERSION").Trim()
            $Tag = "v$RawVersion"
        } catch {
            Write-Fatal "Could not determine latest release version from GitHub."
        }
    }
}

if (-not $Tag.StartsWith("v")) {
    $Tag = "v$Tag"
}
Write-Info "Release version: $Tag"

# 3. Construct URLs & Download
$BinaryFilename = "adbot-$Tag-windows-$Arch.exe"
$DownloadUrl = "$GitHubUrl/releases/download/$Tag/$BinaryFilename"
$ChecksumsUrl = "$GitHubUrl/releases/download/$Tag/checksums.txt"

$TempDir = Join-Path ([System.IO.Path]::GetTempPath()) ([System.Guid]::NewGuid().ToString())
New-Item -ItemType Directory -Force -Path $TempDir | Out-Null
$TempBinary = Join-Path $TempDir $BinaryFilename

try {
    Write-Step "Downloading $BinaryFilename..."
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempBinary -UseBasicParsing

    # 4. Checksum Verification
    Write-Step "Verifying checksum..."
    $ChecksumsFile = Join-Path $TempDir "checksums.txt"
    try {
        Invoke-WebRequest -Uri $ChecksumsUrl -OutFile $ChecksumsFile -UseBasicParsing
        $Lines = Get-Content $ChecksumsFile
        $Match = $Lines | Where-Object { $_ -match $BinaryFilename }
        if ($Match) {
            $ExpectedHash = ($Match -split '\s+')[0].Trim().ToLower()
            $ActualHash = (Get-FileHash -Path $TempBinary -Algorithm SHA256).Hash.ToLower()

            if ($ExpectedHash -ne $ActualHash) {
                Write-Fatal "Checksum mismatch! Expected: $ExpectedHash, Got: $ActualHash"
            }
            Write-Info "Checksum verified ($($ActualHash.Substring(0, 16))...)"
        }
    } catch {
        Write-Info "Checksum file not available; proceeding."
    }

    # 5. Install Binary
    Write-Step "Installing $BinName to $InstallDir..."
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    }

    $TargetExe = Join-Path $InstallDir $BinName
    Copy-Item -Path $TempBinary -Destination $TargetExe -Force

    # 6. Configure PATH
    $UserPath = [System.Environment]::GetEnvironmentVariable('Path', 'User')
    if ($UserPath -notlike "*$InstallDir*") {
        Write-Step "Adding $InstallDir to User PATH..."
        [System.Environment]::SetEnvironmentVariable('Path', "$UserPath;$InstallDir", 'User')
        $env:Path = "$env:Path;$InstallDir"
        Write-Info "PATH updated."
    }

    # 7. Verification
    if (Test-Path $TargetExe) {
        Write-Host ""
        Write-Host "==============================================" -ForegroundColor Green
        Write-Host "        ADBOT INSTALLATION COMPLETE           " -ForegroundColor Green
        Write-Host "==============================================" -ForegroundColor Green
        Write-Host ""
        Write-Info "Installed version: $Tag"
        Write-Info "Location:          $TargetExe"
        Write-Host ""
        Write-Host "Restart your PowerShell/Command Prompt or run:" -ForegroundColor Yellow
        Write-Host "  adbot" -ForegroundColor Cyan
        Write-Host ""
    } else {
        Write-Fatal "Installation failed: executable not found."
    }

} finally {
    if (Test-Path $TempDir) {
        Remove-Item -Recurse -Force -Path $TempDir -ErrorAction SilentlyContinue
    }
}

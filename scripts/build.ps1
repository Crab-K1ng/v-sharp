<#
.SYNOPSIS
    Build script for the project.
.DESCRIPTION
    This PowerShell script automates the build process for the project, including compiling source code,
    running tests, and packaging the application for deployment.
.PARAMETER Configuration
    The build configuration to use (e.g., Debug, Release). Default is 'Release'.
.EXAMPLE
    .\build.ps1 -Configuration Debug -RunTests
    Builds the project using the Debug configuration and runs tests.
.NOTES
    File Name: build.ps1
    Author: Codezz-ops (codezz-ops@protonmail.com)
    Created: 2025-12-26
    Prerequisites: PowerShell 5.1 or higher
    Copyright: (c) 2025 VSharp.
#>

[CmdletBinding()]
param (
    [ValidateSet("Debug", "Release")]
    [string]$Configuration = "Release",
    [switch]$RunTests,
    [switch]$Clean,
    [switch]$VerboseBuild
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Fail {
    param (
        [string]$Message,
        [int]$Code = 1
    )
    Write-Error $Message
    exit $Code
}

function Info {
    param ([string]$Message)
    Write-Host "==> $Message"
}

function Success {
    param ([string]$Message)
    Write-Host "==> $Message"
}

function Warn {
    param ([string]$Message)
    Write-Warning $Message
}

Info "VSharp Build Script"
Info "Configuration: $Configuration"

try {
    Get-Command go -ErrorAction Stop
}
catch {
    Fail "Go is not installed or not found in PATH."
}

Info ("Go version: " + (go version))

try {
    $ProjectRoot = Resolve-Path "$PSScriptRoot/.."
}
catch {
    Fail "Failed to resolve project root."
}

$CmdDir = Join-Path $ProjectRoot "cmd/vsharp"
$BinDir = Join-Path $ProjectRoot "bin"
$GoMod = Join-Path $ProjectRoot "go.mod"

if (-not (Test-Path $GoMod)) {
    Fail "go.mod not found. This is not a Go module."
}

if (-not (Test-Path $CmdDir)) {
    Fail "Command directory not found: $CmdDir"
}

if ($Clean) {
    Info "Cleaning build artifacts"
    if (Test-Path $BinDir) {
        try {
            Remove-Item -Recurse -Force $BinDir
        }
        catch {
            Fail "Failed to clean bin directory."
        }
    }
}

if (-not (Test-Path $BinDir)) {
    try {
        New-Item -ItemType Directory -Path $BinDir | Out-Null
    }
    catch {
        Fail "Failed to create bin directory."
    }
}

$BinaryName = "vsharp"
if ($IsWindows) {
    $BinaryName += ".exe"
}

$OutputPath = Join-Path $BinDir $BinaryName

$Version = "dev"
try {
    $Version = (git describe --tags --dirty --always 2>$null)
}
catch {
    Warn "Git version not available; using '$Version'"
}

if ($RunTests) {
    Info "Running tests"
    Push-Location $ProjectRoot
    try {
        go test ./...
    }
    catch {
        Pop-Location
        Fail "Tests failed. Build aborted."
    }
    Pop-Location
}

$GoArgs = @("build")

if ($VerboseBuild) {
    $GoArgs += "-v"
}

$GoArgs += "-o"
$GoArgs += $OutputPath

$LdFlags = @("-X", "main.version=$Version")

if ($Configuration -eq "Release") {
    $LdFlags += "-s"
    $LdFlags += "-w"
}

$GoArgs += "-ldflags"
$GoArgs += ($LdFlags -join " ")

$GoArgs += "./cmd/vsharp"

Info "Building binary"
Info "Output: $OutputPath"
Info "Version: $Version"

Push-Location $ProjectRoot
try {
    go @GoArgs
}
catch {
    Pop-Location
    Fail "Go build failed."
}
Pop-Location

if (-not (Test-Path $OutputPath)) {
    Fail "Build succeeded but output binary not found."
}

$FileInfo = Get-Item $OutputPath
Info "Binary size: $([math]::Round($FileInfo.Length / 1KB, 2)) KB"

Success "Build completed successfully"
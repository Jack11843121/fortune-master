# 一键构建（Windows）：前端 + 后端 → fortune.exe
$ErrorActionPreference = "Stop"
Set-Location (Join-Path $PSScriptRoot "..")

Write-Host "==> 构建前端 web/dist ..."
Push-Location web
if (-not (Test-Path node_modules)) {
    Write-Host "==> npm install ..."
    npm install
}
npm run build
Pop-Location

Write-Host "==> 编译后端 fortune.exe ..."
go build -o fortune.exe .\cmd\server

Write-Host ""
Write-Host "构建完成：fortune.exe"
Write-Host "运行： .\fortune.exe -config config.yaml  （浏览器打开 http://localhost:8080）"
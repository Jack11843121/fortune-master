package web

import "embed"

// Dist 内嵌前端构建产物（web/dist），随 Go 二进制一起分发，实现“单文件部署”。
//
//go:embed all:dist
var Dist embed.FS
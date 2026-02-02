package web

import "embed"

//go:embed "static" "views/*.html"
var Files embed.FS

// Package web embeds the server-rendered templates and static assets
// (CSS/JS) directly into the compiled binary, so the 32-bit release needs
// no extra files installed alongside it.
package web

import "embed"

//go:embed templates static
var FS embed.FS

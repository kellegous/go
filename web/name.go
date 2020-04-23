package web

import (
	"strings"
)

const encodedIDPrefix = ":"

var bannedNames = map[string]bool{
	"api":     true,
	"edit":    true,
	"healthz": true,
	"links":   true,
	"s":       true,
	"version": true,
	"about":   true,
}

// Parse the shortcut name from the given URL path, given the base URL that is
// handling the request.
func parseName(base, path string) (string, string) {
	t := path[len(base):]
	ix := strings.Index(t, "/")
	if ix == -1 {
		return t, ""
	}
	return t[:ix], t[ix+1:]
}

// Clean a shortcut name. Currently this just means stripping any leading
// ":" to avoid collisions with auto generated names.
func cleanName(name string) string {
	for strings.HasPrefix(name, encodedIDPrefix) {
		name = name[1:]
	}
	return name
}

// Is this name one that was generated from the incrementing id.
func isGenerated(name string) bool {
	return strings.HasPrefix(name, string(genURLPrefix))
}

// isBannedName indicates if the name is one that is reserved by the server?
func isBannedName(name string) bool {
	return bannedNames[name]
}

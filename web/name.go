package web

import "strings"

const encodedIDPrefix = ":"

// Parse the shortcut name from the given URL path, given the base URL that is
// handling the request.
func parseName(base, path string) string {
	t := path[len(base):]
	ix := strings.Index(t, "/")
	if ix == -1 {
		return t
	}
	return t[:ix]
}

// Clean a shortcut name. Currently this just means stripping any leading
// ":" to avoid collisions with auto generated names.
func cleanName(name string) string {
	for strings.HasPrefix(name, encodedIDPrefix) {
		name = name[1:]
	}
	return name
}

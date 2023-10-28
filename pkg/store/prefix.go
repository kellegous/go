package store

import "strings"

func URIPrefix(uri string) string {
	ix := strings.IndexAny(uri, "/?")
	if ix == -1 {
		return uri
	}
	return uri[:ix]
}

package util

import (
	"net/url"
)

func GetURLPathAndQuery(u *url.URL) string {
	r := u.Path
	query := u.Query()
	if len(query) == 0 {
		return r
	}
	return r + "?" + query.Encode()
}

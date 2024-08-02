package api

import (
	"fmt"
	"strings"
)

// Set key value
func (api *API) Set(key, value string) bool {
	return api.db.String.Set(key, value)
}

// Get key
func (api *API) Get(key string) (string, bool) {
	if key == "*" {
		m := api.db.String.GetAll()
		if len(m) == 0 {
			return "", false
		}
		var sb strings.Builder
		for k, v := range m {
			sb.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
		out := sb.String()
		out = strings.TrimSuffix(out, "\n")
		return out, true
	}
	return api.db.String.Get(key)
}

// Del key
func (api *API) Del(key string) bool {
	return api.db.String.Del(key)
}

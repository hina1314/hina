package api

import "encoding/json"

// HSet a key value [key1 value1 ...]
func (api *API) HSet(key string, fields ...string) bool {
	return api.db.HashMap.HSet(key, fields...)
}

// HGet a key
func (api *API) HGet(key, field string) (string, bool) {
	return api.db.HashMap.HGet(key, field)
}

// HGetAll a
func (api *API) HGetAll(key string) (string, bool) {
	m, ok := api.db.HashMap.HGetAll(key)
	jsonData, _ := json.Marshal(m)
	return string(jsonData), ok
}

// HDel a [key ...]
func (api *API) HDel(key string, mem ...string) bool {
	return api.db.HashMap.HDel(key, mem...)
}

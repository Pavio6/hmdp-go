package utils

import "time"

// RedisData mirrors the Java RedisData helper.
type RedisData struct {
	ExpireTime time.Time   `json:"expireTime"`
	Data       interface{} `json:"data"`
}

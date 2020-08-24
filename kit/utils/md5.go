package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(value string) string {
	h := md5.New()
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}

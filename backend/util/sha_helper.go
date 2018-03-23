package util

import (
	"crypto/sha512"
	"encoding/base64"
	"hash"
)

var hasher hash.Hash

var GetHash = func(data []byte) string {
	hasher.Write(data)
	hashed := hasher.Sum(nil)
	hasher.Reset()
	return base64.StdEncoding.EncodeToString(hashed)
}

func init() {
	hasher = sha512.New()
}

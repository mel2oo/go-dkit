package hash

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

func SHA256String(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func MD5String(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

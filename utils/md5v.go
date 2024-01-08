package utils

import (
	"crypto/md5"
	"encoding/hex"
)
const iteration = 2
func MD5V(str string) string {
	salt,_ := GetEnvSalt()
    b := []byte(str)
    s := []byte(salt)
	h := md5.New()
	h.Write(s) 
	h.Write(b)
        var res []byte
	res = h.Sum(nil)
	for i := 0; i < iteration-1; i++ {
		h.Reset()
		h.Write(res)
		res = h.Sum(nil)
	}
	return hex.EncodeToString(res)
}
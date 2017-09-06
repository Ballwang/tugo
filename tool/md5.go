package tool

import (
	"crypto/md5"
	"encoding/hex"
)

//字符串生成MD5
func Md5String(string string) string  {
	h := md5.New()
	h.Write([]byte(string))
	md5String := hex.EncodeToString(h.Sum(nil))
	return md5String
}

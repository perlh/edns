package core

import "encoding/base64"

func Base64Encoding(str string) string { //Base64编码
	src := []byte(str)
	res := base64.StdEncoding.EncodeToString(src) //将编码变成字符串
	return res
}

func Base64Decoding(str string) string { //Base64解码
	res, _ := base64.StdEncoding.DecodeString(str)
	return string(res)
}

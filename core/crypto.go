package core

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
)

func MD5(str string) string {
	data := []byte(str) //切片
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

func GenMd5(code string) string {
	//c1 := md5.Sum([]byte(code)) //返回[16]byte数组

	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}

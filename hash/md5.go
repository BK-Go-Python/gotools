package hash

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

// 将一个文件转为md5
func Filetomd5(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", errors.New(
			fmt.Sprintf("md5.go hash.Filetomd5 os open error %v", err))
	}
	h := md5.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return "", errors.New(fmt.Sprintf("md5.go hash.FileMd5 io copy error %v", err))
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// 将一个字符串转化为md5
func StringMd5(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}
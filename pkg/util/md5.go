package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"strings"
)

func Md5Sum(reader io.Reader) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var returnMD5String string
	hash := md5.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}
func Md5SumString(input string) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	buffer := strings.NewReader(input)
	return Md5Sum(buffer)
}

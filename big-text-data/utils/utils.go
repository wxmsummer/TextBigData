package utils

import "github.com/axgle/mahonia"

// srcStr 待转码的原始字符串 srcEncoding 源编码  dstEncoding 目标编码
func ConvertEncoding(srcStr, srcEncoding, dstEncoding string) (dstStr string, err error) {
	srcDecoder := mahonia.NewDecoder(srcEncoding)
	dstDecoder := mahonia.NewDecoder(dstEncoding)
	utfStr := srcDecoder.ConvertString(srcStr)
	_, dstBytes, err := dstDecoder.Translate([]byte(utfStr), true)
	if err != nil {
		return
	}
	dstStr = string(dstBytes)
	return
}
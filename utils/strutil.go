package utils

import (
	"strings"
)

type StrUtil struct {
}

func (*StrUtil) SplitStrToNum(str string) (intstr, decstr string, err error) {
	return splitStrToNum(str)
}

func splitStrToNum(str string) (intstr, decstr string, err error) {
	err = canPraseBigFloat(str)
	if err != nil {
		return
	}
	arr := strings.Split(str, ".")
	count := len(arr)
	intstr = arr[0]
	if count == 2 {
		decstr = arr[1]
	}
	return
}

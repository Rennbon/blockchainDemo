package utils

import (
	"github.com/go-errors/errors"
	"regexp"
)

type RegUtil struct {
}

//验证字符串是否是满足是
//不支持符号
func (*RegUtil) CanPraseBigFloat(str string) error {
	return canPraseBigFloat(str)
}

func canPraseBigFloat(str string) error {
	reg, err := regexp.Compile(`^[0-9]+\.{0,1}[0-9]*$`)
	if err != nil {
		return err
	}
	bo := reg.FindString(str)
	if bo == "" {
		return errors.New("String validation failed.")
	}
	return nil
}

package utils_test

import (
	"github.com/Rennbon/blockchainDemo/utils"
	"testing"
)

var su utils.StrUtil

func TestStrUtil_SplitStrToNum(t *testing.T) {
	resmap := map[string]bool{
		"12312312332":                           true,
		"1111111.33333":                         true,
		"812391231237123123.111123123812399922": true,
		"000000.0000000":                        true,
		"a11231312313":                          false,
		"123123我.1231231":                       false,
		"落霞与孤鹜齐飞，秋水共长天一色":                       false,
		"1231^&.123123":                         false,
	}
	failure := make([]string, 0, 0)
	for k, v := range resmap {
		l, f, err := su.SplitStrToNum(k)
		if (err == nil) == v {
			t.Logf("字符串:'%s'\r\n验证成功\r\nleft:%s\r\nright:%s\r\n", k, l, f)
		} else {
			t.Errorf("字符串:'%s'\r\n验证失败，错误为:%s\r\n", k, err)
			failure = append(failure, k)
		}
	}
	if len(failure) > 0 {
		t.Fail()
	}
}

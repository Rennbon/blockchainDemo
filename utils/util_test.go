package utils_test

import (
	"github.com/Rennbon/blockchainDemo/utils"
	"testing"
)

var su utils.StrUtil

/////////////////////测试用实体///////////////////////////
type testModle struct {
	from, to, target string
	gap, index       int
}

var slc = []testModle{
	{index: 0, target: "前移2位不溢出", from: "12345.6789", to: "123.456789", gap: 2},
	{index: 1, target: "后移2位不溢出", from: "12345.6789", to: "1234567.89", gap: -2},
	{index: 2, target: "前移10位溢出补0", from: "12345.6789", to: "0.00000123456789", gap: 10},
	{index: 3, target: "后移10位溢出补0，无小数点", from: "12345.6789", to: "123456789000000", gap: -10},
	{index: 4, target: "前移正好", from: "12345.6789", to: "0.123456789", gap: 5},
	{index: 5, target: "后移正好", from: "12345.6789", to: "123456789", gap: -4},
}

////////////////////////////////////////////////////////

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
		l, f, err := su.SplitStrToNum(k, true)
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
func BenchmarkStrUtil_SplitStrToNum(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ { //use b.N for looping
		su.SplitStrToNum("12345.67890", false)
	}
}

func TestStrUtil_MoveDecimalPosition(t *testing.T) {
	for k, v := range slc {
		/*if k != 4 {
			continue
		}*/
		str, err := su.MoveDecimalPosition(v.from, v.gap, false)
		if err != nil {
			t.Errorf("测试失败，下标:%d\r\n")
			t.Fail()
		} else if str != v.to {
			t.Errorf("测试失败,下标：%d\r\n原文参数:%s\r\n预期结果:%s\r\nGap    :%d\r\n实际结果:%s\r\n原始目标:%s\r\n", k, v.from, v.to, v.gap, str, v.target)
			t.Fail()
		} else {
			t.Logf("测试成功,下标：%d\r\nn原文参数:%s\r\n预期结果:%s\r\nGap     :%d\r\n实际结果:%s\r\n原始目标:%s\r\n", k, v.from, v.to, v.gap, str, v.target)
		}
	}

}

func BenchmarkStrUtil_MoveDecimalPosition(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ { //use b.N for looping
		for _, v := range slc {
			su.MoveDecimalPosition(v.from, v.gap, false)

		}
	}
}

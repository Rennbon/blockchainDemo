package coins

import (
	"bytes"
	"github.com/Rennbon/blockchainDemo/errors"
	"github.com/Rennbon/blockchainDemo/utils"
	"math/big"
)

var regutil utils.RegUtil
var strutil utils.StrUtil

//接入币种实现接口
type CoinAmounter interface {
	FloatToCoinAmout(f float64) (*CoinAmount, error)
	//获取新amount
	//num:数值
	//trgt：目标精度
	StringToCoinAmout(num string) (*CoinAmount, error)
	//获取精度及单位
	GetUnitPrec(cu CoinUnit) (cup *CoinUnitPrec)
}

//代币单位
type CoinUnitName string
type CoinUnit int8
type CoinAmount struct {
	Amount        *big.Int //基本代币开始计算，为正整数，如比特币的聪，最基本单位
	*CoinUnitPrec          //默认tostring的单位
}
type CoinUnitPrec struct {
	Prec     int          //小数精度
	UnitName CoinUnitName //单位字符串
}

//单位层次，普通单位上下各三层，一共七层
//如果碰到不够的再加
const (
	CoinBilli    CoinUnit = 9
	CoinMega     CoinUnit = 6
	CoinKilo     CoinUnit = 3
	CoinOrdinary CoinUnit = 0
	CoinMilli    CoinUnit = -3
	CoinMicro    CoinUnit = -6
	CoinBox      CoinUnit = -8
)

func (ca *CoinAmount) Add(amount *CoinAmount) {

}
func (ca *CoinAmount) String() string {
	if &ca == nil {
		return ""
	}
	str := ca.Amount.String()
	length := len(str)
	buff := &bytes.Buffer{}
	if ca.Prec >= length {
		buff.WriteString("0.")
		for i := 0; i < ca.Prec-length+1; i++ {
			buff.WriteString("0")
		}
		buff.WriteString(str)

	} else {
		buff.WriteString(str[:length-ca.Prec])
		buff.WriteString(".")
		buff.WriteString(str[length-ca.Prec:])
	}
	return buff.String()
}

func stringToAmount(str string, cb CoinUnit, gupfunc getUnitPrec, origin CoinUnit) (ca *CoinAmount, err error) {
	l, r, err := strutil.SplitStrToNum(str, false)
	if err != nil {
		return
	}
	gap := int(cb - origin)
	rlen := len(r)
	if rlen > gap {
		err = errors.ERR_PARAM_FAIL
		return
	}
	buff := &bytes.Buffer{}
	buff.WriteString(l)
	buff.WriteString(r)
	for i := 0; i < gap-rlen; i++ {
		buff.WriteString("0")
	}
	amt := big.NewInt(0)
	amt.SetString(buff.String(), 10)
	ca = &CoinAmount{
		Amount:       amt,
		CoinUnitPrec: gupfunc(cb),
	}
	return
}

/*
//转换string类型数字为整数部分和小数部分
func splitStrToNum(str string, cb CoinUnit, gupfunc getUnitPrec, origin CoinUnit) (ca *CoinAmount, err error) {
	l, r, err := strutil.SplitStrToNum(str, false)
	if err != nil {
		return
	}
	ca = &CoinAmount{
		CoinUnit:     cb,
		CoinUnitPrec: gupfunc(cb),
	}
	ltmp := big.NewInt(0)
	bl := false
	if ca.IntPart, bl = ltmp.SetString(l, 10); !bl {
		err = errors.ERR_PARAM_FAIL
		return
	}
	buff := &bytes.Buffer{}
	buff.WriteString(l)
	if r != "" {
		r = "0." + r
		buff.WriteString(r)
		if ca.DecPart, err = strconv.ParseFloat(r, 64); err != nil {
			return
		}
	}
	gap := int(cb - origin)
	if gap > 0 {
		for i := 0; i < gap; i++ {
			buff.WriteString("0")
		}
	}
	ca.Z.SetString(buff.String(), 10)
	return
}*/

/*//转换string类型数字为整数部分和小数部分
func splitStrToNum(str string, cb CoinUnit, gupfunc getUnitPrec, origin CoinUnit) (ca *CoinAmount, err error) {
	l, r, err := strutil.SplitStrToNum(str, false)
	if err != nil {
		return
	}
	ca = &CoinAmount{
		CoinUnit:     cb,
		CoinUnitPrec: gupfunc(cb),
	}
	ltmp := big.NewInt(0)
	bl := false
	if ca.IntPart, bl = ltmp.SetString(l, 10); !bl {
		err = errors.ERR_PARAM_FAIL
		return
	}
	buff := &bytes.Buffer{}
	buff.WriteString(l)
	if r != "" {
		r = "0." + r
		buff.WriteString(r)
		if ca.DecPart, err = strconv.ParseFloat(r, 64); err != nil {
			return
		}
	}
	gap := int(cb - origin)
	if gap > 0 {
		for i := 0; i < gap; i++ {
			buff.WriteString("0")
		}
	}
	ca.Z.SetString(buff.String(), 10)
	return
}*/

/*//性能测试效率较低
func convertCoinUnit1(ca *CoinAmount, cb CoinUnit, gupfunc getUnitPrec, origin CoinUnit) (caout *CoinAmount, err error) {
	if ca == nil {
		err = errors.ERR_PARAM_FAIL
		return
	}
	if ca.CoinUnit == cb {
		caout = ca
		return
	}
	gap := int(cb) - int(ca.CoinUnit)
	newnum, err := strutil.MoveDecimalPosition(ca.String(), gap, false)
	if err != nil {
		return
	}
	caout, err = splitStrToNum(newnum, cb, gupfunc, origin)
	return
}

//转换精度
//效率比convertCoinUnit1相对高
func convertCoinUnit(ca *CoinAmount, cb CoinUnit, gupfunc getUnitPrec) (caout *CoinAmount, err error) {
	if ca == nil {
		err = errors.ERR_PARAM_FAIL
		return
	}
	if ca.CoinUnit == cb {
		caout = ca
		return
	}
	caout = &CoinAmount{
		Z:            ca.Z,
		CoinUnitPrec: gupfunc(cb),
	}
	gap := int(cb - ca.CoinUnit)
	if gap > 0 {
		//往左边移动
		lstr := ca.IntPart.String() //int部分全部字符串
		llen := len(lstr)
		int2dec := "" //int部分同步到dec部分的字符串
		if gap > llen {
			//整数部分变0，全变小数
			caout.IntPart = big.NewInt(0)
			int2dec = lstr //同步全部
		} else {
			llstr := lstr[:llen-gap]
			int2dec = lstr[llen-gap:] //同步末尾部分
			ip := big.NewInt(0)
			caout.IntPart, _ = ip.SetString(llstr, 10)
		}
		//decPart开始转换,补位gap数量的长度
		buff := &bytes.Buffer{}
		buff.WriteString("0.")
		for i := 0; i < gap-len(int2dec); i++ {
			buff.WriteString("0")
		}
		buff.WriteString(int2dec)
		decLeft, _ := strconv.ParseFloat(buff.String(), 64)
		pow := math.Pow10(int(-gap))
		//原小数移位+int部分的移位部分
		caout.DecPart = ca.DecPart*pow + decLeft

	} else {
		gap = -gap
		pow := math.Pow10(int(gap))
		//补位
		ln := big.NewInt(0)
		ln = ln.Mul(ca.IntPart, big.NewInt(int64(pow)))
		r := ca.DecPart * pow
		if cb > 0 {
			ln.Add(ln, big.NewInt(int64(r)))
			caout.DecPart = 0
		} else {
			//取出decpart的整数部分累加到intpart
			intr := math.Floor(r)
			ln.Add(ln, big.NewInt(int64(intr)))
			rstr := strconv.FormatFloat(r, 'f', ca.Prec, 64) //只需保留原来的prec就行
			_, rrstr, errinner := strutil.SplitStrToNum(rstr, false)
			if errinner != nil {
				err = errinner
				return
			}
			caout.DecPart, _ = strconv.ParseFloat("0."+rrstr, 64)
		}
		caout.IntPart = ln
	}
	caout.CoinUnit = cb
	return
}*/

type getUnitPrec func(cu CoinUnit) (cup *CoinUnitPrec)

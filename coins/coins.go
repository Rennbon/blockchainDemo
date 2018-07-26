package coins

import (
	"bytes"
	"math"
	"math/big"
	"strconv"

	"github.com/Rennbon/blockchainDemo/errors"
	"github.com/Rennbon/blockchainDemo/utils"
)

var regutil utils.RegUtil
var strutil utils.StrUtil

//接入币种实现接口
type CoinAmounter interface {
	//获取新amount
	//num:数值
	//trgt：目标精度
	NewCoinAmout(num string) (*CoinAmount, error)
	//转换amount精度
	//ca：当前coinAmount实体
	//trgt:目标精度
	ConvertAmountPrec(ca *CoinAmount, trgt CoinUnit) (caout *CoinAmount, err error)
	//获取精度及单位
	GetUnitPrec(cu CoinUnit) (cup *CoinUnitPrec)
}

//代币单位
type CoinUnitName string
type CoinUnit int8
type CoinAmount struct {
	IntPart  *big.Int //整数部分
	DecPart  float64  //小数部分(整数不存值)
	CoinUnit CoinUnit //单位精度
	*CoinUnitPrec
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

func (ca *CoinAmount) String() string {
	if &ca == nil {
		return ""
	}
	var buf bytes.Buffer

	buf.WriteString(ca.IntPart.String())
	fstr := strconv.FormatFloat(ca.DecPart, 'f', ca.Prec, 64)
	l := len(fstr)
	fstr = fstr[1:l]
	buf.WriteString(fstr)
	return buf.String()
}

//转换string类型数字为整数部分和小数部分
func splitStrToNum(str string, cb CoinUnit, gupfunc getUnitPrec) (ca *CoinAmount, err error) {
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
	if r != "" {
		r = "0." + r
		if ca.DecPart, err = strconv.ParseFloat(r, 64); err != nil {
			return
		}
	}
	return
}

//性能测试效率较低
func convertCoinUnit1(ca *CoinAmount, cb CoinUnit, gupfunc getUnitPrec) (caout *CoinAmount, err error) {
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
	caout, err = splitStrToNum(newnum, cb, gupfunc)
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
}

type getUnitPrec func(cu CoinUnit) (cup *CoinUnitPrec)

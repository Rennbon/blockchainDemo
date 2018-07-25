package coins

import (
	"bytes"
	"math"
	"math/big"
	"strconv"

	"github.com/Rennbon/blockchainDemo/errors"
	"github.com/Rennbon/blockchainDemo/utils"
)

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
	Prec     int
	UnitName CoinUnitName //单位字符串

}

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
	//todo 精度问题，按不同币种来，这里实现有误
	fstr := strconv.FormatFloat(ca.DecPart, 'f', ca.Prec, 64)
	l := len(fstr)
	fstr = fstr[1:l]
	buf.WriteString(fstr)
	return buf.String()
}

var strtuil utils.StrUtil

//转换string类型数字为整数部分和小数部分
func splitStrToNum(str string, cb CoinUnit, gupfunc getUnitPrec) (ca *CoinAmount, err error) {
	l, r, err := strtuil.SplitStrToNum(str, false)
	if err != nil {
		return
	}
	ca = &CoinAmount{
		CoinUnit: cb,
	}
	ca.UnitName = gupfunc(cb).UnitName
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
func ConvertcoinUnit(ca *CoinAmount, cb CoinUnit, gupfunc getUnitPrec) (caout *CoinAmount, err error) {
	if ca == nil {
		err = errors.ERR_PARAM_FAIL
		return
	}
	if ca.CoinUnit == cb {
		caout = ca
		return
	}
	gap := int(cb) - int(ca.CoinUnit)
	newnum, err := strtuil.MoveDecimalPosition(ca.String(), gap, false)
	if err != nil {
		return
	}
	caout, err = splitStrToNum(newnum, cb, gupfunc)
	return
}

//转换精度
func ConvertcoinUnit1(ca *CoinAmount, cb CoinUnit, gupfunc getUnitPrec) (caout *CoinAmount, err error) {
	if ca == nil {
		err = errors.ERR_PARAM_FAIL
		return
	}
	if ca.CoinUnit == cb {
		caout = ca
		return
	}
	gap := int(ca.CoinUnit) - int(cb)
	ln := big.NewInt(0)
	pow := math.Pow10(int(gap))
	if gap > 0 {
		//往左边移动
		lstr := ca.IntPart.String()
		llen := len(lstr)
		if gap > llen {
			caout.IntPart = big.NewInt(0)
		} else {
			llstr := lstr[:llen-gap]
			lrstr := lstr[llen-gap:]
			ip := big.NewInt(0)
			ca.IntPart, _ = ip.SetString(llstr, 10)
			lrstr = "0." + lrstr
			lrf, _ := strconv.ParseFloat(lrstr, 64)
			pow := math.Pow10(int(-gap))
			caout.DecPart = ca.DecPart*pow + lrf
		}
	} else {
		gap = -gap
		//补位
		l := ln.Mul(ca.IntPart, big.NewInt(int64(pow)))
		r := ca.DecPart * pow
		if cb > 0 {
			l.Add(l, big.NewInt(int64(r)))
			caout.DecPart = 0
		} else {
			//取出decpart的整数部分累加到intpart
			intr := math.Floor(r)
			l.Add(l, big.NewInt(int64(intr)))
			rstr := strconv.FormatFloat(r, 'f', int(math.Abs(float64(ca.CoinUnit))), 64)
			_, rrstr, errinner := strtuil.SplitStrToNum(rstr, false)
			if errinner != nil {
				err = errinner
				return
			}
			caout.DecPart, _ = strconv.ParseFloat(rrstr, 64)
		}
		caout.IntPart = l
	}
	caout.CoinUnit = cb

	caout.UnitName = gupfunc(cb).UnitName
	return
}

type getUnitPrec func(cu CoinUnit) (cup *CoinUnitPrec)

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

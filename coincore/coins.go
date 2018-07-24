package coincore

import (
	"bytes"
	"math"
	"math/big"
	"strconv"
	"github.com/Rennbon/blockchainDemo/utils"
	"github.com/go-errors/errors"
)


//代币单位
type CoinUnitName string
type CoinUnit int8
type CoinAmount struct {
	IntPart  *big.Int     //整数部分
	DecPart  float64      //小数部分(整数不存值)
	UnitName CoinUnitName //单位字符串
	CoinUnit CoinUnit     //单位精度
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


func (ca *CoinAmount) String(hasUnitName bool) string {
	if &ca == nil {
		return ""
	}
	var buf bytes.Buffer
	buf.WriteString(ca.IntPart.String())
	if ca.CoinUnit < 0 {
		prec := int(math.Abs(float64(ca.CoinUnit)))
		fstr := strconv.FormatFloat(ca.DecPart, 'f', prec, 64)
		l := len(fstr)
		fstr = fstr[1:l]
		buf.WriteString(fstr)
	}
	if hasUnitName {
		buf.WriteString(" ")
		buf.WriteString(string(ca.UnitName))
	}
	return buf.String()
}

var strtuil utils.StrUtil

//转换string类型数字为整数部分和小数部分
func splitStrToNum(str string,cb CoinUnit,method getUnitName)(ca *CoinAmount,err error ){
	l,r,err:= strtuil.SplitStrToNum(str)
	if err!=nil {
		return
	}
	ca =&CoinAmount{
		UnitName:method(cb),
		CoinUnit:cb,
	}
	ltmp:=big.NewInt(0)
	bl :=false
	if ca.IntPart,bl =ltmp.SetString(l,10);!bl{
		err = errors.New("parse fail in splitStrToNum")
		return
	}
	if r!= ""{
		if ca.DecPart,err = strconv.ParseFloat(r,64);err!=nil{
			return
		}
	}
	return
}

type getUnitName func(CoinUnit) (CoinUnitName)


//接入币种实现接口
type CoinAmounter interface {
	//获取新amount
	//num:数值
	//trgt：目标精度
	GetNewOrdinaryAmount(num string)( *CoinAmount,error)
	//转换amount精度
	//ca：当前coinAmount实体
	//trgt:目标精度
	ConvertAmountPrec(ca *CoinAmount, trgt CoinUnit)error
	//获取单位名称
	GetBtcUnitName(CoinUnit) CoinUnitName
}

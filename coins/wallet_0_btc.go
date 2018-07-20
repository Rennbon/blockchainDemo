package coins

import (
	"blockchainDemo/cert"
	"blockchainDemo/errors"
	"encoding/hex"
	"fmt"
	"log"

	"blockchainDemo/database"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type BtcService struct {
	client *rpcclient.Client
}

var (
	certSrv cert.BtcCertService
	btcSrv  BtcService
)

func initBtcClinet() {
	cli, err := rpcclient.New(btcConn, nil)
	if err != nil {
		panic("btc rpcclient error.")
	}
	btcSrv.client = cli
	log.Println("coins=>btc_wallet=>initClinet sccuess.")
}

/*
*获取新的地址
*account:账户名
 */
func (*BtcService) GetNewAddress(account string, mode AcountRunMode) (address, accountOut string, err error) {
	key, err := certSrv.GenerateSimpleKey()
	if err != nil {
		return
	}
	if err = dhSrv.AddAccount(account, key.PrivKey, key.PubKey, key.Address, key.Seed, database.BTC); err != nil {
		return
	}
	switch mode {
	case NoneMode:
		break
	case PrvMode:
		account, err = btcSrv.AddPrvkeyToWallet(key.PrivKey, account)
		break
	case PubMode:
		account, err = btcSrv.AddPubkeyToWallet(key.PubKey, account)
		break
	case AddrMode:
		account, err = btcSrv.AddAddressToWallet(key.PubKey, account)
		break
	default:
		break
	}
	if err != nil {
		return
	}
	return key.Address, account, nil
}

/* 导入privatekey,这个导入能在coin.core上直接listaccounts查看到，因为有私钥了 */
func (*BtcService) AddPrvkeyToWallet(prvkey, accoutIn string) (accountOut string, err error) {
	wif, err := btcutil.DecodeWIF(prvkey)
	if err != nil {
		return
	}
	if err = btcSrv.client.ImportPrivKeyLabel(wif, accoutIn); err != nil {
		return
	}
	return accoutIn, nil
}

/* 将publickey对应的address添加到链中，
pubKey 公钥
account 地址自定义名称 */
func (*BtcService) AddPubkeyToWallet(pubKey, accountIn string) (accountOut string, err error) {
	//验证地址是否已存在
	err = btcSrv.CheckAddressExists(pubKey)
	if err != nil {
		return
	}
	if err = btcSrv.client.ImportPubKey(pubKey); err != nil {
		return
	}
	address, _ := btcutil.DecodeAddress(pubKey, &chaincfg.RegressionNetParams)
	//修改名字 忽略错误
	if err = btcSrv.client.SetAccount(address, accountIn); err != nil {
		return
	}
	return accountIn, nil
}

/* 将publickey对应的address添加到链中，
pubKey 公钥
account 地址自定义名称 */
func (*BtcService) AddAddressToWallet(pubKey, accountIn string) (accountOut string, err error) {
	//验证地址是否已存在
	err = btcSrv.CheckAddressExists(pubKey)
	if err != nil {
		return
	}
	address, _ := btcutil.DecodeAddress(pubKey, &chaincfg.RegressionNetParams)
	if err = btcSrv.client.ImportAddress(address.EncodeAddress()); err != nil {
		return
	}
	//修改名字 忽略错误
	if btcSrv.client.SetAccount(address, accountIn) != nil {
		return
	}
	return accountIn, nil
}

/* 验证publickey对应的地址是否已存在于链中
pubkey 公钥 */
func (*BtcService) CheckAddressExists(pubKey string) error {
	address, err := btcutil.DecodeAddress(pubKey, &chaincfg.RegressionNetParams)
	addrValid, err := btcSrv.client.ValidateAddress(address)
	if err != nil {
		return err
	}
	if addrValid.IsWatchOnly {
		return errors.ERR_DATA_EXISTS
	}
	return nil
}

/*
*获取所有account
 */
func (*BtcService) GetAccounts() (accounts []*Account, err error) {
	accs, err := btcSrv.client.ListAccounts()
	if err != nil {
		return
	}
	for k, v := range accs {
		accounts = append(accounts, &Account{
			Amount: v.ToBTC(),
			Name:   k,
			Unit:   "BTC",
		})
	}
	return accounts, nil
}
func (*BtcService) GetBalanceInAddress(address string) (balance float64, err error) {
	unspents, err := btcSrv.GetUnspentByAddress(address)
	if err != nil {
		return
	}
	for _, v := range unspents {
		balance += v.Amount
	}
	return
}

//根据address获取未花费的tx
func (*BtcService) GetUnspentByAddress(address string) (unspents []btcjson.ListUnspentResult, err error) {
	btcAdd, err := btcutil.DecodeAddress(address, &chaincfg.RegressionNetParams)
	if err != nil {
		return
	}
	adds := [1]btcutil.Address{btcAdd}
	unspents, err = btcSrv.client.ListUnspentMinMaxAddresses(1, 999999, adds[:])
	if err != nil {
		return
	}
	return
}

//转账
//addrForm来源地址，addrTo去向地址
//transfer 转账金额
//fee 小费
func (*BtcService) SendAddressToAddress(addrFrom, addrTo string, transfer, fee float64) (txId string, err error) {
	//数据库获取prv pub key等信息，便于调试--------START------
	actf, err := dhSrv.GetAccountByAddress(addrFrom)
	if err != nil {
		return
	}
	//----------------------------------------END-----------

	unspents, err := btcSrv.GetUnspentByAddress(addrFrom)
	if err != nil {
		return
	}
	//各种参数声明 可以构建为内部小对象
	outsu := float64(0)                 //unspent单子相加
	feesum := fee                       //交易费总和
	totalTran := transfer + feesum      //总共花费
	var pkscripts [][]byte              //txin签名用script
	tx := wire.NewMsgTx(wire.TxVersion) //构造tx

	for _, v := range unspents {
		if v.Amount == 0 {
			continue
		}
		if outsu < totalTran {
			outsu += v.Amount
			{
				//txin输入-------start-----------------
				hash, _ := chainhash.NewHashFromStr(v.TxID)
				outPoint := wire.NewOutPoint(hash, v.Vout)
				txIn := wire.NewTxIn(outPoint, nil, nil)

				tx.AddTxIn(txIn)

				//设置签名用script
				txinPkScript, errInner := hex.DecodeString(v.ScriptPubKey)
				if errInner != nil {
					err = errInner
					return
				}
				pkscripts = append(pkscripts, txinPkScript)
			}
		} else {
			break
		}
	}
	//家里穷钱不够
	if outsu < totalTran {
		err = errors.ERR_NOT_ENOUGH_COIN
		return
	}
	// 输出1, 给form----------------找零-------------------
	addrf, err := btcutil.DecodeAddress(addrFrom, &chaincfg.RegressionNetParams)
	if err != nil {
		return
	}
	pkScriptf, err := txscript.PayToAddrScript(addrf)
	if err != nil {
		return
	}
	baf := int64((outsu - totalTran) * 1e8)
	tx.AddTxOut(wire.NewTxOut(baf, pkScriptf))
	//输出2，给to------------------付钱-----------------
	addrt, err := btcutil.DecodeAddress(addrTo, &chaincfg.RegressionNetParams)
	if err != nil {
		return
	}
	pkScriptt, err := txscript.PayToAddrScript(addrt)
	if err != nil {
		return
	}
	bat := int64(transfer * 1e8)
	tx.AddTxOut(wire.NewTxOut(bat, pkScriptt))
	//-------------------输出填充end------------------------------
	err = sign(tx, actf.PrvKey, pkscripts) //签名
	if err != nil {
		return
	}
	//广播
	txHash, err := btcSrv.client.SendRawTransaction(tx, false)
	if err != nil {
		return
	}
	//这里最好也记一下当前的block count,以便监听block count比此时高度
	//大6的时候去获取当前TX是否在公链有效
	dhSrv.AddTx(txHash.String(), addrFrom, []string{addrFrom, addrTo})
	fmt.Println("Transaction successfully signed")
	fmt.Println(txHash.String())
	return txHash.String(), nil
}

//这个方法ListAddressTransactions method not found;btcd NOTE: This is a btcwallet extension.
func (*BtcService) GetTxByAddress(addrs []string, name string) (interface{}, error) {
	ct := len(addrs)
	addresses := make([]btcutil.Address, 0, ct)
	for _, v := range addrs {
		address, err := btcutil.DecodeAddress(v, &chaincfg.RegressionNetParams)
		if err != nil {
			log.Println("一个废物")
		} else {
			addresses = append(addresses, address)
		}
	}

	txs, err := btcSrv.client.ListAddressTransactions(addresses, name)
	if err != nil {
		return nil, err
	}
	return txs, nil
}

//验证交易是否被公链证实
//txid:交易id
func (*BtcService) CheckTxMergerStatus(txId string) error {
	txHash, err := chainhash.NewHashFromStr(txId)
	if err != nil {
		return err
	}
	txResult, err := btcSrv.client.GetTransaction(txHash)
	if err != nil {
		return err
	}
	//pow共识机制当6个块确认后很难被修改
	if txResult.Confirmations < 6 {
		return errors.ERR_UNCONFIRMED
	}
	return nil
}

//签名
//privkey的compress方式需要与TxIn的
func sign(tx *wire.MsgTx, privKey string, pkScripts [][]byte) error {
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return err
	}
	/* lookupKey := func(a btcutil.Address) (*btcec.PrivateKey, bool, error) {
		return wif.PrivKey, false, nil
	} */
	for i, _ := range tx.TxIn {
		script, err := txscript.SignatureScript(tx, i, pkScripts[i], txscript.SigHashAll, wif.PrivKey, false)
		//script, err := txscript.SignTxOutput(&chaincfg.RegressionNetParams, tx, i, pkScripts[i], txscript.SigHashAll, txscript.KeyClosure(lookupKey), nil, nil)
		if err != nil {
			return err
		}
		tx.TxIn[i].SignatureScript = script
		vm, err := txscript.NewEngine(pkScripts[i], tx, i,
			txscript.StandardVerifyFlags, nil, nil, -1)
		if err != nil {
			return err
		}
		err = vm.Execute()
		if err != nil {
			return err
		}
		log.Println("Transaction successfully signed")
	}
	return nil
}

//离线签名signTxOut是获取keyDB使用，区分addres的compress状态
func mkGetKey(keys map[string]addressToKey) txscript.KeyDB {
	if keys == nil {
		return txscript.KeyClosure(func(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
			return nil, false, errors.ERR_NOPE
		})
	}
	return txscript.KeyClosure(func(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
		a2k, ok := keys[addr.EncodeAddress()]
		if !ok {
			return nil, false, errors.ERR_NOPE
		}
		return a2k.key, a2k.compressed, nil
	})
}

type addressToKey struct {
	key        *btcec.PrivateKey
	compressed bool
}

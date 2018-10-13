package wallet

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"crypto/elliptic"
	"io/ioutil"
	"../utils"
)

type Wallets struct {
	WalletMap map[string]*Wallet
}

/**
从本地读取文件解析出钱包集
 */
func GetWallets() *Wallets {
	if !utils.WalletsExists() {
		fmt.Println("本地钱包文件不存在。")
		w := &Wallets{}
		w.WalletMap = make(map[string]*Wallet)
		return w
	}

	walletsByte, err := ioutil.ReadFile(utils.WalletName())
	utils.PanicErr(err)
	return DeserializeWallets(walletsByte)
}

/**
新建一个钱包地址，并更新本地的钱包文件
 */
func (w *Wallets) NewWalletAndPersist() {
	wallet := NewWallet()
	address := wallet.GetAddress()
	fmt.Printf("新生成地址：%s\n", address)
	w.WalletMap[string(address)] = wallet

	walletsByte := w.serializeWallets()
	// WriteFile writes data to a file named by filename.
	// If the file does not exist, WriteFile creates it with permissions perm;
	// otherwise WriteFile truncates it before writing.
	err := ioutil.WriteFile(utils.WalletName(), walletsByte, 0644)
	utils.PanicErr(err)
}

/**
序列化钱包集
 */
func (w *Wallets) serializeWallets() []byte {
	var buff bytes.Buffer
	gob.Register(elliptic.P256()) //要注册使用到的interface

	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(*w)
	utils.PanicErr(err)
	return buff.Bytes()
}

/**
反序列化钱包集
 */
func DeserializeWallets(walletsByte []byte) *Wallets {
	var w *Wallets
	gob.Register(elliptic.P256()) //要注册使用到的interface

	reader := bytes.NewReader(walletsByte)
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(&w)
	utils.PanicErr(err)
	return w
}

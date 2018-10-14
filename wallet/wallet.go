package wallet

import (
	"crypto/ecdsa"
	"../utils"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte
}

func NewWallet()*Wallet  {
	privateKey, publicKey := utils.NewKeyPair()
	return &Wallet{privateKey,publicKey}
}

/**
根据公钥hash，生成地址
 */
func (w *Wallet) GetAddress() []byte {
	// step1. 原始公钥-->sha256-->ripemd160: 公钥哈希
	publicKeyHash := utils.TwiceHash(w.PublicKey)

	// step2. (版本号+公钥哈希)--->sha256 ---> sha256: 校验码
	afterAddVersion := append([]byte{utils.AddressVersion}, publicKeyHash...)
	checkSumBytes := utils.CheckSum(afterAddVersion)

	// step3. (版本号+公钥哈希+校验码)--->进行Base58编码：地址
	payload := append(afterAddVersion, checkSumBytes...)
	address := utils.Base58Encode(payload)
	return address
}

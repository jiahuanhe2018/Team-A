package utils

import (
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"bytes"
)

/**
由原始公钥进行两次哈希得到公钥hash
 */
func TwiceHash(publicKey []byte) []byte {
	// step1. 原始公钥sha256哈希
	s256 := sha256.Sum256(publicKey)

	// step2. 继续进行一次ripemd160哈希
	r160 := ripemd160.New()
	r160.Write(s256[:])
	hashResult := r160.Sum(nil)

	return hashResult
}

/**
根据地址获取公钥哈希
 */
func GetPubKeyHashFromAddr(address string) []byte {
	// step1. 地址Base58解码：版本号+公钥哈希+校验码
	payload := Base58Decode([]byte(address))

	// step2. 拆解出公钥哈希
	publicKeyHash := payload[1 : len(payload)-AddressCheckSumLength]
	return publicKeyHash
}

/**
根据校验和检查地址是否合法
 */
func IsValidAddress(address []byte) bool {
	// step1. 地址Base58解码：版本号+公钥哈希+校验码
	payload := Base58Decode(address)

	// step2. 剥离出原始的校验码
	if len(payload)-AddressCheckSumLength < 0 {
		return false
	}
	checkSumBytes := payload[len(payload)-AddressCheckSumLength:]

	// step3. 计算真实的校验码
	afterAddVersion := payload[:len(payload)-AddressCheckSumLength]
	realCheckSumBytes := CheckSum(afterAddVersion)

	return bytes.Compare(checkSumBytes, realCheckSumBytes) == 0
}

/**
两次sha256生成校验码
 */
func CheckSum(input []byte) []byte {
	hash1 := sha256.Sum256(input)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:AddressCheckSumLength]
}

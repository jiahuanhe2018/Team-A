package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
)

/**
生成公私钥对
 */
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	PanicErr(err)

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, publicKey
}

func Sign(privateKey ecdsa.PrivateKey, toSignData []byte) (signature []byte) {
	r, s, err := ecdsa.Sign(rand.Reader, &privateKey, toSignData)
	PanicErr(err)

	signature = append(r.Bytes(), s.Bytes()...)
	return
}

func Verify(publicKey []byte, signature []byte, toVerifyData []byte) bool {
	curve := elliptic.P256()
	x, y := getXY(publicKey)
	rawPublicKey := ecdsa.PublicKey{curve, &x, &y}

	r, s := getXY(signature)
	return ecdsa.Verify(&rawPublicKey, toVerifyData, &r, &s)
}

func getXY(data []byte) (big.Int, big.Int) {
	x := big.Int{}
	y := big.Int{}
	x.SetBytes(data[:len(data)/2])
	y.SetBytes(data[len(data)/2:])
	return x, y
}

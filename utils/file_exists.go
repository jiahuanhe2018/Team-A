package utils

import "os"

func WalletsExists()  bool{
	_, err := os.Stat(WalletName())
	return !os.IsNotExist(err)
}

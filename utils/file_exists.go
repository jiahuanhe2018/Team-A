package utils

import "os"

func WalletsExists()  bool{
	return FileIsExist(WalletName())
}

func FileIsExist(file string) bool {
	fi, err := os.Stat(file)
	if err != nil {
		return os.IsExist(err)
	}
	return !fi.IsDir()
}

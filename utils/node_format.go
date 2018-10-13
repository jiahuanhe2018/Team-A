package utils

import (
	"fmt"
	"os"
)

func WalletName() string {
	return fmt.Sprintf(WalletsLocalFileFormat, os.Getenv(NodeIdEnv))
}

func NodeAddress() string {
	return fmt.Sprintf(NodeAddressFormat, os.Getenv(NodeIdEnv))
}

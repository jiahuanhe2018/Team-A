package pow

import (
	"strings"
	"../block"
	"../utils"
	"strconv"
	"fmt"
)

func MineHash(block *block.Block) {
	var hash string
	var nonce = 0
	for {
		block.Nonce = strconv.Itoa(nonce)
		hash = block.CalculateHash()
		fmt.Printf("\r%d --- %s", nonce, hash)
		if isHashValid(hash, utils.PowDifficulty) {
			block.Hash = hash
			fmt.Println()
			break
		}
		nonce ++
	}
}

func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

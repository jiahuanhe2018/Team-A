package pos

import (
	"sync"
	"log"
)

var validators = make(map[string]int)
var valiMutex = &sync.Mutex{}

func AddValidator(address string, balance int) {
	valiMutex.Lock()
	defer valiMutex.Unlock()
	validators[address] = balance
	log.Println("当前所有验证节点为，", validators)
}

func DeleteValidator(address string) {
	valiMutex.Lock()
	defer valiMutex.Unlock()
	delete(validators, address)
	log.Println("当前所有验证节点为，", validators)
}

func GetValidators() map[string]int {
	valiMutex.Lock()
	defer valiMutex.Unlock()
	result := validators
	return result
}

package blockchain

import (
	"../block"
	"../utils"
	"../account"
	TX "../transaction"
	"sync"
	"log"
	"path/filepath"
	"os"
	"encoding/gob"
	"fmt"
)

type Blockchain struct {
	Blocks  []block.Block
	TxPool  *TX.TxPool
	DataDir string
}

var BlcIns = Blockchain{
	TxPool: TX.NewTxPool(),
}

var blcMutex = &sync.Mutex{}

func (blc *Blockchain) NewTransaction(sender string, recipient string, amount uint64, data []byte) *TX.Transaction {
	return &TX.Transaction{
		Amount:    amount,
		Recipient: recipient,
		Sender:    sender,
		Data:      data,
	}
}

func (blc *Blockchain) AddBlock(b block.Block) {
	blcMutex.Lock()
	defer blcMutex.Unlock()
	blc.Blocks = append(blc.Blocks, b)
}

func (blc *Blockchain) LastBlock() block.Block {
	blcMutex.Lock()
	defer blcMutex.Unlock()
	return blc.Blocks[len(blc.Blocks)-1]
}

func (blc *Blockchain) GetBalance(address string) uint64 {
	accounts := blc.LastBlock().Accounts
	if acunt, ok := accounts[address]; ok {
		return acunt.Balance
	}
	return 0
}

func (blc *Blockchain) PackageTx(newBlock *block.Block) {
	newBlock.Transactions = blc.TxPool.AllTx
	accountsMap := blc.LastBlock().Accounts
	for addr, acunt := range accountsMap {
		fmt.Println(addr, "---", acunt)
	}

	unusedTx := make([]TX.Transaction, 0)

	for _, tx := range blc.TxPool.AllTx {
		if senderAccount, ok := accountsMap[tx.Sender]; ok {
			if senderAccount.Balance < tx.Amount {
				// 余额不足以支付交易
				unusedTx = append(unusedTx, tx)
				continue
			}
			senderAccount.Balance -= tx.Amount
			senderAccount.State += 1
			accountsMap[tx.Sender] = senderAccount

			if recipientAccount, ok := accountsMap[tx.Recipient]; ok {
				recipientAccount.Balance += tx.Amount
				accountsMap[tx.Recipient] = recipientAccount
			} else {
				newAccount := account.Account{
					Balance: tx.Amount,
					State:   0,
				}
				accountsMap[tx.Recipient] = newAccount
			}
		}
	}

	// 余额不够的交易放回交易池
	blc.TxPool.Clear()
	if len(unusedTx) > 0 {
		for _, tx := range unusedTx {
			blc.TxPool.AddTx(&tx)
		}
	}

	newBlock.Accounts = accountsMap
}

func (blc *Blockchain) WriteData2File() {
	if blc.DataDir == "" {
		log.Println("未指定区块链文件目录，无法写入")
		return
	}
	filePath := filepath.Join(blc.DataDir, utils.BlockchainDataFileName)

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Println("无法写入区块链数据到文件中，打开文件失败：", err)
		return
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	if err := enc.Encode(blc); err != nil {
		log.Fatal("区块链编码失败：", err)
	}

	fmt.Printf("成功写入区块链到本地目录：%s\n", filePath)
}

func (blc *Blockchain) ReadDataFromFile() {
	if blc.DataDir == "" {
		log.Println("未指定区块链文件目录，无法读取")
		return
	}
	filePath := filepath.Join(blc.DataDir, utils.BlockchainDataFileName)
	if !utils.FileIsExist(filePath) {
		log.Println("未指定区块链文件，无法读取")
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Println("无法从文件读取区块链数据，打开文件失败：", err)
		return
	}
	defer file.Close()
	dec := gob.NewDecoder(file)

	var blcFromFile Blockchain
	if err := dec.Decode(&blcFromFile); err != nil {
		log.Fatal("区块链解码失败：", err)
	}
	BlcIns = blcFromFile
}

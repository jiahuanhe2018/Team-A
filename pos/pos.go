package pos

import (
	"../block"
	blc "../blockchain"
	"sync"
	"time"
	"log"
	"math/rand"
)

var tempBlocks []*block.Block
var candidateBlocks = make(chan *block.Block)

var hasBeenValid = make(chan int, 1)

var posMutex = &sync.Mutex{}

func PoS() {
	//	从通道接收候选区块并放入临时区块池中
	go func() {
		for candidate := range candidateBlocks {
			posMutex.Lock()
			tempBlocks = append(tempBlocks, candidate)
			posMutex.Unlock()
		}
	}()

	//	定期从临时区块池中用pos算法选出最终区块
	go func() {
		for {
			pickWinner()
		}
	}()
}

func pickWinner() {
	time.Sleep(10 * time.Second)
	posMutex.Lock()
	pickedTempBlocks := tempBlocks
	tempBlocks = []*block.Block{}
	posMutex.Unlock()
	pickedValidators := GetValidators()

	if len(pickedTempBlocks) == 0 {
		log.Println("此轮竞选无有效区块")
		return
	}

	// 构建奖池
	lotteryPool := make([]string, 0)
OUTER:
	for _, b := range pickedTempBlocks {
		for _, node := range lotteryPool {
			// 奖池中每个见证节点只处理一次
			if node == b.Validator {
				continue OUTER
			}
		}

		if balance, ok := pickedValidators[b.Validator]; ok {
			for i := 0; i < balance; i ++ {
				lotteryPool = append(lotteryPool, b.Validator)
			}
		}
	}

	//	抽奖选出此轮见证节点
	if len(lotteryPool) == 0 {
		log.Println("此轮竞选无有效见证节点")
		return
	}
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	winner := lotteryPool[r.Intn(len(lotteryPool))]

	// 将见证节点的区块选一个出来作为最终区块，保存到区块链
	for _, b := range pickedTempBlocks {
		if b.Validator == winner {
			blc.BlcIns.AddBlock(*b)
			hasBeenValid <- 1
			break
		}
	}
}

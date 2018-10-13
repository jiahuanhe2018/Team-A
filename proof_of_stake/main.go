package main

import (
	"sync"
	"time"
	"log"
	"github.com/davecgh/go-spew/spew"
	"strconv"
	"encoding/hex"
	"crypto/sha256"
	"net"
	"github.com/joho/godotenv"
	"os"
	"math/rand"
	"fmt"
	"io"
	"bufio"
	"encoding/json"
)

const EnvFileName = "myenv.env"

type BlockZLH struct {
	IndexZLH     int
	TimestampZLH string
	BPMZLH       int
	HashZLH      string
	PrevHashZLH  string
	ValidatorZLH string
}

var BlockchainZLH []*BlockZLH
var tempBlocksZLH []*BlockZLH

var candidateBlocksZLH = make(chan *BlockZLH)
var announcementsZLH = make(chan string)

var BlcMutexZLH = &sync.Mutex{}

var validatorsZLH = make(map[string]int)

func main() {
	err := godotenv.Load(EnvFileName)
	if err != nil {
		log.Fatal(err)
	}

	// 创建创世区块
	genesisBlock := generateGenesisBlockZLH()
	BlockchainZLH = append(BlockchainZLH, genesisBlock)

	// 启动服务端
	port := os.Getenv("PORT")
	server, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("开始监听端口：", port)
	defer server.Close()

	//	从通道接收候选区块并放入临时区块池中
	go func() {
		for candidateBlock := range candidateBlocksZLH {
			BlcMutexZLH.Lock()
			tempBlocksZLH = append(tempBlocksZLH, candidateBlock)
			BlcMutexZLH.Unlock()
		}
	}()

	//	定期从临时区块池中用pos算法选出最终区块
	go func() {
		for {
			pickWinnerZLH()
		}
	}()

	//	处理客户端请求
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnZLH(conn)
	}
}

func pickWinnerZLH() {
	time.Sleep(10 * time.Second)
	BlcMutexZLH.Lock()
	pickedTempBlocks := tempBlocksZLH
	tempBlocksZLH = []*BlockZLH{}
	pickedValidators := validatorsZLH
	BlcMutexZLH.Unlock()

	if len(pickedTempBlocks) == 0 {
		log.Println("此轮竞选无有效区块")
		return
	}

	// 构建奖池
	lotteryPool := make([]string, 0)
OUTER:
	for _, block := range pickedTempBlocks {
		for _, node := range lotteryPool {
			// 奖池中每个见证节点只处理一次
			if node == block.ValidatorZLH {
				continue OUTER
			}
		}

		if balance, ok := pickedValidators[block.ValidatorZLH]; ok {
			for i := 0; i < balance; i ++ {
				lotteryPool = append(lotteryPool, block.ValidatorZLH)
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
	for _, block := range pickedTempBlocks {
		if block.ValidatorZLH == winner {
			BlcMutexZLH.Lock()
			BlockchainZLH = append(BlockchainZLH, block)
			BlcMutexZLH.Unlock()
			for range validatorsZLH {
				announcementsZLH <- fmt.Sprintf("此轮见证节点为%s,BPM为%d\n", winner, block.BPMZLH)
			}
			break
		}
	}

}

func handleConnZLH(conn net.Conn) {
	defer conn.Close()

	// 接收广播消息
	go func() {
		for {
			msg := <-announcementsZLH
			io.WriteString(conn, msg)
		}
	}()

	// 创建新的见证节点
	var address string
	io.WriteString(conn, "请输入持有的Token数量：\n")
	scanBalance := bufio.NewScanner(conn)
	for scanBalance.Scan() {
		balance, err := strconv.Atoi(scanBalance.Text())
		if err != nil {
			log.Printf("%v 不是数字，错误：%v", scanBalance.Text(), err)
			io.WriteString(conn, fmt.Sprintf("%v不是数字，无法创建节点\n", scanBalance.Text()))
			conn.Close()
			return
		}
		address = calculateHashZLH(time.Now().String())
		BlcMutexZLH.Lock()
		validatorsZLH[address] = balance
		BlcMutexZLH.Unlock()
		fmt.Println("所有见证节点为，", validatorsZLH)
		break
	}
	// 产生候选区块
	io.WriteString(conn, "输入BPM：\n")
	scanner := bufio.NewScanner(conn)

	go func() {
		for scanner.Scan() {
			bpm, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Printf("%v不是数字，错误：%v", scanner.Text(), err)
				BlcMutexZLH.Lock()
				delete(validatorsZLH, address)
				BlcMutexZLH.Unlock()
				io.WriteString(conn, fmt.Sprintf("%v不是数字，已自动退出见证节点\n", scanner.Text()))
				conn.Close()
				return
			}

			BlcMutexZLH.Lock()
			lastBlock := BlockchainZLH[len(BlockchainZLH)-1]
			BlcMutexZLH.Unlock()
			newBlock := generateBlockZLH(lastBlock, bpm, address)
			if isBlockValidZLH(newBlock, lastBlock) {
				candidateBlocksZLH <- newBlock
			}
			io.WriteString(conn, "继续输入BPM：\n")
		}
	}()

	// 定时向客户端发送完整的区块链
	for {
		time.Sleep(20 * time.Second)
		BlcMutexZLH.Lock()
		output, err := json.Marshal(BlockchainZLH)
		if err != nil {
			log.Fatal(err)
		}
		BlcMutexZLH.Unlock()
		io.WriteString(conn, fmt.Sprintf("当前区块链：\n%s\n", string(output)))
	}

}

func calculateHashZLH(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func calculateBlockHashZLH(block *BlockZLH) string {
	record := strconv.Itoa(block.IndexZLH) + block.TimestampZLH + strconv.Itoa(
		block.BPMZLH) + block.PrevHashZLH
	h := sha256.Sum256([]byte(record))
	return hex.EncodeToString(h[:])
}

func generateGenesisBlockZLH() *BlockZLH {
	genesisBlock := &BlockZLH{
		IndexZLH:     0,
		TimestampZLH: time.Now().String(),
		BPMZLH:       0,
		HashZLH:      "",
		PrevHashZLH:  "",
		ValidatorZLH: "",
	}
	genesisBlock.HashZLH = calculateBlockHashZLH(genesisBlock)
	log.Println("创建创世区块：")
	spew.Dump(genesisBlock)
	return genesisBlock
}

func generateBlockZLH(oldBlock *BlockZLH, BPM int, address string) *BlockZLH {
	newBlock := &BlockZLH{
		IndexZLH:     oldBlock.IndexZLH + 1,
		TimestampZLH: time.Now().String(),
		BPMZLH:       BPM,
		HashZLH:      "",
		PrevHashZLH:  oldBlock.HashZLH,
		ValidatorZLH: address,
	}
	newBlock.HashZLH = calculateBlockHashZLH(newBlock)
	log.Println("创建新的候选区块：")
	spew.Dump(newBlock)
	return newBlock
}

func isBlockValidZLH(newBlock, oldBlock *BlockZLH) bool {
	if oldBlock.IndexZLH+1 != newBlock.IndexZLH {
		return false
	}

	if oldBlock.HashZLH != newBlock.PrevHashZLH {
		return false
	}

	if calculateBlockHashZLH(newBlock) != newBlock.HashZLH {
		return false
	}

	return true
}

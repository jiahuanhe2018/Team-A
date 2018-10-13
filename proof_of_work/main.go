package main

import (
	"sync"
	"github.com/joho/godotenv"
	"log"
	"time"
	"github.com/davecgh/go-spew/spew"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"net/http"
	"encoding/json"
	"io"
	"strings"
	"fmt"
	"github.com/gorilla/mux"
	"os"
)

const EnvFileName = "myenv.env"
const Difficulty = 3

type BlockZLH struct {
	IndexZLH      int
	TimestampZLH  string
	BPMZLH        int
	HashZLH       string
	PrevHashZLH   string
	DifficultyZLH int
	NonceZLH      string
}

var BlockchainZLH []*BlockZLH

type MessageZLH struct {
	BPM int
}

var BlcMutexZLH = &sync.Mutex{}

func main() {
	err := godotenv.Load(EnvFileName)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		genesisBlock, _ := generateGenesisBlockZLH()
		BlcMutexZLH.Lock()
		BlockchainZLH = append(BlockchainZLH, genesisBlock)
		BlcMutexZLH.Unlock()
	}()

	log.Fatal(runZLH())
}

func runZLH() error {
	httpPort := os.Getenv("PORT")
	log.Println("Http监听端口：", httpPort)
	server := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        makeMuxRouterZLH(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func makeMuxRouterZLH() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlcZLH).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func handleGetBlcZLH(w http.ResponseWriter, r *http.Request) {
	blcBytes, err := json.MarshalIndent(BlockchainZLH, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(blcBytes))
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m MessageZLH

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		responseWithJSONZLH(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	BlcMutexZLH.Lock()
	newBlock := generateBlockZLH(BlockchainZLH[len(BlockchainZLH)-1], m.BPM)
	if isBlockValidZLH(
		newBlock, BlockchainZLH[len(BlockchainZLH)-1]) {
		BlockchainZLH = append(BlockchainZLH, newBlock)
		log.Println("最新区块链信息：")
		spew.Dump(BlockchainZLH)
	}
	BlcMutexZLH.Unlock()

	responseWithJSONZLH(w, r, http.StatusCreated, newBlock)
}

func responseWithJSONZLH(w http.ResponseWriter, r *http.Request,
	code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(payload, "", "   ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func generateBlockZLH(oldBlock *BlockZLH, BPM int) *BlockZLH {
	var nonce = 0
	var hash string

	newBlock := &BlockZLH{
		IndexZLH:      oldBlock.IndexZLH + 1,
		TimestampZLH:  time.Now().String(),
		BPMZLH:        BPM,
		HashZLH:       "",
		PrevHashZLH:   oldBlock.HashZLH,
		DifficultyZLH: Difficulty,
	}
	for {
		newBlock.NonceZLH = strconv.Itoa(nonce)
		hash = calculateHashZLH(newBlock)
		fmt.Printf("\r%d --- %s", nonce, hash)
		if isHashValidZLH(hash, Difficulty) {
			newBlock.HashZLH = hash
			fmt.Println()
			break
		}
		nonce ++
	}
	return newBlock
}

func isHashValidZLH(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

func calculateHashZLH(block *BlockZLH) string {
	record := strconv.Itoa(block.IndexZLH) + block.TimestampZLH + strconv.Itoa(
		block.BPMZLH) + block.PrevHashZLH + block.NonceZLH
	h := sha256.Sum256([]byte(record))
	return hex.EncodeToString(h[:])
}

func generateGenesisBlockZLH() (*BlockZLH, error) {
	genesisBlock := &BlockZLH{
		IndexZLH:      0,
		TimestampZLH:  time.Now().String(),
		BPMZLH:        0,
		HashZLH:       "",
		PrevHashZLH:   "",
		DifficultyZLH: Difficulty,
		NonceZLH:      "",
	}
	log.Println("创建创世区块：")
	spew.Dump(genesisBlock)
	return genesisBlock, nil
}

func isBlockValidZLH(newBlock, oldBlock *BlockZLH) bool {
	if oldBlock.IndexZLH+1 != newBlock.IndexZLH {
		return false
	}

	if oldBlock.HashZLH != newBlock.PrevHashZLH {
		return false
	}

	if calculateHashZLH(newBlock) != newBlock.HashZLH {
		return false
	}

	return true
}

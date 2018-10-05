package test

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	//"github.com/davecgh/go-spew/spew"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-host"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/libp2p/go-libp2p-net"
	"bufio"
)

const difficulty = 1

type Block struct {
	Index      int
	Timestamp  string
	Result     int
	Hash       string
	PrevHash   string
	Difficulty int
	Nonce      string
}

var blockchain []Block

var tempBlocks []Block

var work chan Block

var recv chan Block

func init() {
	blockchain = make([]Block, 0)
}

var m = &sync.RWMutex{}

func main() {
	go func() {
		makeGenesisBlock()
	}()
	go func() {
		worker()
	}()
	host, _ := makeBasicHost(6060, true, 6060)
	host.SetStreamHandler("/p2p/1.0.0", HandleStream)
	for {
		select {
		case b := <-work:
			appendAndBroadcast(b)
		case b := <-recv:
			handleTempBlock(b)
		}
	}
}

func HandleStream(s net.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go ReadData(rw)
	go WriteData(rw)
}
func WriteData(writer *bufio.ReadWriter) {
	
}
func ReadData(writer *bufio.ReadWriter) {
	
}

func makeBasicHost(listenPort int, secio bool, randseed int64) (host.Host, error) {

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	// Generate a key pair for this host. We will use it
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(priv),
	}

	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := basicHost.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)
	log.Printf("I am %s\n", fullAddr)
	if secio {
		log.Printf("Now run \"go run main.go -c chain -l %d -d %s -secio\" on a different terminal\n", listenPort+2, fullAddr)
	} else {
		log.Printf("Now run \"go run main.go -c chain -l %d -d %s\" on a different terminal\n", listenPort+2, fullAddr)
	}

	return basicHost, nil
}

func handleTempBlock(block Block) {
	//TODO
}

func broadcast(block Block) {

}

func appendAndBroadcast(b Block) {
	m.Lock()
	defer m.Unlock()
	blockchain = append(blockchain, b)
	broadcast(b)
}

func worker() {
NEXT_BLOCK:
	for {
		b := prepareBlock()
		l := len(blockchain)
		for i := 0; l == blen(); i++ {
			hex := fmt.Sprintf("%x", i)
			b.Nonce = hex
			if !isHashValid(calculateHash(b), b.Difficulty) {
				fmt.Println(calculateHash(b), " do more work!")
				continue
			} else {
				fmt.Println(calculateHash(b), " work done!")
				b.Hash = calculateHash(b)
				work <- b
				break NEXT_BLOCK
			}
		}
	}
}

func blen() int {
	m.RLock()
	defer m.RUnlock()
	return len(blockchain)
}

func makeGenesisBlock() {
	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{0, t.String(), 0, calculateHash(genesisBlock), "", difficulty, ""}
	m.Lock()
	defer m.Unlock()
	blockchain = append(blockchain, genesisBlock)
}

func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.Result) + block.PrevHash + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func prepareBlock() Block {
	m.RLock()
	defer m.RUnlock()

	t := time.Now()
	oldBlock := blockchain[len(blockchain)-1]
	var newBlock Block
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Result = 1
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = difficulty
	return newBlock
}


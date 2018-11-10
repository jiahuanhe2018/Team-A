package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"github.com/36dfinity/tbls"
	"strings"
	"time"
	"sort"
)

var (
	T          = 2
	N          = 3
	seed       = "bls_chain"
	listener   *net.UDPConn
	addr       *net.UDPAddr
	all        map[int]*Node
	sk         bls.SecretKey
	pk         *bls.PublicKey
	aggrSk     bls.SecretKey
	aggrPk     *bls.PublicKey
	poly       []bls.SecretKey
	commitment string
	wg         sync.WaitGroup
)

func main() {
	InitNetwork(os.Args[1])
	go read(listener)
	InitBls()

	//TODO
	time.Sleep(10 * time.Second)
	fmt.Println("allCommitmentsReady", allCommitmentsReady())
	fmt.Println("allSharesReady", allSharesReady())

	aggrSk = getAggrSk()
	aggrPk = aggrSk.GetPublicKey()
	fmt.Println("aggrSk", aggrSk, "aggrPk", aggrPk)

	sign := sign(seed)
	for i := 0; i < 30; i++ {
		time.Sleep(time.Second)
		broadcast("sign " + sign)
	}
	fmt.Println("allSignsReady", allSignsReady())

	gs := seed
	for _, node := range all {
		gs += node.sign
		var shareSK bls.SecretKey
		shareSK.SetByCSPRNG()
		shareSK.SetHexString(node.share)

	}

	b := make([]byte, 1)
	os.Stdin.Read(b)
}

type AtomicBoolean struct {
	b bool
	*sync.RWMutex
}

func (a *AtomicBoolean) Get() bool {
	a.RLock()
	return a.b
}

func (a *AtomicBoolean) Set(b bool) {
	a.Lock()
	a.b = b
}

type Node struct {
	udpAddr    *net.UDPAddr
	commitment string // broadcast
	share      string //unicast
	sign       string
}

func InitNetwork(port string) {
	p, _ := strconv.Atoi(port)
	addr = &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: p}
	l, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	listener = l
	all = make(map[int]*Node, N)
	for i := 0; i < N; i++ {
		n := &Node{}
		n.udpAddr = &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8964 + i}
		all[8964+i] = n
	}
}

func InitBls() {
	bls.Init(bls.CurveFp382_1)
	sk.SetByCSPRNG()
	pk = sk.GetPublicKey()
	poly = sk.GetMasterSecretKey(T)
	pks := make([]string, 0, T)
	for _, p := range poly {
		pks = append(pks, p.GetPublicKey().GetHexString())
	}
	commitment = strings.Join(pks, ":")

	// TODO: 可靠传输保证参与初始化的 n 个节点收到 commitment
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		// TODO: 用 tendermint 私钥签名并将签名一起发送
		broadcast("commitment " + commitment)
	}
	shares := make([]string, 0)
	for i := 0; i < N; i++ {
		var s bls.SecretKey
		s.SetByCSPRNG()
		var id bls.ID
		id.SetLittleEndian([]byte{1, 2, 3, 4, 5, byte(i)})
		s.Set(poly, &id)
		shares = append(shares, fmt.Sprintf("share %s", s.GetHexString()))
	}
	for k := 0; k < 10; k++ {
		time.Sleep(time.Second)
		for i, v := range shares {
			addr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8964 + i}
			// TODO:eth 椭圆曲线加密
			unicast(v, addr)
		}
	}
}

func read(conn *net.UDPConn) {
	for {
		data := make([]byte, 4096)
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		inputs := strings.Split(string(data[:n]), " ")
		port := remoteAddr.Port
		switch inputs[0] {
		case "commitment":
			// TODO: 用 tendermint 的公钥验签
			all[port].commitment = inputs[1]
			splits := strings.Split(inputs[1], ":")
			for _, split := range splits {
				var pk bls.PublicKey
				pk.SetHexString(split)

			}
		case "share":
			// TODO:eth 的椭圆曲线解密
			all[port].share = inputs[1]
		case "sign":
			// TODO: commitment 验签
			all[port].sign = inputs[1]
		}
	}
}

// TODO: reflect
func allCommitmentsReady() bool {
	for _, v := range all {
		if v.commitment == "" {
			return false
		}
	}
	return true
}

func allSharesReady() bool {
	for _, v := range all {
		if v.share == "" {
			return false
		}
	}
	return true
}

func allSignsReady() bool {
	for _, v := range all {
		if v.sign == "" {
			return false
		}
	}
	return true
}

func unicast(msg string, dest *net.UDPAddr) {
	listener.WriteToUDP([]byte(msg), dest)
}

func broadcast(msg string) {
	for _, v := range all {
		unicast(msg, v.udpAddr)
	}
}

func sign(msg string) string {
	sign := aggrSk.Sign(msg)
	return sign.GetHexString()
}

func getAggrSk() bls.SecretKey {
	keys := make([]int, 0)
	for k, _ := range all {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	var aggrSk bls.SecretKey
	aggrSk.SetByCSPRNG()
	for _, k := range keys {
		var s bls.SecretKey
		s.SetByCSPRNG()
		node := all[k]
		s.SetHexString(node.share)
		aggrSk.Add(&s)
	}
	return aggrSk
}

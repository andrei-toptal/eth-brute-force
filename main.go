package main

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var iterations uint64

func getAddressForPrivateKey(privateKey *ecdsa.PrivateKey) common.Address {
	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		panic("failed to get public key from the given private key")
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA)
}

func bruteForceRoutine(ethAddress []byte, doneCh chan bool) {
	for {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		address := getAddressForPrivateKey(privateKey)
		if bytes.Equal(address.Bytes(), ethAddress) {
			privateKeyBytes := crypto.FromECDSA(privateKey)
			fmt.Printf("Private key found: %s\n", hexutil.Encode(privateKeyBytes))
			if err := os.WriteFile("private_key", privateKeyBytes, 0644); err != nil {
				panic(err)
			}
			doneCh <- true
		}
		atomic.AddUint64(&iterations, 1)
	}
}

func iterationsWatcher() {
	var prevIterations uint64
	oneSecondTimer := time.NewTimer(time.Second)
	for range oneSecondTimer.C {
		currentIterations := atomic.LoadUint64(&iterations)
		fmt.Printf("Speed: %d iterations per second\n", currentIterations-prevIterations)
		prevIterations = currentIterations
		oneSecondTimer.Reset(time.Second)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: eth-brute-force ethereum_hex_address")
		return
	}

	ethAddressHex := os.Args[1]
	ethAddress, err := hexutil.Decode(ethAddressHex)
	if err != nil {
		fmt.Printf("wrong ethereum address: %v\n", err)
	}

	fmt.Printf("Searching private key for %s\n", hexutil.Encode(ethAddress))
	threads := runtime.NumCPU()
	fmt.Printf("CPU threads: %d\n", threads)

	go iterationsWatcher()

	done := make(chan bool)
	for i := 0; i < threads; i++ {
		go bruteForceRoutine(ethAddress, done)
	}
	<-done
}

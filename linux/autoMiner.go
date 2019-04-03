package main

import (
	"fmt"
	"github.com/usechain/go-usedrpc"
	"math/big"
	"time"
	"strconv"
	"github.com/peterh/liner"
	"strings"

)

func main() {
	rpc := usedrpc.NewUseRPC("http://127.0.0.1:8848")
	accounts, _ := rpc.UseAccounts()

	if len(accounts) == 0{
		fmt.Println("No account!!!!!!")
		return
	} else {
		fmt.Println("Please choose one account to mine, enter account number")
		num, err := reader()
		if err != nil {
			return
		}
		accountNum,err := strconv.Atoi(num)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Enter account password")
		passwd, err := reader()
		if err != nil {
			return
		}

		unlock, err := rpc.UnlockAccount(accounts[accountNum], passwd)
		if err != nil {
			fmt.Println("unlock account ", err)
		}
		fmt.Println("unlock account", unlock)
	}

	queryCert, err := rpc.GetCertifications(accounts[0])
	if err != nil {
		fmt.Println("GetCertifications", err)
	}
	if queryCert != "0x1" {
		//sendCreditRegister
		tx := usedrpc.T {
			From: accounts[0],
			To:   "0xfffffffffffffffffffffffffffffffff0000001",
			Value: big.NewInt(0),
			Data:  "",
			Gas:4000000,
			GasPrice: big.NewInt(40000000000),
		}
		txhash, err := rpc.SendCreditRegisterTransaction(tx)
		if err != nil {
			fmt.Println("SendCreditRegister", err)
		}
		fmt.Println("The SendCreditRegister Transaction hash:", err, txhash)
		time.Sleep(5*time.Second)

		queryCert2, err := rpc.GetCertifications(accounts[0])
		if queryCert2 == "0x1" {
			miner(accounts,rpc)
		}
	} else if queryCert == "0x1" {
		miner(accounts,rpc)
	}
}

func miner(accounts []string, rpc *usedrpc.UseRPC) {
	n := new(big.Int)
	n, ok := n.SetString("51000000000000000000", 0)
	if ok {
		//fmt.Println("register miner use value: ", n)
	} else {
		fmt.Println("SetString: error")
	}

	//query miner
	queryMiner, err := rpc.UseIsMiner(accounts[0], "latest")
	fmt.Println("queryMiner register result", queryMiner)
	if queryMiner {
		//miner start
		err = rpc.MinerStart()
		if err != nil {
			fmt.Println("Miner start failed, please open your rpcapi 'miner'", err)
			return
		}
		return
	}

	// miner register
	var registerMinnerTX usedrpc.T
	registerMinnerTX = usedrpc.T {
		From: accounts[0],
		To:   "0xfffffffffffffffffffffffffffffffff0000002",
		Value: n,
		Data:  "0x819f163a",
		GasPrice: big.NewInt(40000000000),
	}
	res, err := rpc.UseSendTransaction(registerMinnerTX)
	fmt.Println("Send  Miner register transaction hash:", err, res)
	time.Sleep(2 * time.Second)

	//miner start
	err = rpc.MinerStart()
	if err != nil {
		fmt.Println("Miner start failed, please open your rpcapi 'miner'", err)
		return
	}
	fmt.Println("Miner start: true")

}

func reader() (string, error) {
	Prompter := liner.NewLiner()
	defer Prompter.Close()
	Prompter.SetCtrlCAborts(true)
	active := true
	for active {
		typed, err := Prompter.Prompt("> ")
		if err != nil {
			return "", err
		}
		// If at least a character is typed
		var arr []string
		if arr = strings.Fields(typed); len(arr) > 0 {
			active = false
			return arr[0], nil
		}
	}
	return "", nil
}

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
	time.Sleep(2 * time.Second)
	accounts, err := rpc.UseAccounts()
	if err != nil {
		fmt.Println("UseAccounts", err)
	}
	var accountNum int
	if len(accounts) == 0{
		fmt.Println("No account!!!!!!")
		return
	} else {
		fmt.Println("Please choose one account to mine, enter account number>>>")
		num, err := reader()
		if err != nil {
			return
		}
		accountNum,err = strconv.Atoi(num)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Enter account password>>>")
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

	queryCert, err := rpc.GetCertifications(accounts[accountNum])
	if err != nil {
		fmt.Println("GetCertifications", err)
	}
	fmt.Println("GetCertifications: ", queryCert)
	if queryCert != "0x1" {
		//sendCreditRegister
		tx := usedrpc.T {
			From: accounts[accountNum],
			To:   "UmixYUgBHA9vJj47myQKn8uZAm4an7zyYJ8",
			Value: big.NewInt(0),
			Data:  "",
			Gas:4000000,
			GasPrice: big.NewInt(40000000000),
		}
		txhash, err := rpc.SendCreditRegisterTransaction(tx, false)
		if err != nil {
			fmt.Println("SendCreditRegister", err)
		}

		fmt.Println("The SendCreditRegister Transaction hash:", txhash)
		fmt.Println("Please wait 30 second for committee to confirm your credit register......")
		time.Sleep(60 * time.Second)

		queryCert2, err := rpc.GetCertifications(accounts[accountNum])
		if queryCert2 == "0x1" {
			miner(accounts,rpc,accountNum)
		} else {
			fmt.Println("Your credit register still unconfirmed, this program shudown now, please try again later")
		}
	} else if queryCert == "0x1" {
		miner(accounts,rpc,accountNum)
	}
}

func miner(accounts []string, rpc *usedrpc.UseRPC, accountNum int) {
	n := new(big.Int)
	n, ok := n.SetString("50000000000000000000", 0)
	if ok {
		//fmt.Println("register miner use value: ", n)
	} else {
		fmt.Println("SetString: error")
	}

	//query isMiner
	queryMiner, err := rpc.UseIsMiner(accounts[accountNum], "latest")
	fmt.Println("queryMiner register result", queryMiner)
	if queryMiner {
		//miner start
		err = rpc.MinerStart()
		if err != nil {
			fmt.Println("Miner start failed, please open your rpcapi 'miner'", err)
			return
		}
		return
	} else {
		//if punished
		isPunished, err := rpc.UseIsPunishedMiner(accounts[accountNum], "latest")
		if err != nil {
			fmt.Println("UseIsPunishedMiner", err)
		}
		if  isPunished {
			// miner register
			var unregisterMinnerTX usedrpc.T
			unregisterMinnerTX = usedrpc.T {
				From: accounts[accountNum],
				To:   "UmixYUgBHA9vJj47myQKn8uZAm4anEfrG78",
				Data:  "0x6d3a3f8d",
				GasPrice: big.NewInt(40000000000),
			}
			res, err := rpc.UseSendTransaction(unregisterMinnerTX)
			if err != nil {
				fmt.Println("Send Miner unregister transaction error", err)
				return
			}
			fmt.Println("Send  Miner unregister transaction hash:", res)
			time.Sleep(2 * time.Second)
		}
	}

	// miner register
	var registerMinnerTX usedrpc.T
	registerMinnerTX = usedrpc.T {
		From: accounts[accountNum],
		To:   "UmixYUgBHA9vJj47myQKn8uZAm4anEfrG78",
		Value: n,
		Data:  "0x819f163a",
		GasPrice: big.NewInt(40000000000),
	}
	res, err := rpc.UseSendTransaction(registerMinnerTX)
	if err != nil {
		fmt.Println("Send Miner register transaction error", err)
	}
	fmt.Println("Send  Miner register transaction hash:", res)
	time.Sleep(10 * time.Second)

	//miner start
	err = rpc.MinerStart()
	if err != nil {
		fmt.Println("Miner start failed, please open your rpcapi 'miner'", err)
		return
	}
	time.Sleep(2 * time.Second)

	mining, err := rpc.UseMining()
	if err != nil {
		fmt.Println("useMing error", err)
	}
	fmt.Println("Miner start: ", mining)
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



package main

import (
	ot "backend-go/onetoken-go-sdk/onetoken"
	"fmt"
)

func main() {
	contract := "okex/btc.usdt"
	//contract1 := "huobip/btc.usdt"

	//go getTick(contract, tickData)
	tickData, err := getTick(contract)
	zhubiData, err := getZhubi(contract)
	candleData, err := getCandle(contract)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("connecting.....")
	//time.Sleep(time.Millisecond * 500)
	/*for {
		time.Sleep(time.Second * 2)
		fmt.Printf("received data: %s\n", data)
	}*/

	for {
		select {
		case tickMsg := <-tickData:
			//time.Sleep(time.Second * 1)
			fmt.Printf("received tickdata: %s\n", string(tickMsg))
		case zhubiMsg := <-zhubiData:
			fmt.Printf("received zhubidata: %s\n", string(zhubiMsg))
		case candleMsg := <-candleData:
			fmt.Printf("received candledata: %s\n", string(candleMsg))
		}
	}
}

//获取实时tick数据
func getTick(contract string) (chan []byte, error) {

	tickData, err := ot.ContractsTick(contract)
	return tickData, err
}

//获取实时逐笔交易数据
func getZhubi(contract string) (chan []byte, error) {

	zhubiData, err := ot.ContractsZhubi(contract)
	return zhubiData, err
}

//获取实时candle数据
func getCandle(contract string) (chan []byte, error) {

	candleData, err := ot.ContractsCandle(contract)
	return candleData, err
}

package main

import (
	ot "backend-go/onetoken-go-sdk/onetoken"
	"fmt"
)

var acc1 = ot.Account{
	OtKey:    "",
	OtSecret: "",
	Accounts: "",
}

func main() {
	msg, err := getInfo()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(msg)
	fmt.Println("hehe")
}

//获取账户信息
func getInfo() (string, error) {
	contract := "okex/btc.usdt"
	infomsg, err := acc1.GetAccountInfo(contract)
	if err != nil {
		return "", err
	}
	return infomsg, nil
}

//查询订单信息
func getOrderList() (string, error) {
	contract := "okex/btc.usdt"
	state := "end"
	orderListmsg, err := acc1.GetAccountOrder(contract, state)
	if err != nil {
		return "", err
	}
	return orderListmsg, nil
}

//获取最近成交记录
func getTrans() (string, error) {
	contract := "okex/btc.usdt"
	count := "100"
	transmsg, err := acc1.GetAccountTrans(contract, count)
	if err != nil {
		return "", err
	}
	return transmsg, nil
}

//创建订单
func createOrder() (string, error) {
	contract := "okex/btc.usdt"
	data := map[string]interface{}{
		"contract": contract,
		"bs":       "b",
		"price":    3100,
		"amount":   0.03,
	}
	ordermsg, err := acc1.CreateOrder(contract, data)
	if err != nil {
		return "", err
	}
	return ordermsg, nil
}

// 取消订单
func deleteOrder() (string, error) {
	exchangeOid := "okex/btc.usdt-z2b3s1kqlf8lx9iow8scn1e8jxr"
	dordermsg, err := acc1.DeleteOrder(exchangeOid)
	if err != nil {
		return "", err
	}
	return dordermsg, nil
}

//取消所有订单
func deleteAllOrder() (string, error) {
	contract := "okex/btc.usdt"
	dAllordermsg, err := acc1.DeleteAllOrder(contract)
	if err != nil {
		return "", err
	}
	return dAllordermsg, nil
}

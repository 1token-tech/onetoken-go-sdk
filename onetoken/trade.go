package onetoken

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

//Account 账户struct
type Account struct {
	OtKey    string
	OtSecret string
	Accounts string
}

var urlHeader string

func init() {
	urlHeader = "https://1token.trade/api/v1/trade"
}

//GetAccountInfo 获取账户信息
func (acc *Account) GetAccountInfo(contract string) (string, error) {
	verb := "GET"
	path := "/" + acc.Accounts + "/info"
	nonce := getNonce()
	signData := getSign(acc.OtSecret, verb, path, nonce, "")
	url := urlHeader + path
	//log.Println(signUri)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	getHander(nonce, acc.OtKey, signData, req)

	resp, err := client.Do(req)
	//resp, err := http.Get("https://1token.trade/api/v1/trade/info")

	if err != nil {
		log.Println(err)
		return "", err
	}

	defer resp.Body.Close() //关闭
	body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	return string(body), nil

}

//GetAccountOrder 查询订单信息
func (acc *Account) GetAccountOrder(contract, state string) (string, error) {
	verb := "GET"
	path := "/" + acc.Accounts + "/orders"
	nonce := getNonce()
	signData := getSign(acc.OtSecret, verb, path, nonce, "")
	url := urlHeader + path + "?contract=" + contract + "&state=" + state
	//log.Println(signUri)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	getHander(nonce, acc.OtKey, signData, req)

	resp, err := client.Do(req)
	//resp, err := http.Get("https://1token.trade/api/v1/trade/info")

	if err != nil {
		log.Println(err)
		return "", err
	}

	defer resp.Body.Close() //关闭
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

//GetAccountTrans 获取最近成交记录
func (acc *Account) GetAccountTrans(contract, count string) (string, error) {
	verb := "GET"
	path := "/" + acc.Accounts + "/trans"
	nonce := getNonce()
	signData := getSign(acc.OtSecret, verb, path, nonce, "")
	url := urlHeader + path + "?contract=" + contract + "&count=" + count
	//log.Println(signUri)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	getHander(nonce, acc.OtKey, signData, req)

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return "", err
	}

	defer resp.Body.Close() //关闭
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

//CreateOrder 创建订单
func (acc *Account) CreateOrder(contract string, data map[string]interface{}) (string, error) {
	verb := "POST"
	path := "/" + acc.Accounts + "/orders"

	datas, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	nonce := getNonce()
	signData := getSign(acc.OtSecret, verb, path, nonce, string(datas))
	url := urlHeader + path
	//log.Println(signUri)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(datas)))
	if err != nil {
		return "", err
	}
	getHander(nonce, acc.OtKey, signData, req)

	resp, err := client.Do(req)
	/*resp, err := http.Post("https://1token.trade/api/v1/trade/orders",
	                             "application/x-www-form-urlencode",// Content-Type post请求必须设置
		                         strings.NewReader("name=abc"))*/

	if err != nil {
		log.Println(err)
		return "", err
	}

	defer resp.Body.Close() //关闭
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

//DeleteOrder 取消订单
func (acc *Account) DeleteOrder(exchangeOid string) (string, error) {
	verb := "DELETE"
	path := "/" + acc.Accounts + "/orders"
	nonce := getNonce()
	signData := getSign(acc.OtSecret, verb, path, nonce, "")
	url := urlHeader + path + "?exchange_oid=" + exchangeOid
	//log.Println(signUri)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return "", err
	}
	getHander(nonce, acc.OtKey, signData, req)

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return "", err
	}

	defer resp.Body.Close() //关闭
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

//DeleteAllOrder 取消所有订单
func (acc *Account) DeleteAllOrder(contract string) (string, error) {
	verb := "DELETE"
	path := "/" + acc.Accounts + "/orders/all"
	nonce := getNonce()
	signData := getSign(acc.OtSecret, verb, path, nonce, "")
	url := urlHeader + path + "?contract=" + contract
	//log.Println(signUri)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return "", err
	}
	getHander(nonce, acc.OtKey, signData, req)

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return "", err
	}

	defer resp.Body.Close() //关闭
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func getSign(key, verb, path, nonce, data string) string {
	hmac := hmac.New(sha256.New, []byte(key))
	var dataStr string
	if data == "" {
		dataStr = ""
	} else {
		dataStr = data
	}
	datas := verb + path + nonce + dataStr
	hmac.Write([]byte(datas))

	b := hmac.Sum(nil)
	rsp := hex.EncodeToString(b)
	//log.Println(rsp)
	return rsp
}

func getNonce() string {
	timestamp := time.Now().Unix()
	return strconv.Itoa(int(timestamp * 1000000))
}

func getHander(nonce, otKey, signData string, req *http.Request) {
	req.Header.Set("Api-Nonce", nonce)
	req.Header.Set("Api-Key", otKey)
	req.Header.Set("Api-Signature", signData)
	req.Header.Set("Content-Type", "application/json")
}

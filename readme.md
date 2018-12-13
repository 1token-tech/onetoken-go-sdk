# OneToken SDK golang版本

## SDK介绍

OneToken SDK 是一个api接口，用户可以通过这个接口获取tick、zhubi、candle行情信息和处理订单。它包含ws API和restful API接口，
ws API接口用于获取行情信息，restful API接口用于订单的处理。用户可以根据自己的喜好直接使用API。

## 支持的交易所

* bitmex
* okex
* binance
* bithumb
* huobi.pro
* bitfinex
* bitstar
* bittrex
* poloniex
* gate
* exx
* coinegg

其他交易所的支持正在开发中

## golang版本需求

go1.10.5及以上版本

## SDK的安装

go get github.com/1token-trade/onetoken-go-sdk/onetoken  
或下载到本地后 go install ...

## 名词解释

### tick

 买卖挡位数据  
 包含属性：`time, price, volume, bids, asks, contract, last, exchange_time, amount`  
 属性含义见如下json

 ```json
 {
    "uri":"single-tick-verbose",
    "data":
    {
         "asks":
         [
             {"price": 9218.5, "volume": 1.7371},
             ...
         ],
         "bids":
         [
             {"price": 9218.4, "volume": 0.81871728},
             ...
         ],
         "contract": "bitfinex/btc.usd",
         "last": 9219.3,  // 最新成交价
         "time": "2018-05-03T16:16:41.630400+08:00",  // 1token的系统时间 ISO 格式 带时区
         "exchange_time": "2018-05-03T16:16:41.450400+08:00",  // 交易所给的时间 ISO 格式 带时区
         "amount": 16592.4,  //成交额 (CNY)
         "volume": 0.3   // 成交量
   }
}
```

### zhubi

成交记录  
包含属性：`amount, bs, contract， exchange_time， price， time`  
 属性含义见如下json

 ```json
{
    "uri":"single-zhubi-verbose",
    "data":
    [
        {
            "amount": 0.21,
            "bs": "s",
            "contract": "bitfinex/btc.usd",
            "exchange_time": "2018-05-03T08:14:20.307000+00:00",
            "price": 9231.8,
            "time": "2018-05-03T16:14:20.541068+08:00"
        }
    ]
}
```

订阅逐笔数据, 如果需要请求多个contract的逐笔数据， 可以在同一个websocket里面发送多个subscribe-single-zhubi-verbose的请求, 每个请求带着不同的contract

### candle

k线数据  
包含属性：`amount, close， high， low，open，volume，contract，duration, time`  
属性含义见如下json

 ```json
{
    "amount": 16592.4, //成交量
    "close": 9220.11,
    "high": 9221,
    "low": 9220.07,
    "open": 9220.07,
    "volume": 0.3, //成交额
    "contract": "huobip/btc.usdt",
    "duration": "1m",
    "time": "2018-05-03T07:30:00Z" // 时间戳 isoformat
}
```

### otKey、otSecret

onetoken上生成的“我的OT Key”

### accounts

交易所开通的账户 “交易所/账户名”

### contract

交易所和交易对 "huobip/btc.usdt"

### state

订单状态 active和end两种基本状态

### count

订单成交记录数量

### data

创建订单的参数  
例如：

```golang
data := map[string]interface{}{
    "contract": acc.contract,
    "bs":       "b",
    "price":    3100,
    "amount":   0.03,
}
```

## SDK的使用

详细例子见examples  
用封装的restfull接口交易：

```golang
//创建订单
var acc1 = ot.Account{
    OtKey:    "",
    OtSecret: "",
    Accounts: "",
}
contract := "okex/btc.usdt"
data := map[string]interface{}{
    "contract": acc.contract,
    "bs":       "b",
    "price":    3100,
    "amount":   0.03,
}
ordermsg, err := acc1.CreateOrder(contract, data)
if err != nil {
    return "", err
}
```

用封装的websocket接口获取行情：

```golang
//获取实时tick数据
contract := "okex/btc.usdt"
tickData, err := quote.ContractsTick(contract)
if err != nil {
    fmt.Println("error:", err)
    return
}
```
package onetoken

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var lock *sync.Mutex

//var conn *websocket.Conn

func init() {
	lock = new(sync.Mutex)
	//lock = &sync.Mutex{}
	//conn = connecting()
}

var addr = flag.String("addr", "1token.trade", "http service address")

//ContractsTick tick信息ws接口
func ContractsTick(clientRequest string) (chan []byte, error) {
	conn := tickConnecting()
	count := 0
	go tickWrite(conn, clientRequest)
	//go zhubi(conn, clientRequest)
	go heartBeatLoop(conn)

	messageData := make(chan []byte, 10)
	go func() {
		for {
			if conn == nil {
				if count > 100 {
					fmt.Println("connect outtime")
					conn.Close()
					count = 0
					break
				}
				count++
				time.Sleep(2 * time.Second)
				fmt.Println("waiting for connecting....")
				conn = tickConnecting()
				go tickWrite(conn, clientRequest)
				go heartBeatLoop(conn)
			} else {
				_, message, err := conn.ReadMessage()
				if err != nil {
					fmt.Println("read:", err)
					conn = tickConnecting()
					go tickWrite(conn, clientRequest)
					go heartBeatLoop(conn)
					fmt.Println("reconnect....")
					continue
				} else {
					count = 0
				}

				messageData <- message
			}
			//fmt.Printf("received: %s\n", *messageData)
		}
	}()
	return messageData, nil
}

//ContractsZhubi 逐笔信息ws接口
func ContractsZhubi(clientRequest string) (chan []byte, error) {
	conn := tickConnecting()
	count := 0
	go zhubiWrite(conn, clientRequest)
	go heartBeatLoop(conn)

	messageData := make(chan []byte, 10)
	go func() {
		for {
			if conn == nil {
				if count > 100 {
					fmt.Println("connect outtime")
					conn.Close()
					count = 0
					break
				}
				count++
				time.Sleep(2 * time.Second)
				fmt.Println("waiting for connecting....")
				conn = tickConnecting()
				go zhubiWrite(conn, clientRequest)
				go heartBeatLoop(conn)
			} else {
				_, message, err := conn.ReadMessage()
				if err != nil {
					fmt.Println("read:", err)
					conn = tickConnecting()
					go zhubiWrite(conn, clientRequest)
					go heartBeatLoop(conn)
					fmt.Println("reconnect....")
					continue
				} else {
					count = 0
				}

				messageData <- message
			}
			//fmt.Printf("received: %s\n", *messageData)
		}
	}()
	return messageData, nil
}

//ContractsCandle candle信息ws接口
func ContractsCandle(clientRequest string) (chan []byte, error) {
	conn := candleConnecting()
	count := 0
	go candleWrite(conn, clientRequest)
	go heartBeatLoop(conn)

	// var messageData chan []byte
	messageData := make(chan []byte, 10)
	go func() {
		for {
			if conn == nil {
				if count > 100 {
					fmt.Println("connect outtime")
					conn.Close()
					count = 0
					break
				}
				count++
				time.Sleep(2 * time.Second)
				fmt.Println("waiting for connecting....")
				conn = candleConnecting()
				go candleWrite(conn, clientRequest)
				go heartBeatLoop(conn)
			} else {
				_, message, err := conn.ReadMessage()
				if err != nil {
					fmt.Println("read:", err)
					conn = candleConnecting()
					go candleWrite(conn, clientRequest)
					go heartBeatLoop(conn)
					fmt.Println("reconnect....")
					continue
				} else {
					count = 0
				}

				messageData <- message
			}
			//fmt.Printf("received: %s\n", *messageData)
		}
	}()
	return messageData, nil
}

func heartBeatLoop(conn *websocket.Conn) {
	for {
		if conn == nil {
			time.Sleep(1 * time.Second)
			break
		} else {
			p := make(map[string]string)
			p["uri"] = "ping"
			j, _ := json.Marshal(p)
			time.Sleep(time.Second * 20)
			err := conn.WriteMessage(websocket.TextMessage, j)
			if err != nil {
				fmt.Println("writeHeartErr:", err)
				break
			}
		}

	}
}

func tickWrite(conn *websocket.Conn, clientRequest string) {
	for {
		if conn == nil {
			time.Sleep(1 * time.Second)
			break
		} else {
			request := make(map[string]string)
			request["uri"] = "subscribe-single-tick-verbose"
			request["contract"] = clientRequest

			requestm, _ := json.Marshal(request)
			time.Sleep(time.Second * 2)
			//conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))

			lock.Lock()
			err := conn.WriteMessage(websocket.TextMessage, requestm)
			lock.Unlock()
			if err != nil {
				fmt.Println("write tick err:", err)
				break
			}
		}
	}
}

func zhubiWrite(conn *websocket.Conn, clientRequest string) {
	for {
		if conn == nil {
			time.Sleep(1 * time.Second)
			break
		} else {
			request := make(map[string]string)
			request["uri"] = "subscribe-single-zhubi-verbose"
			request["contract"] = clientRequest

			requestm, _ := json.Marshal(request)
			time.Sleep(time.Second * 2)
			//conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
			lock.Lock()
			err := conn.WriteMessage(websocket.TextMessage, requestm)
			lock.Unlock()
			if err != nil {
				fmt.Println("write zhubi err:", err)
				break
			}
		}
	}
}

func candleWrite(conn *websocket.Conn, clientRequest string) {
	for {
		if conn == nil {
			time.Sleep(1 * time.Second)
			break
		} else {
			request := make(map[string]string)
			request["duration"] = "1m"
			request["contract"] = clientRequest

			requestm, _ := json.Marshal(request)
			time.Sleep(time.Second * 2)
			//conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
			lock.Lock()
			err := conn.WriteMessage(websocket.TextMessage, requestm)
			lock.Unlock()
			if err != nil {
				fmt.Println("write candle err:", err)
				break
			}
		}
	}
}

func candleConnecting() (conn *websocket.Conn) {
	u := url.URL{Scheme: "wss", Host: *addr, Path: "/api/v1/ws/candle"}
	var dialer *websocket.Dialer
	//conn, _, err := dialer.Dial("ws://localhost:9876/ws", nil)
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return conn
}

func tickConnecting() (conn *websocket.Conn) {
	u := url.URL{Scheme: "wss", Host: *addr, Path: "/api/v1/ws/tick"}
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return conn
}

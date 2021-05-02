package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket" //这里使用的是 gorilla 的 websocket 库
)

func main() {
	//创建一个拨号器，也可以用默认的 websocket.DefaultDialer
	dialer := websocket.Dialer{}
	//向服务器发送连接请求，websocket 统一使用 ws://，默认端口和http一样都是80
	connect, r, err := dialer.Dial("wss://im-api-vip6-v2.easemob.com/ws/032/xvrhfd2j/websocket", nil)

	log.Println(r.Body)

	if nil != err {
		log.Println(err)
		return
	}
	//离开作用域关闭连接，go 的常规操作
	defer connect.Close()

	//定时向客户端发送数据
	//go tickWriter(connect)

	//启动数据读取循环，读取客户端发送来的数据
	for {
		//从 websocket 中读取数据
		//messageType 消息类型，websocket 标准
		//messageData 消息数据
		messageType, messageData, err := connect.ReadMessage()
		if nil != err {
			log.Println(err)
			break
		}
		switch messageType {
		case websocket.TextMessage: //文本数据
			fmt.Println(string(messageData))
			if string(messageData) == "o" {
				//connect.WriteJSON()
				log.Println(0)
				err := connect.WriteMessage(websocket.TextMessage, []byte(`["CAASPAoOY3gtZGV2I2N4c3R1ZHkSCDc1MDUwMDQ4GgtlYXNlbW9iLmNvbSITd2ViaW1fMTYxOTk3NTc4OTIxMhqFASR0JFlXTXRQWE9WWEt0cUVldXNWYm11ZzRkQU9LS2xFNEQtNlJIa3NZZ0p4U3NiOGpMS2tLcHcyV29SNlkyY3M1b05Cc0w2QXdNQUFBRjVMaFdOeUFCUEdnRGVrdEdzc295TEFDODctNFZFakRSS1FaZ05tMXNrUDR6X2RJMVIwOTNiOXdAA0rAAQgQEgUzLjAuMCgAMABKDTE2MTk5NzU3ODkyMTJiBXdlYmltahN3ZWJpbV8xNjE5OTc1Nzg5MjEycoUBJHQkWVdNdFBYT1ZYS3RxRWV1c1ZibXVnNGRBT0tLbEU0RC02Ukhrc1lnSnhTc2I4akxLa0twdzJXb1I2WTJjczVvTkJzTDZBd01BQUFGNUxoV055QUJQR2dEZWt0R3Nzb3lMQUM4Ny00VkVqRFJLUVpnTm0xc2tQNHpfZEkxUjA5M2I5d1AAWAA="]`))
				if nil != err {
					log.Println(err)
					break
				}
			}
		case websocket.BinaryMessage: //二进制数据
			fmt.Println(messageData)
		case websocket.CloseMessage: //关闭
		case websocket.PingMessage: //Ping
		case websocket.PongMessage: //Pong
		default:

		}
	}

	time.Sleep(100 * time.Second)
}

func tickWriter(connect *websocket.Conn) {
	for {
		//向客户端发送类型为文本的数据
		err := connect.WriteMessage(websocket.TextMessage, []byte("from client to server"))
		if nil != err {
			log.Println(err)
			break
		}
		//休息一秒
		time.Sleep(time.Second)
	}
}

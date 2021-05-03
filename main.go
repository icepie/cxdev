package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/icepie/cxdev/client"

	"github.com/gorilla/websocket"
)

func main() {

	cxuser, err := client.NewCXUser("B19071121", "", "1283")
	if err != nil {
		log.Println(err)
	}

	log.Println(cxuser.RealName)
	log.Println(cxuser.Token)
	log.Println(cxuser.Cooikes)
	log.Println(cxuser.TUID)
	log.Println(cxuser.IMToken)

	//创建一个拨号器，也可以用默认的 websocket.DefaultDialer
	dialer := websocket.Dialer{}
	//向服务器发送连接请求，websocket 统一使用 ws://，默认端口和http一样都是80
	connect, r, _ := dialer.Dial("wss://im-api-vip6-v2.easemob.com/ws/032/xvrhfd2j/websocket", nil)

	log.Println(r.Status)

	// s := "46447381"

	// data, err := hex.DecodeString(s)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("% x", data)

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
				err := connect.WriteMessage(websocket.TextMessage, []byte(`["CAASPAoOY3gtZGV2I2N4c3R1ZHkSCDc1MDUwMDQ4GgtlYXNlbW9iLmNvbSITd2ViaW1fMTYxOTk4NDQzOTEzORqFASR0JFlXTXQ3Y0pQNkt0cUVldTNVNmZxWVB0VExhS2xFNEQtNlJIa3NZZ0p4U3NiOGpMS2tLcHcyV29SNlkyY3M1b05Cc0w2QXdNQUFBRjVMaG9ST2dCUEdnRDBaS2ZHdjV4Nm1zUTJXV0pfVkJiczVwRDNLVE8zZnV4NTU3d2lkTHpadlFAA0rAAQgQEgUzLjAuMCgAMABKDTE2MTk5ODQ0MzkxMzliBXdlYmltahN3ZWJpbV8xNjE5OTg0NDM5MTM5coUBJHQkWVdNdDdjSlA2S3RxRWV1M1U2ZnFZUHRUTGFLbEU0RC02Ukhrc1lnSnhTc2I4akxLa0twdzJXb1I2WTJjczVvTkJzTDZBd01BQUFGNUxob1JPZ0JQR2dEMFpLZkd2NXg2bXNRMldXSl9WQmJzNXBEM0tUTzNmdXg1NTd3aWRMelp2UVAAWAA="]`))
				if nil != err {
					log.Println(err)
					break
				}
			}
			if strings.HasPrefix(string(messageData), "a[\"") {
				sEnc := (string(messageData))

				sEnc = strings.TrimLeft(sEnc, "a[\"")

				sEnc = strings.TrimRight(sEnc, "\"]")

				sDec, _ := base64.StdEncoding.DecodeString(sEnc)
				fmt.Println(string(sDec))
				fmt.Println(sDec)

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

package client

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func (u *cxUser) imConn() {
	dialer := websocket.Dialer{}
	var err error
	u.conn, _, err = dialer.Dial(imURL, nil)
	if err != nil {
		log.Println("IM连接失败!")
		return
	}
	log.Println("IM连接成功!")
}

func (u *cxUser) IMStart() {
	u.imConn()
	//离开作用域关闭连接，go 的常规操作
	//defer u.conn.Close()
	for {
		var msg connMsg
		messageType, messageData, err := u.conn.ReadMessage()
		if nil != err {
			log.Println(err)
			log.Println("IM已断开连接!")
			log.Println("准备重试...")
			time.Sleep(5 * time.Second)
			u.imConn()
			//continue
		}
		msg.Type = messageType

		switch messageType {
		case websocket.TextMessage: //文本数据
			fmt.Println(string(messageData))
			if string(messageData) == "o" {
				//connect.WriteJSON()
				log.Println(0)
				err := u.conn.WriteMessage(websocket.TextMessage, []byte(`["CAASPAoOY3gtZGV2I2N4c3R1ZHkSCDc1MDUwMDQ4GgtlYXNlbW9iLmNvbSITd2ViaW1fMTYxOTk4NDQzOTEzORqFASR0JFlXTXQ3Y0pQNkt0cUVldTNVNmZxWVB0VExhS2xFNEQtNlJIa3NZZ0p4U3NiOGpMS2tLcHcyV29SNlkyY3M1b05Cc0w2QXdNQUFBRjVMaG9ST2dCUEdnRDBaS2ZHdjV4Nm1zUTJXV0pfVkJiczVwRDNLVE8zZnV4NTU3d2lkTHpadlFAA0rAAQgQEgUzLjAuMCgAMABKDTE2MTk5ODQ0MzkxMzliBXdlYmltahN3ZWJpbV8xNjE5OTg0NDM5MTM5coUBJHQkWVdNdDdjSlA2S3RxRWV1M1U2ZnFZUHRUTGFLbEU0RC02Ukhrc1lnSnhTc2I4akxLa0twdzJXb1I2WTJjczVvTkJzTDZBd01BQUFGNUxob1JPZ0JQR2dEMFpLZkd2NXg2bXNRMldXSl9WQmJzNXBEM0tUTzNmdXg1NTd3aWRMelp2UVAAWAA="]`))
				if nil != err {
					log.Println(err)
					//continue
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
			log.Println("IM已断开连接!")
			log.Println("准备重试...")
			time.Sleep(5 * time.Second)
			u.imConn()
		case websocket.PingMessage: //Ping
		case websocket.PongMessage: //Pong
		default:

		}
	}

}

func (u *cxUser) IMHandle() {
}

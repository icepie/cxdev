package client

import (
	"bytes"
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

func (u *cxUser) sendPackage(data []byte) {

	log.Println("CXIM: Message send: " + string(data))

	err := u.conn.WriteMessage(websocket.TextMessage, data)
	if nil != err {
		log.Println(err)
	}
}

func (u *cxUser) imLogin() {

	timeUnix := time.Since(time.Unix(0, 0)).Milliseconds()

	buf := bytes.Buffer{}
	//data := []byte{0x08, 0x00, 0x12, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	buf.Write([]byte{0x08, 0x00, 0x12})
	buf.WriteByte(byte(52 + len(fmt.Sprint(u.TUID))))
	buf.Write([]byte{0x0a, 0x0e})
	buf.Write([]byte("cx-dev#cxstudy"))
	buf.Write([]byte{0x12})
	buf.WriteByte(byte(len(fmt.Sprint(u.TUID))))
	buf.Write([]byte(fmt.Sprint(u.TUID)))
	buf.Write([]byte{0x1a, 0x0b})
	buf.Write([]byte("easemob.com"))
	buf.Write([]byte{0x22, 0x13})
	buf.Write([]byte("webim_" + fmt.Sprint(timeUnix)))
	buf.Write([]byte{0x1a, 0x85, 0x01})
	buf.Write([]byte("$t$"))
	buf.Write([]byte(u.IMToken))
	buf.Write([]byte{0x40, 0x03, 0x4a, 0xc0, 0x01, 0x08, 0x10, 0x12, 0x05, 0x33, 0x2e, 0x30, 0x2e, 0x30, 0x28, 0x00, 0x30, 0x00, 0x4a, 0x0d})
	buf.Write([]byte(fmt.Sprint(timeUnix)))
	buf.Write([]byte{0x62, 0x05, 0x77, 0x65, 0x62, 0x69, 0x6d, 0x6a, 0x13, 0x77, 0x65, 0x62, 0x69, 0x6d, 0x5f})
	buf.Write([]byte(fmt.Sprint(timeUnix)))
	buf.Write([]byte{0x72, 0x85, 0x01, 0x24, 0x74, 0x24})
	buf.Write([]byte(u.IMToken))
	buf.Write([]byte{0x50, 0x00, 0x58, 0x00})

	data := `["` + base64.StdEncoding.EncodeToString(buf.Bytes()) + `"]`

	u.sendPackage([]byte(data))
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
			log.Println("CXIM: IM已断开连接!")
			log.Println("CXIM: 准备重试...")
			time.Sleep(5 * time.Second)
			u.imConn()
		}
		msg.Type = messageType

		switch messageType {
		case websocket.TextMessage: //文本数据

			//log.Println("CXIM: Message received:", string(messageData))

			if strings.HasPrefix(string(messageData), "o") {
				log.Println("CXIM: 准备登陆...")
				u.imLogin()
			} else if strings.HasPrefix(string(messageData), "a") {
				sEnc := (string(messageData))
				sEnc = strings.TrimLeft(sEnc, "a[\"")
				sEnc = strings.TrimRight(sEnc, "\"]")

				sDec, _ := base64.StdEncoding.DecodeString(sEnc)

				//log.Println(string(sDec))
				// fmt.Println(sDec)

				if bytes.HasPrefix(sDec, []byte{0x08, 0x00, 0x40, 0x03, 0x4a}) {
					//log.Println("CXIM: 登陆成功!")

					// 解析客户端ID
					var d bytes.Buffer
					for i, v := range sDec {
						if i >= 14 {
							d.Write([]byte{v})
						}
					}

					clientID := strings.Trim(d.String(), `"`)

					log.Println("CXIM: 登陆成功:", clientID)

				} else if bytes.HasPrefix(sDec, []byte{0x08, 0x00, 0x40, 0x02, 0x4a}) {

					len := int(sDec[9])
					if err != nil {
						log.Println(err)
					}

					if len > 0 {
						// 解析ChatID
						var d bytes.Buffer
						for i, v := range sDec {
							if i >= 10 && i <= 10+len {
								d.Write([]byte{v})
							}
						}

						chatID := strings.Trim(d.String(), `"`)

						log.Println("CXIM: 收到消息通知 -> ChatID:", chatID)

						log.Println("CXIM: 正在获取消息详情...")

						// 正在构造获取消息的数据包
						s := sDec
						s[3] = 0x00
						s[6] = 0x1a
						ss := bytes.NewBuffer(s)
						ss.Write([]byte{0x58, 0x00})
						data := `["` + base64.StdEncoding.EncodeToString(ss.Bytes()) + `"]`

						// 发送数据包
						u.sendPackage([]byte(data))

						// 标为已读
						u.MarkRead(chatID)
					}

				} else if bytes.HasPrefix(sDec, []byte{0x08, 0x00, 0x40, 0x00, 0x4a}) {

					if strings.Contains(string(sDec), `*]`) {

						tmp := fmt.Sprintf(`%d`, u.TUID)

						start := bytes.LastIndex(sDec, []byte(tmp)) + len(tmp) + 6

						end := bytes.LastIndex(sDec, []byte{42, 93}) - 1

						buf := bytes.Buffer{}

						log.Println("Start: ", start, "End: ", end)

						for i, b := range sDec {
							if i >= start && i <= end {
								buf.Write([]byte{b})
							}
						}

						log.Println("CXIM: 成功获取消息详情 -> Msg:", buf.String())
					}
				}

			} else if strings.HasPrefix(string(messageData), "h") {
				log.Println("CXIM: 心跳")
			}
		case websocket.BinaryMessage: //二进制数据
			fmt.Println(messageData)
		case websocket.CloseMessage: //关闭
			log.Println("CXIM: IM已断开连接!")
			log.Println("CXIM: 准备重试...")
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

func (u *cxUser) MarkRead(chatID string) {

	buf := bytes.Buffer{}

	buf.Write([]byte{8, 0, 64, 0, 74, 22, 16, 220, 148, 128, 164, 216, 221, 151, 151, 12, 26, 10, 18, 8})
	buf.Write([]byte(chatID))
	buf.Write([]byte{88, 0})

	data := `["` + base64.StdEncoding.EncodeToString(buf.Bytes()) + `"]`

	u.sendPackage([]byte(data))
}

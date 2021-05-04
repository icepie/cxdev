package main

import (
	"log"

	"github.com/icepie/cxdev/client"
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

	courses, err := cxuser.GetCourses()
	if err != nil {
		log.Println(err)
	}

	log.Println(courses)

	//log.Println(fmt.Sprintf("%x", buf.Bytes()))
	//s := "48656c6c6f"

	// a := "73616d706c65"
	// bs, err := hex.DecodeString(a)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(string(bs))

	// src := []byte("Hello Gopher!")

	// dst := make([]byte, hex.EncodedLen(len(src)))
	// hex.Encode(dst, src)

	// log.Println(dst)

	// one := make([]byte, 2)
	// two := make([]byte, 1)
	// //16进制
	// one[0] = 0x08
	// one[1] = 0x00
	// two[0] = 0x12
	// test, err := hex.DecodeString("145")

	// src1 := []byte("t6")

	// dst1 := make([]byte, hex.EncodedLen(len(src1)))
	// hex.Encode(dst1, src1)

	// fmt.Println(test)
	// if err != nil {
	// 	fmt.Println(err)

	// }
	// //two[1] = 0x030
	// res1 := append(one[:], two[:]...)
	// res1 = append(dst1[:], res1[:]...)
	// fmt.Println(res1)

	// log.Println(a)

	// imgRte, err := cxuser.UploadImage("./test.png")
	// if err != nil {
	// 	log.Println(err)
	// }

	// log.Println(imgRte)

	//cxuser.GetActivelist(204942811, 10173395)

	cxuser.IMStart()
}

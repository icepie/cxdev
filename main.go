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

	// imgRte, err := cxuser.UploadImage("./test.png")
	// if err != nil {
	// 	log.Println(err)
	// }

	// log.Println(imgRte)

	//cxuser.GetActivelist(204942811, 10173395)

	cxuser.IMStart()
}

package client

import (
	"fmt"

	"github.com/asmcos/requests"
)

type CXUser struct {
	Username string
	Password string
	SchoolID string
	UID      uint64
}

func NewXUser(Username string, Password string, SchoolID string) error {

	req := requests.Requests()

	p := requests.Params{
		"name":     Username,
		"pwd":      Password,
		"schoolid": SchoolID,
		"verify":   "0",
	}

	resp, _ := req.Get("https://passport2.chaoxing.com/api/login", p)

	println(resp.Text())

	coo := resp.Cookies()
	// coo is [] *http.Cookies
	println("********cookies*******")
	for _, c := range coo {
		fmt.Println(c.Name, c.Value)
	}
	return nil
}

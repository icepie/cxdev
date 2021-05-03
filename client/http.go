package client

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/asmcos/requests"
)

type CXUser struct {
	Username  string
	Password  string
	SchoolID  uint64
	UID       uint64
	RealName  string
	Uname     string
	RoleID    string
	Phone     string
	CxID      uint64
	DxfID     uint64
	Email     string
	IsCertify uint64
	OpacPwd   bool
	Token     string
	TUID      uint64
	IMToken   string
	Cooikes   []*http.Cookie
}

func NewCXUser(Username string, Password string, SchoolID string) (CXUser, error) {

	// 初始化
	var u CXUser
	u.Username = Username
	u.Password = Password

	// 尝试登陆
	req := requests.Requests()

	p := requests.Params{
		"name":     Username,
		"pwd":      Password,
		"schoolid": SchoolID,
		"verify":   "0",
	}

	resp, err := req.Get("https://passport2.chaoxing.com/api/login", p)
	if err != nil {
		return u, err
	}

	log.Println(resp.Text())

	var rte LoginRte
	err = json.Unmarshal([]byte(resp.Text()), &rte)
	if err != nil {
		return u, err
	} else {
		u.Cooikes = resp.Cookies()
		u.UID = rte.UID
		u.SchoolID = rte.SchoolID
		u.RealName = rte.RealName
		u.Uname = rte.Uname
		u.RoleID = rte.RoleID
		u.CxID = rte.CxID
		u.DxfID = rte.DxfID
		u.Email = rte.Email
		u.OpacPwd = rte.OpacPwd
		u.getToken()
		// 初始化im
		u.initIM()
	}

	if !rte.Result {
		return u, errors.New(rte.ErrorMsg)
	}

	return u, nil
}

// 获取 Token
func (u *CXUser) getToken() {

	var token string

	// 获取token
	req := requests.Requests()

	for _, cooike := range u.Cooikes {
		req.SetCookie(cooike)
	}

	resp, err := req.Get("https://pan-yz.chaoxing.com/api/token/uservalid")
	if err != nil {
		return
	}

	var rte TokenRte
	err = json.Unmarshal([]byte(resp.Text()), &rte)
	if err != nil {
		return
	}

	if rte.Result {
		token = rte.Token
		log.Println("test")
	} else {
		return
	}

	u.Token = token
}

// 获取 Token
func (u *CXUser) initIM() {

	// 获取token
	req := requests.Requests()

	for _, cooike := range u.Cooikes {
		req.SetCookie(cooike)
	}

	resp, err := req.Get("https://im.chaoxing.com/webim/me")
	if err != nil {
		return
	}

	reg := regexp.MustCompile(`loginByToken\('(\d+?)', '([^']+?)'\);`)
	if reg == nil {
		log.Println("MustCompile err")
		return
	}

	//提取关键信息
	result := reg.FindAllStringSubmatch(resp.Text(), -1)

	tuid, err := strconv.Atoi(result[0][1])

	if err != nil {
		return
	}

	u.TUID = uint64(tuid)
	u.IMToken = result[0][2]
}

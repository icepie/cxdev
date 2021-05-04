package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/asmcos/requests"
	"github.com/gorilla/websocket"
)

type connMsg struct {
	Type int
	Data []byte
}

type cxUser struct {
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
	Running   bool
	Cooikes   []*http.Cookie
	conn      *websocket.Conn
	MsgChan   chan connMsg
}

func NewCXUser(Username string, Password string, SchoolID string) (cxUser, error) {

	// 初始化
	var u cxUser
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

	resp, err := req.Get(loginURL, p)
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
func (u *cxUser) getToken() {

	var token string

	// 获取token
	req := requests.Requests()

	for _, cooike := range u.Cooikes {
		req.SetCookie(cooike)
	}

	resp, err := req.Get(getPanTokenURL)
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
	} else {
		return
	}

	u.Token = token
}

// 获取 IM Token
func (u *cxUser) initIM() {

	// 获取token
	req := requests.Requests()

	for _, cooike := range u.Cooikes {
		req.SetCookie(cooike)
	}

	resp, err := req.Get(getIMTokenURL)
	if err != nil {
		return
	}

	reg := regexp.MustCompile(`loginByToken\('(\d+?)', '([^']+?)'\);`)

	//提取关键信息
	result := reg.FindAllStringSubmatch(resp.Text(), -1)

	if len(result) == 0 {
		return
	}

	tuid, err := strconv.Atoi(result[0][1])

	if err != nil {
		return
	}

	u.TUID = uint64(tuid)
	u.IMToken = result[0][2]
}

func (u cxUser) GetCourses() (courses []Course, err error) {

	req := requests.Requests()

	for _, cooike := range u.Cooikes {
		req.SetCookie(cooike)
	}

	p := requests.Params{
		"rss":        "1",
		"start":      "0",
		"size":       "500",
		"catalogId":  "0",
		"searchname": "",
	}

	resp, err := req.Get(getCoursesURL, p)
	if err != nil {
		return
	}

	// resp, err := req.Get("http://mooc1-2.chaoxing.com/visit/interaction")
	// if err != nil {
	// 	return
	// }

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.Text()))
	if err != nil {
		return
	}

	//log.Println(resp.Text())

	// 找到课程
	ul := doc.Find("ul#courseList.course-list").First()

	// 遍历课程
	ul.Find("li.course").Each(func(index int, li *goquery.Selection) {
		courseCover := li.Find("div.course-cover").First()
		courseInfo := li.Find("div.course-info").First()

		var course Course

		clazzId, e := courseCover.Find("input.clazzId").First().Attr("value")
		if e {
			classID, _ := strconv.Atoi(clazzId)
			course.ClassID = uint64(classID)
			//log.Println("clazzId: ", clazzId)
		}

		courseId, e := courseCover.Find("input.courseId").First().Attr("value")
		if e {
			courseID, _ := strconv.Atoi(courseId)
			course.CourseID = uint64(courseID)
			//log.Println("courseId: ", courseId)
		}

		role, e := courseCover.Find("input.role").First().Attr("value")
		if e {
			course.Role, _ = strconv.Atoi(role)
			//log.Println("role: ", role)
		}

		imgurl, e := courseCover.Find("a>img").First().Attr("src")
		if e {
			course.ImgURL = imgurl
		}

		name := courseInfo.Find("a.color1>span.course-name").First()
		course.Name = name.Text()

		info := courseInfo.Find("p.color2").First()
		course.Info = info.Text()
		//log.Println("info: ", info.Text())

		teacher := courseInfo.Find("p.color3").First()
		course.Teacher = teacher.Text()
		//log.Println("teacher: ", teacher.Text())

		// 设置班级信息
		classRte, err := u.getClassDetail(courseId, clazzId)
		if err == nil {
			course.BBSID = classRte.Data.BBSID

			chatID, _ := strconv.Atoi(classRte.Data.ChatID)
			course.Class.ChatID = uint64(chatID)

			course.Class.Name = classRte.Data.Name
			course.Class.StudentCount = classRte.Data.StudentCount

			CreatorUserID, _ := strconv.Atoi(classRte.Data.CreatorUserID)
			course.Class.CreatorUserID = uint64(CreatorUserID)
		}

		courses = append(courses, course)
		//  else {
		// 	log.Println(err)
		// }
		// log.Println(courseCover.Html())
		// log.Println(courseInfo.Html())
	})

	return
	//doc.

	// reg := regexp.MustCompile(`\?courseid=(\d+?)&clazzid=(\d+?)&cpi=\d+`)

	// //提取关键信息
	// result := reg.FindAllStringSubmatch(resp.Text(), -1)

	// if len(result) == 0 {
	// 	return
	// }

	// log.Println(result)

}

func (u cxUser) getClassDetail(courseId string, classId string) (rte ClassRte, err error) {

	req := requests.Requests()

	for _, cooike := range u.Cooikes {
		req.SetCookie(cooike)
	}

	p := requests.Params{
		"fid":      fmt.Sprint(u.SchoolID),
		"courseId": courseId,
		"classId":  classId,
	}

	resp, err := req.Get(getClassDetailURL, p)
	if err != nil {
		return
	}

	//log.Println(resp.Text())

	err = json.Unmarshal([]byte(resp.Text()), &rte)
	if err != nil {
		return
	}

	if rte.Result == 0 {
		return rte, errors.New(fmt.Sprint(rte.ErrorMsg))
	}

	return
}

func (u cxUser) UploadImage(path string) (rte UploadImageRte, err error) {

	if !PathExist(path) {
		return rte, errors.New("the file is not available")
	}

	req := requests.Requests()

	for _, cooike := range u.Cooikes {
		req.SetCookie(cooike)
	}

	p := requests.Params{
		"puid":   fmt.Sprint(u.UID),
		"_token": u.Token,
	}

	f := requests.Files{
		"file": path,
	}

	resp, err := req.Post(uploadPanURL, p, f)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp.Text()), &rte)
	if err != nil {
		return
	}

	if !rte.Result {
		return rte, errors.New(rte.Msg)
	}

	return
}

func (u cxUser) GetActivelist(courseId uint64, classId uint64) {
	req := requests.Requests()

	for _, cooike := range u.Cooikes {
		req.SetCookie(cooike)
	}

	p := requests.Params{
		"fid":      fmt.Sprint(u.SchoolID),
		"courseId": fmt.Sprint(courseId),
		"classId":  fmt.Sprint(classId),
	}

	resp, err := req.Get(getActiveListURL, p)
	if err != nil {
		return
	}

	log.Println(resp.Text())
}

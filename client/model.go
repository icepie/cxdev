package client

// LoginRte 登陆响应
type LoginRte struct {
	CxID      uint64 `json:"cxid"`
	DxfID     uint64 `json:"dxfid"`
	Email     string `json:"email"`
	IsCertify uint64 `json:"isCertify"`
	OpacPwd   bool   `json:"opacPwd"`
	Phone     string `json:"phone"`
	RealName  string `json:"realname"`
	Result    bool   `json:"result"`
	RoleID    string `json:"roleid"`
	SchoolID  uint64 `json:"schoolid"`
	Status    string `json:"status"`
	UID       uint64 `json:"uid"`
	Uname     string `json:"uname"`
	ErrorMsg  string `json:"errorMsg"`
	URL       string `json:"url"`
}

// TokenRte 获取Token响应
type TokenRte struct {
	Token  string `json:"_token"`
	Result bool   `json:"result"`
}

type ClassRte struct {
	Data struct {
		BBSID  string `json:"bbsid"`
		ChatID string `json:"chatid"`
		Course struct {
			Data []struct {
				ID       uint64 `json:"id"`
				ImageURL string `json:"imageurl"`
				Name     string `json:"name"`
			} `json:"data"`
		} `json:"course"`
		CourseID      uint64 `json:"courseid"`
		CreatorUserID string `json:"creatoruserid"`
		Name          string `json:"name"`
		StudentCount  uint64 `json:"studentcount"`
	} `json:"data"`
	ErrorMsg interface{} `json:"errorMsg"`
	Msg      interface{} `json:"msg"`
	Result   int         `json:"result"`
}

// Class 班级模型

type Class struct {
	Name          string
	StudentCount  uint64
	ChatID        uint64
	CreatorUserID uint64
}

// Course 课程模型
type Course struct {
	Name     string
	ClassID  uint64
	CourseID uint64
	Role     int
	Info     string
	Teacher  string
	BBSID    string
	ImgURL   string
	Class    Class
}

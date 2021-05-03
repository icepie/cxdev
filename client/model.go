package client

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

type TokenRte struct {
	Token  string `json:"_token"`
	Result bool   `json:"result"`
}

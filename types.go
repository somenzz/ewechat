package ewechat

import (
	"sync"
	"time"
)

type EWechat struct {
	CorpID      string
	CorpSecret  string
	AgentID     int
	token       string
	tokenExpiry time.Time
	mutex       sync.Mutex
}

type WechatAccessTokenResponse struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type WechatPostMessageResponse struct {
	ErrCode        int    `json:"errcode"`
	ErrMsg         string `json:"errmsg"`
	InvalidUser    string `json:"invaliduser"`
	InvalidParty   string `json:"invalidparty"`
	InvalidTag     string `json:"invalidtag"`
	UnlicensedUser string `json:"unlicenseduser"`
	MsgID          string `json:"msgid"`
	ResponseCode   string `json:"response_code"`
}

type TextMessage struct {
	ToUser                 string      `json:"touser"`
	ToParty                string      `json:"toparty"`
	ToTag                  string      `json:"totag"`
	MsgType                string      `json:"msgtype"`
	AgentID                int         `json:"agentid"`
	Text                   TextContent `json:"text"`
	Safe                   int         `json:"safe"`
	EnableIDTrans          int         `json:"enable_id_trans"`
	EnableDuplicateCheck   int         `json:"enable_duplicate_check"`
	DuplicateCheckInterval int         `json:"duplicate_check_interval"`
}

type TextContent struct {
	Content string `json:"content"`
}

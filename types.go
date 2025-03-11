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

// UploadMediaResponse 定义接口返回的 JSON 结构体
type UploadMediaResponse struct {
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}

// MessageType 定义消息类型常量
type MessageType string

const (
	MsgTypeText     MessageType = "text"
	MsgTypeImage    MessageType = "image"
	MsgTypeVoice    MessageType = "voice"
	MsgTypeVideo    MessageType = "video"
	MsgTypeFile     MessageType = "file"
	MsgTypeTextCard MessageType = "textcard"
	MsgTypeNews     MessageType = "news"
	MsgTypeMPNews   MessageType = "mpnews"
	MsgTypeMarkdown MessageType = "markdown"
)

// BaseMessage 基础消息结构
type BaseMessage struct {
	ToUser                 string      `json:"touser,omitempty"`
	ToParty                string      `json:"toparty,omitempty"`
	ToTag                  string      `json:"totag,omitempty"`
	MsgType                MessageType `json:"msgtype"`
	AgentID                int         `json:"agentid,omitempty"`
	Safe                   int         `json:"safe,omitempty"`
	EnableIDTrans          int         `json:"enable_id_trans,omitempty"`
	EnableDuplicateCheck   int         `json:"enable_duplicate_check,omitempty"`
	DuplicateCheckInterval int         `json:"duplicate_check_interval,omitempty"`
}

// MessageContent 消息内容联合类型
type MessageContent struct {
	Text     *TextContent     `json:"text,omitempty"`
	Image    *ImageContent    `json:"image,omitempty"`
	Voice    *VoiceContent    `json:"voice,omitempty"`
	Video    *VideoContent    `json:"video,omitempty"`
	File     *FileContent     `json:"file,omitempty"`
	TextCard *TextCardContent `json:"textcard,omitempty"`
	News     *NewsContent     `json:"news,omitempty"`
	MPNews   *MPNewsContent   `json:"mpnews,omitempty"`
	Markdown *MarkdownContent `json:"markdown,omitempty"`
}

type ImageContent struct {
	MediaID string `json:"media_id"`
}

type VoiceContent struct {
	MediaID string `json:"media_id"`
}

type VideoContent struct {
	MediaID     string `json:"media_id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type FileContent struct {
	MediaID string `json:"media_id"`
}

type TextCardContent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	BtnTxt      string `json:"btntxt,omitempty"`
}

type NewsArticle struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
	PicURL      string `json:"picurl,omitempty"`
	AppID       string `json:"appid,omitempty"`
	PagePath    string `json:"pagepath,omitempty"`
}

type NewsContent struct {
	Articles []NewsArticle `json:"articles"`
}

type MPNewsArticle struct {
	Title            string `json:"title"`
	ThumbMediaID     string `json:"thumb_media_id"`
	Author           string `json:"author,omitempty"`
	ContentSourceURL string `json:"content_source_url,omitempty"`
	Content          string `json:"content"`
	Digest           string `json:"digest,omitempty"`
}

type MPNewsContent struct {
	Articles []MPNewsArticle `json:"articles"`
}

type MarkdownContent struct {
	Content string `json:"content"`
}

// Response 响应结构
type Response struct {
	ErrCode        int    `json:"errcode"`
	ErrMsg         string `json:"errmsg"`
	InvalidUser    string `json:"invaliduser,omitempty"`
	InvalidParty   string `json:"invalidparty,omitempty"`
	InvalidTag     string `json:"invalidtag,omitempty"`
	UnlicensedUser string `json:"unlicenseduser,omitempty"`
	MsgID          string `json:"msgid,omitempty"`
	ResponseCode   string `json:"response_code,omitempty"`
}

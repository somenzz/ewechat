package ewechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (e *EWechat) getToken() (string, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// Check if the current token is still valid
	if e.token != "" && time.Now().Before(e.tokenExpiry) {
		return e.token, nil
	}
	// If not, fetch a new token
	apiURL := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + e.CorpID + "&corpsecret=" + e.CorpSecret
	// ... rest of the HTTP request code ...
	request, _ := http.NewRequest("GET", apiURL, nil)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	var model WechatAccessTokenResponse
	err = json.Unmarshal(body, &model)
	if err != nil {
		return "", err
	}

	// Update the token and its expiry time
	if model.ErrCode == 0 { //出错返回码，为0表示成功，非0表示调用失败
		e.token = model.AccessToken
		e.tokenExpiry = time.Now().Add(time.Duration(model.ExpiresIn) * time.Second)
		return e.token, nil
	}
	return "", fmt.Errorf(model.ErrMsg)
}

func (e *EWechat) SendMessage(text string, users string) (string, error) {

	model := TextMessage{}
	model.ToUser = users
	model.MsgType = "text"
	model.AgentID = e.AgentID
	model.Text.Content = text
	model.Safe = 0
	requestJSON, _ := json.Marshal(model)
	token, err := e.getToken()
	if err != nil {
		return "", err
	}
	apiURL := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + token
	request, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestJSON))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	var responseModel WechatPostMessageResponse
	json.Unmarshal(body, &responseModel)
	return responseModel.ErrMsg, nil
}
func (e *EWechat) GetUserID(telephone string) (string, error) {
	token, err := e.getToken()
	if err != nil {
		return "", err
	}

	apiURL := "https://qyapi.weixin.qq.com/cgi-bin/user/getuserid?access_token=" + token
	body := map[string]string{"mobile": telephone}
	requestJSON, _ := json.Marshal(body)
	request, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestJSON))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody, _ := io.ReadAll(response.Body)
	var res map[string]interface{}
	json.Unmarshal(responseBody, &res)

	if res["errmsg"] == "ok" {
		return res["userid"].(string), nil
	} else {
		return "", fmt.Errorf(res["errmsg"].(string))
	}
}

// UploadMedia 上传临时素材到企业微信
func (e *EWechat) UploadMedia(mediaType, filePath string) (*UploadMediaResponse, error) {
	// 验证 mediaType 是否合法
	validTypes := map[string]bool{"image": true, "voice": true, "video": true, "file": true}
	if !validTypes[mediaType] {
		return nil, fmt.Errorf("invalid media type: %s, must be image, voice, video, or file", mediaType)
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// 检查文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}
	fileSize := fileInfo.Size()
	if fileSize < 5 {
		return nil, fmt.Errorf("file size must be greater than 5 bytes, got %d", fileSize)
	}

	// 检查文件大小限制
	switch mediaType {
	case "image":
		if fileSize > 10*1024*1024 {
			return nil, fmt.Errorf("image size exceeds 10MB limit")
		}
	case "voice":
		if fileSize > 2*1024*1024 {
			return nil, fmt.Errorf("voice size exceeds 2MB limit")
		}
	case "video":
		if fileSize > 10*1024*1024 {
			return nil, fmt.Errorf("video size exceeds 10MB limit")
		}
	case "file":
		if fileSize > 20*1024*1024 {
			return nil, fmt.Errorf("file size exceeds 20MB limit")
		}
	}

	// 检查文件格式（根据扩展名简单验证）
	ext := strings.ToLower(filepath.Ext(filePath))
	switch mediaType {
	case "image":
		if ext != ".jpg" && ext != ".png" {
			return nil, fmt.Errorf("image must be JPG or PNG format")
		}
	case "voice":
		if ext != ".amr" {
			return nil, fmt.Errorf("voice must be AMR format")
		}
	case "video":
		if ext != ".mp4" {
			return nil, fmt.Errorf("video must be MP4 format")
		}
	}

	// 构造 multipart/form-data 请求体
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("media", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content: %v", err)
	}
	writer.Close()
	accessToken, err := e.getToken()
	if err != nil {
		return nil, err
	}
	// 构造请求 URL
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/media/upload?access_token=%s&type=%s", accessToken, mediaType)

	// 发送 POST 请求
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 解析返回的 JSON 数据
	var result UploadMediaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// 检查接口返回的错误码
	if result.Errcode != 0 {
		return nil, fmt.Errorf("API error: errcode=%d, errmsg=%s", result.Errcode, result.Errmsg)
	}

	return &result, nil
}

// SendMessage 发送消息的主函数
func (e *EWechat) SendAllTypesMessage(msg BaseMessage, content MessageContent) (*Response, error) {
	// 验证基本参数
	accessToken, err := e.getToken()
	if err != nil {
		return nil, err
	}
	if accessToken == "" {
		return nil, fmt.Errorf("access_token cannot be empty")
	}
	if msg.AgentID == 0 {
		msg.AgentID = e.AgentID
	}
	if msg.MsgType == "" {
		return nil, fmt.Errorf("msg_type cannot be empty")
	}

	// 创建完整消息结构
	message := struct {
		BaseMessage
		MessageContent
	}{
		BaseMessage:    msg,
		MessageContent: content,
	}

	// JSON序列化
	body, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %v", err)
	}

	// 构造请求URL
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", accessToken)

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 发送请求
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// 检查API返回的错误
	if response.ErrCode != 0 {
		return &response, fmt.Errorf("api error: %d - %s", response.ErrCode, response.ErrMsg)
	}

	return &response, nil
}

// // 示例使用
// func Example() {
// 	// 发送文本消息
// 	textMsg := BaseMessage{
// 		ToUser:                 "UserID1|UserID2",
// 		MsgType:                MsgTypeText,
// 		AgentID:                1,
// 		EnableDuplicateCheck:   0,
// 		DuplicateCheckInterval: 1800,
// 	}

// 	content := MessageContent{
// 		Text: &TextContent{
// 			Content: "你的快递已到，请携带工卡前往邮件中心领取。",
// 		},
// 	}

// 	resp, err := SendMessage("YOUR_ACCESS_TOKEN", textMsg, content)
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("Response: %+v\n", resp)
// }

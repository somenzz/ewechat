package ewechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

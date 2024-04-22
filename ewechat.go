package ewechat

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type EWechat struct {
	CorpID     string
	CorpSecret string
	AgentID    int
}

func (e *EWechat) getToken() (string, error) {
	type WechatAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	apiURL := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + e.CorpID + "&corpsecret=" + e.CorpSecret
	request, _ := http.NewRequest("GET", apiURL, nil)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	var model WechatAccessTokenResponse

	json.Unmarshal(body, &model)

	return model.AccessToken, nil
}

func (e *EWechat) SendMessage(text string, users string) (string, error) {
	type SendRequest struct {
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		AgentID int    `json:"agentid"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
		Safe int `json:"safe"`
	}

	type WechatPostMessageReponse struct {
		ErrorMessage string `json:"errmsg"`
	}

	model := SendRequest{}
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
	var responseModel WechatPostMessageReponse

	json.Unmarshal(body, &responseModel)

	return responseModel.ErrorMessage, nil
}

package ai

import (
	"ai-text-helper-server/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type AiRequestStruct struct {
	Url     string `json:"url"`
	Payload string `json:"payload"`
}

type CustomizeRequestStruct struct {
	System         string `json:"system"`
	RequestContent string `json:"requestContent"`
}
var urlMap = map[string]string{
	"translate":  "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/ernie_speed?access_token=",
	"completion": "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/ernie_speed?access_token=",
	"customize":  "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions_pro?access_token=",
}
var payloadMap = map[string]string{
	"translate":  `{"messages":[{"role":"user","content":"请将「%s」进行翻译"}],"temperature":1,"system":"你是一个中英翻译助手，我将给你一段英文或者中文。如果我给你英文，请将之翻译为中文，反之则翻译为英文。"}`,
	"completion": `{"messages":[{"role":"user","content":"请为「%s」进行续写"}],"temperature":1,"system":"你是一个文本续写助手，我将给你一段文本。你需要帮我进行续写，要求符合上下文语意"}`,
	"customize":  `{"messages":[{"role":"user","content":"%s"}],"system":"%s"}`,
}

/**
 * 使用 AK，SK 生成鉴权签名（Access Token）
 * @return string 鉴权签名信息（Access Token）
 */
func InitAccessToken() {
	url := "https://aip.baidubce.com/oauth/2.0/token"
	apiKey, err := utils.GetEnvApiKey()
	if err != nil {
		return
	}
	secretKey, err := utils.GetEnvSecretKey()
	if err != nil {
		return
	}
	postData := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", apiKey, secretKey)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(postData))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	accessTokenObj := map[string]string{}
	json.Unmarshal([]byte(body), &accessTokenObj)
	err = utils.SetEnv("ACCESS_TOKEN", accessTokenObj["access_token"])
	fmt.Println(accessTokenObj["access_token"])
	if err != nil {
		return
	}
}
func GenerateParameter(requestType, content string) (AiRequestStruct, error) {
	accessToken, err := utils.GetEnvAccessToken()
	if err != nil {
		return AiRequestStruct{}, err
	}
	arcUrl, ok := urlMap[requestType]
	if !ok {
		return AiRequestStruct{}, errors.New("请求类型错误！")
	}
	arcPayload, ok := payloadMap[requestType]
	if !ok {
		return AiRequestStruct{}, errors.New("请求类型错误！")
	}
	var ars = AiRequestStruct{
		Url:     arcUrl + accessToken,
		Payload: fmt.Sprintf(arcPayload, content),
	}
	return ars, nil
}
func GenerateCustomizeParameter(requestStruct CustomizeRequestStruct) (AiRequestStruct, error) {
	accessToken, err := utils.GetEnvAccessToken()
	if err != nil {
		return AiRequestStruct{}, err
	}
	var ars = AiRequestStruct{
		Url:     urlMap["customize"] + accessToken,
		Payload: fmt.Sprintf(payloadMap["customize"], requestStruct.RequestContent, requestStruct.System),
	}
	return ars, nil
}
func fetch(ars AiRequestStruct) (map[string]interface{}, error) {
	url := ars.Url
	payload := strings.NewReader(ars.Payload)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	var result map[string]interface{}
	body, err := io.ReadAll(res.Body)
	if err == nil {
		json.Unmarshal([]byte(string(body)), &result)
	}
	return result, nil
}
func AiRequest(requestType, content string) (map[string]interface{}, error) {
	ars, err := GenerateParameter(requestType, content)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := fetch(ars)
	return res, err
}
func CustomizeAiRequest(requestStruct CustomizeRequestStruct) (map[string]interface{}, error) {
	ars, err := GenerateCustomizeParameter(requestStruct)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := fetch(ars)
	return res, err
}

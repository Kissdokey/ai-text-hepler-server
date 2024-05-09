package ai

import (
	"ai-text-helper-server/mysql"
	"ai-text-helper-server/utils"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	_"time"

	"github.com/gin-gonic/gin"
)

type AiRequestStruct struct {
	Url     string `json:"url"`
	Payload string `json:"payload"`
}

type CustomizeRequestStruct struct {
	System         string `json:"system"`
	RequestContent string `json:"requestContent"`
}
type GeneralRequestStruct struct {
	System         string `json:"system"`
	Messages []ChatCyle `json:"messages"`
	Temperature   float64 `json:"temperature"`
	Top_p  float64 `json:"top_p"`
	Penalty_score  float64 `json:"penalty_score"`
	Max_output_tokens 	int `json:"max_output_tokens"`
}

// 目前是小文件，后续增加大文件切片向量化（redis？考虑一下）
type ChatRequestStruct struct {
	FileId      string     `json:"fileId"`
	ChatHistory []ChatCyle `json:"chatHistory"`
}

var urlMap = map[string]string{
	"translate":  "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/ernie_speed?access_token=",
	"completion": "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/ernie_speed?access_token=",
	"customize":  "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions_pro?access_token=",
	"chat":       "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions?access_token=",
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

type ChatCyle struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type ChatPayload struct {
	Messages    []ChatCyle `json:"messages"`
	System      string     `json:"system"`
	Temperature float64    `json:"temperature"`
	Stream      bool       `json:"stream"`
}

func GenerateChatParameter(requestStruct ChatRequestStruct, text string) (AiRequestStruct, error) {
	accessToken, err := utils.GetEnvAccessToken()
	if err != nil {
		return AiRequestStruct{}, err
	}
	payloadBytes, err := json.Marshal(ChatPayload{
		Messages:    requestStruct.ChatHistory,
		System:      "你是智能文本处理助手，我将给你一段文本，和聊天上下文，你需要根据文本和聊天上下文，针对我的问题给出我准确合理的回答。上下文是：" + text,
		Temperature: 1,
		Stream:      true,
	})
	if err != nil {
		return AiRequestStruct{}, err
	}
	var ars = AiRequestStruct{
		Url:     urlMap["chat"] + accessToken,
		Payload: string(payloadBytes),
	}
	return ars, nil
}

func GenerateGeneralParameter(requestStruct GeneralRequestStruct) (AiRequestStruct, error) {
	accessToken, err := utils.GetEnvAccessToken()
	if err != nil {
		return AiRequestStruct{}, err
	}
	payloadBytes, err := json.Marshal(requestStruct)
	if err != nil {
		return AiRequestStruct{}, err
	}
	var ars = AiRequestStruct{
		Url:     urlMap["chat"] + accessToken,
		Payload: string(payloadBytes),
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

func StreamFetch(ars AiRequestStruct, c *gin.Context) error {
	url := ars.Url
	payload := strings.NewReader(ars.Payload)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	c.Stream(func(w io.Writer) bool {
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "event:") {
				event := strings.TrimSpace(strings.Split(line, ":")[1])
				fmt.Fprintf(w, "Event: %s\n", event)
				w.(http.Flusher).Flush()
				if event == "end" {
					break
				}
			} else if strings.Contains(line, "error_code") {
				fmt.Fprintf(w, "%s\n", line)
				w.(http.Flusher).Flush()
				break
			} else if strings.HasPrefix(line, "data:") {
				fmt.Fprintf(w, "{%s}\n", line)
				//防止过快导致前端收到粘黏数据
				// time.Sleep(10 * time.Microsecond)
				w.(http.Flusher).Flush()
			}
		}
		// 返回 false 表示流结束
		return false
	})

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
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
func CustomizeAiRequest(requestStruct GeneralRequestStruct) (map[string]interface{}, error) {
	ars, err := GenerateGeneralParameter(requestStruct)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Print(ars)
	res, err := fetch(ars)
	return res, err
}
func Chat(requestStruct ChatRequestStruct, c *gin.Context) {

	//请求次数减一，返回apiRequestNumber
	username, _ := c.Get("username")
	_, err := mysql.UpdateRequestNumber(username.(string), -1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 501,
			"msg":  "内部错误，请求次数更新失败",
			"data": gin.H{},
		})
		return
	}
	fileId := requestStruct.FileId
	//在数据库中查找filId对应的文件
	file, err := mysql.GetFileInfo(fileId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 501,
			"msg":  "内部错误，文件读取失败",
			"data": gin.H{},
		})
		return
	}
	para, err := GenerateChatParameter(requestStruct, utils.ExtractText(file.Content))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 501,
			"msg":  "内部错误，文件读取失败",
			"data": gin.H{},
		})
		return
	}
	c.Header("Content-Type", "text/event-stream")
	_ = StreamFetch(para, c)
}

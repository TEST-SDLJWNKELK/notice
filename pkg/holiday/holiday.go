package holiday

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var key = ""

// Response 定义结构体以匹配JSON结构
type Response struct {
	Reason    string `json:"reason"`
	Result    Result `json:"result"`
	ErrorCode int    `json:"error_code"`
}

type Result struct {
	Date       string `json:"date"`
	Week       string `json:"week"`
	StatusDesc string `json:"statusDesc"`
	Status     string `json:"status"`
}

func IsHoliday(key string, t time.Time) bool {

	url := "http://apis.juhe.cn/fapig/calendar/day"

	// 接口请求参数
	params := map[string]string{
		"key":  key,
		"date": t.Format("2006-01-02"),
	}

	// 请求头设置
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	rep, err := HttpRequest("GET", url, params, headers, 15)
	if err != nil {
		log.Println("request error：", err.Error())
	} else {
		log.Println("request result：", rep)
	}
	var response Response
	if err = json.Unmarshal([]byte(rep), &response); err != nil {
		log.Println("parse json error:", err)
		return true
	}
	fmt.Println(response)
	if response.Result.StatusDesc == "工作日" {
		return false
	}

	return true
}

// HttpRequest http请求发送
func HttpRequest(method, rawUrl string, bodyMaps, headers map[string]string, timeout time.Duration) (result string, err error) {
	var (
		request  *http.Request
		response *http.Response
		res      []byte
	)
	if timeout <= 0 {
		timeout = 5
	}
	client := &http.Client{
		Timeout: timeout * time.Second,
	}

	// 请求的 body 内容
	data := url.Values{}
	for key, value := range bodyMaps {
		data.Set(key, value)
	}

	jsons := data.Encode()

	if request, err = http.NewRequest(method, rawUrl, strings.NewReader(jsons)); err != nil {
		return
	}

	if method == "GET" {
		request.URL.RawQuery = jsons
	}

	// 增加header头信息
	for key, val := range headers {
		request.Header.Set(key, val)
	}

	// 处理返回结果
	if response, err = client.Do(request); err != nil {
		return "", err
	}

	defer response.Body.Close()

	if res, err = io.ReadAll(response.Body); err != nil {
		return "", err
	}
	return string(res), nil
}

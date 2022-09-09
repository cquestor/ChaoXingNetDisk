package apis

import (
	"ChaoXingNetDisk/datas"
	"ChaoXingNetDisk/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// 登录获取Cookie
func Login(uname, password string) (map[string]string, error) {
	password, err := utils.DesEncrypt([]byte(password), []byte(datas.KEY_DES))
	if err != nil {
		return nil, err
	}
	payload := "fid=-1&uname=" + uname + "&password=" + password + "&refer=http%253A%252F%252Fi.chaoxing.com&t=true&forbidotherlogin=0&validate=&doubleFactorLogin=0&independentId=0"
	req, _ := http.NewRequest(http.MethodPost, datas.URL_LOGIN, strings.NewReader(payload))
	addRequestHeader(req, nil, "application/x-www-form-urlencoded; charset=UTF-8")
	client := &http.Client{Timeout: time.Second * 10, CheckRedirect: disallowRedirect}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var respJson struct {
		Url    string `json:"url"`
		Msg    string `json:"msg2"`
		Status bool   `json:"status"`
	}
	json.Unmarshal(body, &respJson)
	if !respJson.Status {
		return nil, errors.New(respJson.Msg)
	}
	var cookie = make(map[string]string)
	for _, v := range resp.Cookies() {
		cookie[v.Name] = v.Value
	}
	return cookie, nil
}

// 新建文件夹
func NewFolder(cookie map[string]string, folderName string) {
	payload := "parentId=412185986623893504&name=" + folderName + "&selectDlid=allperson&newfileid=0"
	req, _ := http.NewRequest(http.MethodPost, datas.URL_NEW_FOLDER, strings.NewReader(payload))
	addRequestHeader(req, cookie, "application/x-www-form-urlencoded; charset=UTF-8")
	fmt.Println(req.Header)
	client := &http.Client{Timeout: time.Second * 10, CheckRedirect: disallowRedirect}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

// 添加请求头
func addRequestHeader(req *http.Request, cookie map[string]string, contentType string) {
	if cookie != nil {
		var cookies = make([]string, 0, len(cookie))
		for k, v := range cookie {
			temp := k + "=" + v
			cookies = append(cookies, temp)
		}
		req.Header.Add("Cookie", strings.Join(cookies, "; "))
	}
	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36 Edg/105.0.1343.33")
}

// 禁止请求重定向
func disallowRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

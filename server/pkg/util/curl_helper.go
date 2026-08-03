package util

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type OceanEngineBaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func DoGetRequest(accessToken, baseURL string, params url.Values, result interface{}) error {
	fullURL := baseURL + "?" + params.Encode()

	client := &http.Client{Timeout: 30 * time.Second}
	httpReq, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return err
	}

	httpReq.Header.Set("Access-Token", accessToken)

	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, result)
	return err
}

func DoGetRequestWithJsonParams(accessToken, baseURL string, params map[string]interface{}, result interface{}) error {
	queryParts := make([]string, 0)
	for key, value := range params {
		switch v := value.(type) {
		case string:
			queryParts = append(queryParts, url.QueryEscape(key)+"="+url.QueryEscape(v))
		case []string:
			jsonBytes, _ := json.Marshal(v)
			queryParts = append(queryParts, url.QueryEscape(key)+"="+url.QueryEscape(string(jsonBytes)))
		case []interface{}:
			jsonBytes, _ := json.Marshal(v)
			queryParts = append(queryParts, url.QueryEscape(key)+"="+url.QueryEscape(string(jsonBytes)))
		case map[string]interface{}:
			jsonBytes, _ := json.Marshal(v)
			queryParts = append(queryParts, url.QueryEscape(key)+"="+url.QueryEscape(string(jsonBytes)))
		case map[string]string:
			jsonBytes, _ := json.Marshal(v)
			queryParts = append(queryParts, url.QueryEscape(key)+"="+url.QueryEscape(string(jsonBytes)))
		default:
			jsonBytes, _ := json.Marshal(value)
			queryParts = append(queryParts, url.QueryEscape(key)+"="+url.QueryEscape(string(jsonBytes)))
		}
	}
	queryString := ""
	if len(queryParts) > 0 {
		queryString = "?" + strings.Join(queryParts, "&")
	}
	fullURL := baseURL + queryString

	client := &http.Client{Timeout: 30 * time.Second}
	httpReq, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return err
	}

	httpReq.Header.Set("Access-Token", accessToken)

	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, result)
	return err
}

func DoPostRequest(accessToken, baseURL string, params url.Values, result interface{}) error {
	fullURL := baseURL + "?" + params.Encode()

	client := &http.Client{Timeout: 30 * time.Second}
	httpReq, err := http.NewRequest("POST", fullURL, nil)
	if err != nil {
		return err
	}

	httpReq.Header.Set("Access-Token", accessToken)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, result)
	return err
}

func DoPostJSONRequest(accessToken, baseURL string, jsonBody interface{}, result interface{}) error {
	var jsonBytes []byte
	var err error

	if bytes, ok := jsonBody.([]byte); ok {
		jsonBytes = bytes
	} else {
		jsonBytes, err = json.Marshal(jsonBody)
		if err != nil {
			return err
		}
	}

	client := &http.Client{Timeout: 30 * time.Second}
	httpReq, err := http.NewRequest("POST", baseURL, bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Access-Token", accessToken)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, result)
	return err
}

// DoDownloadRequest 下载文件到 writer。
// 返回值：HTTP 状态码、响应头（用于提取 X-Request-ID 等追踪字段）、错误。
func DoDownloadRequest(accessToken, baseURL string, params url.Values, writer io.Writer) (int, http.Header, error) {
	fullURL := baseURL + "?" + params.Encode()

	client := &http.Client{Timeout: 60 * time.Second}
	httpReq, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return 0, nil, err
	}

	httpReq.Header.Set("Access-Token", accessToken)

	resp, err := client.Do(httpReq)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return resp.StatusCode, resp.Header, err
	}

	return resp.StatusCode, resp.Header, nil
}

type MultipartPostOption struct {
	Field    string
	Value    string
	FileName string
	FileData []byte
}

func DoMultipartPost(accessToken, baseURL string, options []MultipartPostOption, timeout time.Duration, result interface{}) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, opt := range options {
		if opt.FileData != nil && opt.FileName != "" {
			part, _ := writer.CreateFormFile(opt.Field, opt.FileName)
			part.Write(opt.FileData)
		} else {
			writer.WriteField(opt.Field, opt.Value)
		}
	}
	writer.Close()

	client := &http.Client{Timeout: timeout}
	httpReq, err := http.NewRequest("POST", baseURL, body)
	if err != nil {
		return err
	}

	httpReq.Header.Set("Access-Token", accessToken)
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBody, result)
	return err
}

package kubesphere

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	netURL "net/url"
	"titan-container-platform/config"

	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("kubesphere")

var (
	serverURL = ""
	account   = ""
	password  = ""
	cluster   = ""

	token = ""
)

// Init initializes the base URL for the application.
func Init(cfg *config.KubesphereAPIConfig) {
	serverURL = cfg.URL
	account = cfg.UserName
	password = cfg.Password
	cluster = cfg.Cluster

	token = getToken()
	// TODO 定时获取token(2小时过期)
}

// func test() {
// 	userName := "cosmos12345670000002"
// 	// err := CreateUserAccount(userName)
// 	// if err != nil {
// 	// 	log.Errorf("CreateUserAccount: %s", err.Error())
// 	// 	return
// 	// }

// 	order := "order00000003"
// 	_, err := createUserSpace(order, userName)
// 	if err != nil {
// 		log.Errorf("CreateUserSpace: %s", err.Error())
// 		return
// 	}

// 	time.Sleep(1 * time.Second)
// 	_, err = changeWorkspaceMembers(order, userName)
// 	if err != nil {
// 		log.Errorf("changeWorkspaceMembers: %s", err.Error())
// 	}
// 	_, err = createUserResourceQuotas(order, 2, 2, 10)
// 	if err != nil {
// 		log.Errorf("CreateUserResourceQuotas: %s", err.Error())
// 	}
// }

func getToken() string {
	data := netURL.Values{}
	data.Set("grant_type", "password")
	data.Set("username", account)
	data.Set("password", password)
	data.Set("client_id", "kubesphere")
	data.Set("client_secret", "kubesphere")

	client := &http.Client{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", serverURL, "/oauth/token"), nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Body = io.NopCloser(strings.NewReader(data.Encode()))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	log.Infof("Response Status: %s\n", resp.Status)
	log.Infof("Response Body: %s\n", string(body))

	var tokenResp tokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	return tokenResp.AccessToken
}

func doRequest(method, path string, body interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s%s", serverURL, path)

	var req *http.Request
	var err error

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

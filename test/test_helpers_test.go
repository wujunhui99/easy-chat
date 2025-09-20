package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	testPassword   = "vk991112"
	defaultDevice  = "mobile"
	defaultDevName = "Pixel 7"
)

var (
	httpClient     = &http.Client{Timeout: 5 * time.Second}
	userAPIBase    = envOrDefault("USER_API_BASE_URL", "http://127.0.0.1:8888")
	socialAPIBase  = envOrDefault("SOCIAL_API_BASE_URL", "http://127.0.0.1:8881")
	errNotRegister = errors.New("account not registered")
)

type apiEnvelope struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

type apiError struct {
	Code int
	Msg  string
}

func (e *apiError) Error() string {
	return fmt.Sprintf("api error: code=%d msg=%s", e.Code, e.Msg)
}

func envOrDefault(key, def string) string {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		return val
	}
	return def
}

func doJSONRequest(tb testing.TB, method, url, token string, payload any) (*apiEnvelope, error) {
	if tb != nil {
		tb.Helper()
	}

	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("marshal payload: %w", err)
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var envelope apiEnvelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, fmt.Errorf("decode envelope: %w (body=%s)", err, string(raw))
	}

	if envelope.Code != 200 {
		return &envelope, &apiError{Code: envelope.Code, Msg: envelope.Msg}
	}

	return &envelope, nil
}

type tokenPayload struct {
	Token  string `json:"token"`
	Expire int64  `json:"expire"`
}

type userInfo struct {
	Id       string `json:"id"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
}

type userInfoData struct {
	Info userInfo `json:"info"`
}

type userContext struct {
	Phone    string
	Nickname string
	Token    string
	UserID   string
}

func ensureTestUsers(t *testing.T) map[string]*userContext {
	t.Helper()

	users := make(map[string]*userContext, 16)
	for i := 1; i <= 16; i++ {
		label := fmt.Sprintf("test%02d", i)
		ctx := bootstrapUser(t, label)
		users[label] = ctx
	}
	return users
}

func bootstrapUser(t *testing.T, label string) *userContext {
	t.Helper()

	loginData, info, err := loginAndFetchInfo(label)
	if err == nil {
		return &userContext{Phone: label, Nickname: label, Token: loginData.Token, UserID: info.Id}
	}

	if errors.Is(err, errNotRegister) {
		if _, regErr := registerUser(t, label); regErr != nil {
			t.Fatalf("register %s: %v", label, regErr)
		}
		loginData, info, err = loginAndFetchInfo(label)
		if err != nil {
			t.Fatalf("login after register %s: %v", label, err)
		}
		return &userContext{Phone: label, Nickname: label, Token: loginData.Token, UserID: info.Id}
	}

	t.Fatalf("login %s: %v", label, err)
	return nil
}

func loginAndFetchInfo(label string) (tokenPayload, userInfo, error) {
	loginData, err := loginUser(label)
	if err != nil {
		if apiErr, ok := err.(*apiError); ok && apiErr.Msg == "手机号没有注册" {
			return tokenPayload{}, userInfo{}, errNotRegister
		}
		return tokenPayload{}, userInfo{}, err
	}

	info, err := fetchUserInfo(loginData.Token)
	if err != nil {
		return tokenPayload{}, userInfo{}, err
	}

	return loginData, info, nil
}

func registerUser(t *testing.T, label string) (tokenPayload, error) {
	t.Helper()

	payload := map[string]any{
		"phone":      label,
		"password":   testPassword,
		"nickname":   label,
		"sex":        1,
		"avatar":     "test.png",
		"deviceType": defaultDevice,
		"deviceName": defaultDevName,
	}

	url := fmt.Sprintf("%s/v1/user/register", userAPIBase)
	envelope, err := doJSONRequest(t, http.MethodPost, url, "", payload)
	if err != nil {
		return tokenPayload{}, err
	}

	var data tokenPayload
	if err := json.Unmarshal(envelope.Data, &data); err != nil {
		return tokenPayload{}, fmt.Errorf("decode register data: %w", err)
	}
	return data, nil
}

func loginUser(label string) (tokenPayload, error) {
	payload := map[string]any{
		"phone":      label,
		"password":   testPassword,
		"deviceType": defaultDevice,
		"deviceName": defaultDevName,
	}

	url := fmt.Sprintf("%s/v1/user/login", userAPIBase)
	envelope, err := doJSONRequest(nil, http.MethodPost, url, "", payload)
	if err != nil {
		return tokenPayload{}, err
	}

	var data tokenPayload
	if err := json.Unmarshal(envelope.Data, &data); err != nil {
		return tokenPayload{}, fmt.Errorf("decode login data: %w", err)
	}
	return data, nil
}

func fetchUserInfo(token string) (userInfo, error) {
	url := fmt.Sprintf("%s/v1/user/me", userAPIBase)
	envelope, err := doJSONRequest(nil, http.MethodGet, url, token, nil)
	if err != nil {
		return userInfo{}, err
	}

	var data userInfoData
	if err := json.Unmarshal(envelope.Data, &data); err != nil {
		return userInfo{}, fmt.Errorf("decode user info: %w", err)
	}

	return data.Info, nil
}

// ensure the compiled package depends on testing for shared helpers.

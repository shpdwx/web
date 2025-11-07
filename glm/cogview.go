package glm

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/shpdwx/web/conf"
)

type CogViewImageResp struct {
	Created int64    `json:"created"`
	Data    []data   `json:"data"`
	Filter  []filter `json:"content_filter"`
}

type data struct {
	Url string `json:"url"`
}

type filter struct {
	Role  string `json:"role"`
	Level int8   `json:"level"`
}

type cogviewReq struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Size   string `json:"size"`
}

func Image(ctx context.Context, c conf.CogView, txt string) (*CogViewImageResp, error) {

	var (
		rawUrl = c.Api
		token  = c.Token
	)

	if rawUrl == "" || token == "" {
		return nil, errors.New("请先配置请求信息")
	}

	// 请求体
	cvq := cogviewReq{Size: "1024x1024"}
	cvq.Model = c.Model
	cvq.Prompt = txt

	b, err := json.Marshal(cvq)
	if err != nil {
		return nil, err
	}

	payload := strings.NewReader(string(b))

	req, err := http.NewRequestWithContext(ctx, "POST", rawUrl, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	// 请求
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp CogViewImageResp

	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

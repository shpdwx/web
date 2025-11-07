package s3

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

type ImageMeta struct {
	Name        string
	ContentType string
	Size        int64
	Origin      string
	RequestId   string
	Desc        string
	Byte        []byte
}

func Fetch(rawUrl string) (*ImageMeta, error) {
	resp := &ImageMeta{}

	if rawUrl == "" {
		return resp, nil
	}

	// 获取图片名
	u, err := url.Parse(rawUrl)
	if err != nil {
		return resp, err
	}

	// 判断文件类型
	isImage := false
	ext := path.Ext(u.Path)
	allows := [...]string{".png", ".jpg", ".jpeg", ".gif"}

	for _, t := range allows {
		if t == ext {
			isImage = true
			break
		}
	}

	if !isImage {
		return resp, errors.New("文件不是图片类型")
	}

	r, err := http.Get(rawUrl)
	if err != nil {
		return resp, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("下载图片失败 Http Status Code 是 %d", r.StatusCode)
	}

	// 读取图片到缓冲区
	buf := new(bytes.Buffer)
	size, err := buf.ReadFrom(r.Body)
	if err != nil {
		return resp, fmt.Errorf("读取数据到缓冲区失败 %v", err)
	}

	// 图片信息
	resp.Size = size
	resp.Name = path.Base(u.Path)
	resp.ContentType = r.Header.Get("Content-Type")
	resp.Origin = "web"
	resp.Byte = buf.Bytes()

	if resp.ContentType == "" {
		resp.ContentType = "application/octet-stream"
	}

	return resp, nil
}

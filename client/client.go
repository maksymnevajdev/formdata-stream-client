package client

import (
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"github.com/valyala/fasthttp"
)

type HttpClient struct {
	client *fasthttp.Client
}

func NewClient() *HttpClient {
	dialer := &fasthttp.TCPDialer{
		Concurrency:      4096,
		DNSCacheDuration: time.Hour,
	}

	return &HttpClient{client: &fasthttp.Client{
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
		Dial:                          dialer.Dial,
	},
	}
}

func (h *HttpClient) SendFile(url string, filename string, file io.Reader) ([]byte, int, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodPost)

	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()
		part, err := m.CreateFormFile("file", filename)

		if err != nil {
			fmt.Println(err)
			return
		}
		if _, err = io.Copy(part, file); err != nil {
			return
		}
	}()

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetBodyStream(r, -1)
	req.Header.Set("Content-Type", m.FormDataContentType())

	fmt.Println(m.FormDataContentType())
	err := h.client.Do(req, resp)
	if err != nil {
		fmt.Printf("Возникла ошибка отправки файла %s на url %s Текст ошибки %s ", filename, url, err)
		return resp.Body(), resp.StatusCode(), err
	}
	statusCode := resp.StatusCode()

	switch statusCode {
	case fasthttp.StatusOK:
		return resp.Body(), statusCode, nil
	default:
		// Ошибка при коде отличном от 200
		return resp.Body(), statusCode, fmt.Errorf("Возникла ошибка отправки файла %s на url %s Сервер ответил %d %s", filename, url, statusCode, resp.Body())
	}
}

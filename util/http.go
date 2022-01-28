package util

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GET(proxy string, webUrl string, headers []string) ([]byte, error) {

	request, err := http.NewRequest("GET", webUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Connection", "keep-alive")

	for _, v := range headers {
		if strings.ContainsAny(v, ":") {
			kv := strings.Split(v, ":")
			request.Header.Set(kv[0], kv[1])
		}
	}

	//build a client
	client := func() *http.Client {
		if proxy == "" {
			return http.DefaultClient
		}
		proxy, _ := url.Parse(proxy)
		tr := &http.Transport{
			Proxy:           http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		return &http.Client{
			Transport: tr,
			Timeout:   time.Second * 10, //超时时间
		}

	}()

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

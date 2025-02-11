package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

type TTS struct{}

func NewTTS() *TTS {
	return &TTS{}
}

func (t *TTS) Speak(text, language string) ([]byte, error) {
	return t.request(
		fmt.Sprintf("https://translate.google.com/translate_tts?ie=UTF-8&total=1&client=tw-ob&q=%s&tl=%s", text, language),
		10*time.Second,
		1*time.Second,
		2,
		false,
		nil,
	)
}

func (t *TTS) request(url string, timeout, retryInterval time.Duration, maxRetries int, trimSpace bool, customTransport http.RoundTripper, headers ...map[string]string) ([]byte, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true, MinVersion: tls.VersionTLS12}
	client := http.Client{Timeout: timeout, Transport: customTransport}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for _, header := range headers {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}

	err = fmt.Errorf("GET request failed")
	for retries := 0; retries < maxRetries; retries++ {
		resp, err := client.Do(req)
		if err != nil {
			time.Sleep(retryInterval)
			continue
		}

		var buf bytes.Buffer
		_, _ = buf.ReadFrom(resp.Body)
		resp.Body.Close()
		b := buf.Bytes()

		if trimSpace {
			for i := 0; i < len(b); i++ {
				if b[i] == ' ' {
					b = append(b[:i], b[i+1:]...)
				}
			}
		}

		return b, nil
	}

	return nil, err
}

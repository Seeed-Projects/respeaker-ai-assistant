package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Gemini struct {
	apiKey string

	ctx    context.Context
	client *genai.Client
	model  *genai.GenerativeModel
}

type GeminiAttachment struct {
	FileName string
	MIMEType string
	Data     []byte
}

func NewGemini(apiKey, model string) (*Gemini, error) {
	obj := &Gemini{
		ctx:    context.Background(),
		apiKey: apiKey,
	}

	client, err := genai.NewClient(obj.ctx, option.WithAPIKey(obj.apiKey))
	if err != nil {
		return nil, err
	}

	obj.client = client
	obj.model = client.GenerativeModel(model)

	return obj, nil
}

func (g *Gemini) Generate(prompt string, attachments ...GeminiAttachment) (string, error) {
	fileDataArr := make([]genai.Part, len(attachments))
	for idx, attachment := range attachments {
		ioReader := io.NopCloser(bytes.NewReader(attachment.Data))
		file, err := g.client.UploadFile(g.ctx, "", ioReader, &genai.UploadFileOptions{
			MIMEType:    attachment.MIMEType,
			DisplayName: attachment.FileName,
		})
		if err != nil {
			return "", err
		}
		fileDataArr[idx] = genai.FileData{URI: file.URI}
	}

	resp, err := g.model.GenerateContent(g.ctx, append([]genai.Part{genai.Text(prompt)}, fileDataArr...)...)
	if err != nil {
		return "", err
	}

	var responseText string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				responseText += fmt.Sprintf("%v", part)
			}
		}
	}

	return strings.TrimSuffix(responseText, "\n"), nil
}

func (g *Gemini) Close() {
	g.client.Close()
}

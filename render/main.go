package main

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

type response struct {
	Body            string            `json:"body"`
	Headers         map[string]string `json:"headers"`
	Cookies         []string          `json:"cookies"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
	StatusCode      int               `json:"statusCode"`
}

func handler(ctx context.Context) (response, error) {
	r := response{
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		Cookies:         []string{},
		IsBase64Encoded: false,
	}
	b, err := createBody()
	if err != nil {
		r.StatusCode = http.StatusInternalServerError
		return r, err
	}
	r.Body = b
	r.StatusCode = http.StatusOK
	return r, nil
}

type data struct {
	Title   string
	Message string
	Time    time.Time
}

func createBody() (string, error) {
	t, err := template.New("tmpl").Parse(html)
	if err != nil {
		return "html parse error", err
	}
	var w bytes.Buffer
	d := data{
		Title:   "サーバーレスCMS",
		Message: "ようこそ",
		Time:    time.Now(),
	}
	if err := t.Execute(&w, d); err != nil {
		return "html template execute error", err
	}
	return w.String(), nil
}

const html = `<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
</head>
<body>
<h1>{{.Title}}</h1>
<p>{{.Message}}</p>
<p>米国の現在時間 {{.Time.Format "2006/1/2 15:04:05"}}</p>
</body>
</html>
`

package connector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type requestType string

const (
	Get  requestType = "GET"
	Post requestType = "POST"
)

type Header struct {
	Name  string
	Value string
}

type Request struct {
	Endpoint string
	Type     *requestType
	Headers  []Header
	Body     any
}

func (h Header) Register(req *http.Request) {
	req.Header.Add(h.Name, h.Value)
}

func (req Request) Make() Response {
	var query *http.Request

	content, err := json.Marshal(req.Body)
	if err != nil {
		return BuildError(err, http.StatusBadRequest)
	}

	if query, err = http.NewRequest(
		fmt.Sprint(req.Type),
		req.Endpoint,
		bytes.NewBuffer(content),
	); err != nil {
		return BuildError(err, http.StatusBadRequest)
	}

	for _, header := range req.Headers {
		header.Register(query)
	}

	query.Header.Add("Content-Type", "application/json")

	res, err := (&http.Client{}).Do(query)
	if err != nil {
		return BuildError(err, http.StatusBadRequest)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return BuildError(err, http.StatusBadRequest)
	}

	return Response{
		Data:   body,
		Status: res.StatusCode,
	}
}

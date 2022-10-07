package conn

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type requestType string

var (
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
	var payload []byte

	if req.Body != nil {
		var err error
		payload, err = json.Marshal(req.Body)
		if err != nil {
			return BuildError(err, http.StatusBadRequest)
		}
	}

	query, err := http.NewRequest(
		string(*req.Type),
		req.Endpoint,
		bytes.NewBuffer(payload),
	)

	if err != nil {
		return BuildError(err, http.StatusInternalServerError)
	}

	for _, header := range req.Headers {
		header.Register(query)
	}

	query.Header.Add("Content-Type", "application/json")

	res, err := (&http.Client{}).Do(query)
	if err != nil {
		return BuildError(err, res.StatusCode)
	}

	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return BuildError(err, res.StatusCode)
	}

	return Response{
		Data:   string(content),
		Status: res.StatusCode,
	}
}

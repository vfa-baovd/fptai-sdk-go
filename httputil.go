package sdk

import (
	"bytes"
	"encoding/json"
	// "errors"
	"io/ioutil"
	"net/http"
	"time"
)

type param struct {
	Method      string
	URI         string
	ContentType string
	Data        []byte
}

func request(p *param) ([]byte, error) {
	req, err := http.NewRequest(p.Method, p.URI, bytes.NewReader(p.Data))
	if err != nil {
		return nil, err
	}

	if p.ContentType != "" {
		req.Header.Set("Content-Type", p.ContentType)
	}

	c := &http.Client{
		Timeout: time.Duration(TIMEOUT) * time.Second,
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var err Error
		json.Unmarshal(body, &err)
		return nil, err
	}

	return body, nil
}

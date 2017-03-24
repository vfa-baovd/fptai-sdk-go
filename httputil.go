package sdk

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
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
		log.Error("failed to create new http request: ", err)
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
		log.Errorf("failed to request to FPT.AI. Error: %s\n", err.Error())
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("failed to read FPT.AI response body. Error: %s. Body: %s", err.Error(), string(body))
		return nil, err
	}
	defer res.Body.Close()
	log.Debugf("FPTAI response body: Status %d. Body: %s\n", res.StatusCode, string(body))

	if res.StatusCode != 200 {
		var err Error
		json.Unmarshal(body, &err)
		log.Error("failed to request to FPT.AI: ", err)
		return nil, err
	}

	return body, nil
}

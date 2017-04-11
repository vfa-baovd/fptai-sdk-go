package sdk

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
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
		return nil, errors.Wrapf(err, "create new request failed. param = %+v\n", p)
	}

	if p.ContentType != "" {
		req.Header.Set("Content-Type", p.ContentType)
	}

	c := &http.Client{
		Timeout: time.Duration(TIMEOUT) * time.Second,
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Do request failed, param = %+v\n", p)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "ReadAll failed. param = %+v; body = %s\n", p, string(body))
	}
	defer res.Body.Close()
	log.Debugf("Request to FPTAI:\nRequest:\tParam = %+v\nResponse: \tStatus %d. Body: %s\n", p, res.StatusCode, string(body))

	if res.StatusCode != 200 {
		var err error
		json.Unmarshal(body, &err)
		log.Error("failed to request to FPT.AI: ", err)
		return nil, err
	}

	return body, nil
}

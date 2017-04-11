package sdk

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

type Application struct {
	client *Client
	code   string
	token	string
}

func (irs IntentResponses) Top() *IntentResponse {
	var top_confidence float64
	var top_intent IntentResponse

	for _, ir := range irs.Data {
		confidence, err := strconv.ParseFloat(ir.Confidence, 64)
		if err != nil {
			// ignore this one
			log.Errorf("failed to parse string to float. Error: %s. String: %s\n", err.Error(), ir.Confidence)
		}
		if confidence > top_confidence {
			top_confidence = confidence
			top_intent = ir
		}
	}

	return &top_intent
}

func (a *Application) Recognize(text string) (*IntentResponse, error) {
	var irs IntentResponses

	v := url.Values{}
	v.Set("application_code", a.code)
	v.Set("api_key", a.token) // Why don't just call it token?
	v.Set("content", text)
	v.Set("type", "intent")
	p := param{
		Method: "POST",
		URI:    fmt.Sprintf("%s/recognition?%s", FPTAIEndpoint, v.Encode()),
	}

	resp, err := request(&p)
	if err != nil {
		log.Error("failed to request Recognition to FPT.AI: ", err)
		return nil, err
	}

	if err := json.Unmarshal(resp, &irs); err != nil {
		log.Error("failed to unmarshal: ", string(resp), err)
		return nil, err
	}

	return irs.Top(), nil
}

func (a *Application) Train() error {
	log.Info("training...")
	v := url.Values{}
	v.Set("application_code", a.code)
	v.Set("type", "intent")
	v.Set("api_key", a.token)
	v.Set("session_id", a.client.SessionID())

	p := param{
		Method: "GET",
		URI:    fmt.Sprintf("%s/train?%s", FPTAIEndpoint, v.Encode()),
	}

	resp, err := request(&p)
	if err != nil {
		log.Error("failed to train: ", string(resp))
		return err
	}

	log.Info("training done")
	return nil
}

// Intents returns all intents of the application
func (a *Application) Intents() (intents []*Intent, err error) {
	v := url.Values{}
	v.Set("session_id", a.client.SessionID())
	v.Set("application_code", a.code)

	p := param{
		Method: "GET",
		URI:    fmt.Sprintf("%s/intent_man?%s", PrincipalEndpoint, v.Encode()),
	}

	resp, err := request(&p)
	if err != nil {
		log.Error("failed to get all intents: ", string(resp))
		return intents, err
	}

	tmp := struct {
		Total int       `json:"total"`
		Data  []*Intent `json:"data"`
	}{}
	if err := json.Unmarshal(resp, &tmp); err != nil {
		log.Error("failed to unmarshal all intents: ", err, string(resp))
		return intents, err
	}

	return tmp.Data, nil
}

func (a *Application) DeleteIntent(intent Intent) error {
	v := url.Values{}
	v.Set("session_id", a.client.SessionID())
	v.Set("intent_code", intent.Code)

	p := param{
		Method: "DELETE",
		URI:    fmt.Sprintf("%s/intent_man?%s", PrincipalEndpoint, v.Encode()),
	}

	_, err := request(&p)
	if err != nil {
		log.Error("failed to delete intent: ", err)
		return  err
	}

	return nil
}

func (a *Application) CreateIntent(label, description string) (*Intent, error) {
	v := url.Values{}
	v.Set("session_id", a.client.SessionID())
	v.Set("application_code", a.code)
	v.Set("label", label)
	v.Set("description", description)

	p := param{
		Method: "POST",
		URI:    fmt.Sprintf("%s/intent_man?%s", PrincipalEndpoint, v.Encode()),
	}

	resp, err := request(&p)
	if err != nil {
		log.Error("failed to create new intent: ", err)
		return nil, err
	}

	var i Intent
	if err := json.Unmarshal(resp, &i); err != nil {
		log.Error("failed to unmarshal create intent response: ", err.Error(), string(resp))
		return nil, err
	}

	return &i, nil
}

func (a *Application) AddSampleByCode(code, sample string) error {
	v := url.Values{}
	v.Set("session_id", a.client.SessionID())
	v.Set("intent_code", code)
	v.Set("content", sample)
	v.Set("application_code", a.code)

	p := param{
		Method: "POST",
		URI: fmt.Sprintf("%s/sample_intent_man?%s", PrincipalEndpoint, v.Encode()),
	}

	_, err := request(&p)
	if err != nil {
		log.Error("failed to add sample: ", err)
		return err
	}

	return nil
}

func (a *Application) AddSampleByLabel(label, sample string) error {
	v := url.Values{}
	v.Set("session_id", a.client.SessionID())
	v.Set("intent_label", label)
	v.Set("content", sample)
	v.Set("application_code", a.code)

	p := param{
		Method: "POST",
		URI: fmt.Sprintf("%s/sample_intent_man?%s", PrincipalEndpoint, v.Encode()),
	}

	_, err := request(&p)
	if err != nil {
		log.Error("failed to add sample: ", err)
		return err
	}

	return nil
}

// TODO: This is temporary. Waiting for API Add Sample by Label
func (a *Application) LabelCodeMap() (map[string]string, error) {
	m := make(map[string]string)

	intents, err := a.Intents()
	if err != nil {
		log.Error("failed to making map of label to code: ", err)
		return m, err
	}

	for _, intent := range intents {
		m[intent.Label] = intent.Code
	}

	return m, nil
}

// ExtractEntities returns entities it found in the input text
// API:
	// http://{{package_domain}}/recognition?
	// application_code={{application_code}}&
	// content=I want to change HubID because I couldn't login.&
	// type=entity
func (a *Application) ExtractEntities(text string) ([]Entity, error) {
	v := url.Values{}
	v.Set("application_code", a.code)
	v.Set("content", text)
	v.Set("type", "entity")
	v.Set("api_key", a.token)

	p := param{
		Method: "POST",
		URI: fmt.Sprintf("%s/recognition?%s", FPTAIEndpoint, v.Encode()),
	}

	resp, err := request(&p)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	log.Error(string(resp))

	var er EntityResponse
	if err := json.Unmarshal(resp, &er); err != nil {
		return nil, errors.Wrap(err, "Unmarshal failed. Response = " + string(resp))
	}

	return er.Data, nil
}
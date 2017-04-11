package sdk

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
)

const SESSION_TIMEOUT float64 = 25.0 // minutes

type Client struct {
	sessionID string
	username  string
	password  string
	lastCall 	time.Time

	Timeout int
}

func NewClient(username, password string) (*Client, error) {
	var c Client = Client{
		username: username,
		password: password,
		lastCall: time.Now(),
		Timeout:  TIMEOUT,
	}

	return &c, nil
}

func (c *Client) SessionID() string {
	defer func(c *Client){c.lastCall = time.Now()}(c)

	if c.sessionID == "" || time.Since(c.lastCall).Minutes() > SESSION_TIMEOUT {
		c.getSessionID()
	}

	return c.sessionID
}

func (c *Client) getSessionID() error {
	p := param{
		Method: "POST",
		URI:    fmt.Sprintf("%s/%s?username=%s&password=%s", OpenFPTEndpoint, SessionPath, c.username, c.password),
	}

	resp, err := request(&p)
	if err != nil {
		log.Errorf("failed to request %s. Error: %s", p.URI, err.Error())
		return err
	}

	c.sessionID = string(resp)
	return nil
}

func (c *Client) GetApp(code, token string) *Application {
	return &Application{
		client: c,
		code:   code,
		token: token,
	}
}

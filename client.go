package sdk

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
)

const SESSION_TIMEOUT float64 = 25.0 // minutes

type client struct {
	sessionID string
	username  string
	password  string
	lastCall 	time.Time

	Timeout int
}

func NewClient(username, password string) (*client, error) {
	var c client = client{
		username: username,
		password: password,
		lastCall: time.Now(),
		Timeout:  TIMEOUT,
	}

	return &c, nil
}

func (c *client) SessionID() string {
	defer func(c *client){c.lastCall = time.Now()}(c)

	if c.sessionID == "" || time.Since(c.lastCall).Minutes() > SESSION_TIMEOUT {
		c.getSessionID()
	}

	return c.sessionID
}

func (c *client) getSessionID() error {
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

func (c *client) GetApp(code, token string) *Application {
	return &Application{
		client: c,
		code:   code,
		token: token,
	}
}

package sdk

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

type client struct {
	sessionID string
	username  string
	password  string

	Timeout int
}

func NewClient(username, password string) (*client, error) {
	var c client = client{
		username: username,
		password: password,
		Timeout:  TIMEOUT,
	}

	return &c, nil
}

func (c *client) SessionID() string {
	if c.sessionID == "" {
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

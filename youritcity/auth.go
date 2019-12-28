package youritcity

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"

	"golang.org/x/net/publicsuffix"
)

// SDKVersion the version of this SDK
const SDKVersion = "0.1.1"

// SDKName is used for identification
const SDKName = "go-sdk_" + SDKVersion

// DefaultYourITcityURI is the URI to access yourITcity services change it for local develepment
const DefaultYourITcityURI = "https://yicauth.cleverapps.io"

// const YourITcityURI = "https://localhost:2015"

// App contains information to connect to yourITcity services
type App struct {
	Token string `json:"app_token"`
	Name  string `json:"name"`
	Valid bool   `json:"valid"`
}

// Roles is a list of name with their authorization (r: read, w: write)
type Roles map[string]map[string]string

// Connection to yourITcity services
type Connection struct {
	client *http.Client
	ConnectionConfig
}

// NewConnection create a new connection to yourITcity services
func NewConnection() *Connection {
	return NewConnectionWithConfig(ConnectionConfig{})
}

// ConnectionConfig to yourITcity services
type ConnectionConfig struct {
	URI         string
	InsecureSSL bool
}

// NewConnectionWithConfig create a new connection to yourITcity services
func NewConnectionWithConfig(config ConnectionConfig) *Connection {
	if config.URI == "" {
		config.URI = DefaultYourITcityURI
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't create cookie jar: %s\n", err)
		os.Exit(1)
	}

	conn := &Connection{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: config.InsecureSSL},
			},
			Jar: jar,
		},
		ConnectionConfig: config,
	}
	return conn
}

// Client return the *http.Client you can use to call yourITcity with authentifacation already handle
func (conn *Connection) Client() *http.Client {
	return conn.client
}

// Signup to yourITcity
func (conn *Connection) Signup(email string, appType string) (*App, error) {
	return conn.login(email, appType, conn.URI+"/auth/signup")
}

// Login to yourITcity
func (conn *Connection) Login(email string, appType string) (*App, error) {
	return conn.login(email, appType, conn.URI+"/auth/login")
}

func (conn *Connection) login(email, appType string, url string) (*App, error) {
	auth := struct {
		Email string `json:"email"`
		Name  string `json:"name,omitempty"`
		Type  string `json:"type"`
	}{
		Email: email,
		Type:  appType,
		Name:  SDKName,
	}
	body, err := json.Marshal(auth)
	if err != nil {
		return nil, err
	}
	resp, err := conn.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("%s: %s", resp.Status, b)
	}

	app := &App{
		Name: auth.Name,
	}
	err = json.NewDecoder(resp.Body).Decode(app)
	if err != nil {
		return nil, err
	}
	return app, nil
}

// GetRoles find every possibles authorization roles
func (conn *Connection) GetRoles() (Roles, error) {
	resp, err := conn.client.Get(conn.URI + "/auth/roles")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("%s: %s", resp.Status, b)
	}

	var roles Roles
	err = json.NewDecoder(resp.Body).Decode(&roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// Renew the session, by getting an HTTP Authorization header and cookie
func (conn *Connection) Renew(app App) (string, error) {
	body := bytes.NewBufferString("{\"app_token\":\"" + app.Token + "\"}")
	resp, err := conn.client.Post(conn.URI+"/renew", "application/json", body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Renew %v", resp.Status)
	}

	return resp.Header.Get("Authorization"), nil
}

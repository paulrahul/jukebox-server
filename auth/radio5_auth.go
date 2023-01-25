package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	. "github.com/paulrahul/jukebox-server"
)

type Credentials struct {
	UserName      string
	Password      string
	ContributorID string
}

type Radio5Client struct {
	Credentials Credentials
	Cookies     map[string]string
}

type Radio5Auth struct {
	User   User
	Client Radio5Client
}

var r5AuthInstance *Radio5Auth

var BASE_URL string

func Radio5Init() {
	BASE_URL = "https://radiooooo.com/"
}

func readCreds(credentials *Credentials) {
	un := os.Getenv("JUKEBOX_R5UN")
	if un == "" {
		panic("Radio5 username must be set")
	}

	pwd := os.Getenv("JUKEBOX_R5PWD")
	if pwd == "" {
		panic("Radio5 password must be set")
	}

	contID := os.Getenv("RADIO5_CONTRIBUTOR_ID")
	if contID == "" {
		panic("Radio5 contributor ID must be set")
	}
	credentials.ContributorID = contID

	credentials.setCredentials(un, pwd)
}

func GetRadio5Auth() *Radio5Auth {
	log.Debug("GetRadio5Auth called.")

	if r5AuthInstance == nil {
		Radio5Init()
		r5AuthInstance = &Radio5Auth{}
	}

	return r5AuthInstance
}

func (r5 Radio5Auth) Login(w http.ResponseWriter, r *http.Request) {
	log.Debug("Radio5Auth.Login called.")

	// Here redirect to a page which asks the user their Radio5 login.
	// Then redirect here, where we store those credentials in encrypted
	// form.

	// TODO (fix this)
	// For now, just redirect to the auth_callback url for radio5, and in
	// the handler for that, hardcode set the credentials.
	http.Redirect(w, r, "/radio5_auth_callback", http.StatusTemporaryRedirect)
}

func (r5 *Radio5Auth) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Radio5Auth.RedirectHandler called.")

	// Set hardcoded credentials for now.
	// TODO - remove
	readCreds(&r5.Client.Credentials)

	// Call Radio5 login API.
	postBody, err := json.Marshal(map[string]string{
		"email":    r5.Client.Credentials.UserName,
		"password": r5.Client.Credentials.Password,
	})
	if err != nil {
		http.Error(w, "Couldn't login: "+err.Error(), http.StatusInternalServerError)
		return
	}
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(
		BASE_URL+"account/login", "application/json", responseBody)
	if err != nil {
		http.Error(w, "Couldn't login: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Now set the cookies if the request was successful.
	r5.Client.Cookies = make(map[string]string)
	if resp.StatusCode == http.StatusOK {
		cookies := resp.Cookies()

		for _, c := range cookies {
			r5.Client.Cookies[c.Name] = c.Value
		}
	} else {
		http.Error(w, "Couldn't login: "+fmt.Sprint(resp.StatusCode), http.StatusInternalServerError)
		return
	}

	// TODO - fetch the un and pwd from request object.

	http.Redirect(w, r, "/user", http.StatusTemporaryRedirect)
}

func (r Radio5Auth) GetUser() User {
	return r.User
}

func (c *Credentials) setCredentials(u, p string) {
	c.UserName = u
	c.Password = p
}

func (client *Radio5Client) addCredHeaders(req *http.Request) {
	for k, v := range client.Cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}
}

func (client *Radio5Client) Get(url string) (int, string, error) {
	log.WithFields(log.Fields{
		"url": url,
	}).Debug("Issuing GET call")

	request, err := http.NewRequest(http.MethodGet, BASE_URL+url, http.NoBody)
	if err != nil {
		return -1, "", err
	}

	client.addCredHeaders(request)
	resp, err := http.DefaultClient.Do(request)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, "", err
	}

	return resp.StatusCode, string(body), nil
}

func (client *Radio5Client) Post(
	url string, contentType string, data map[string]interface{}) (int, error) {

	log.WithFields(log.Fields{
		"url":         url,
		"contentType": contentType,
		"data":        data,
	}).Debug("Issuing POST call")

	postBody, err := json.Marshal(data)
	if err != nil {
		return -1, err
	}
	responseBody := bytes.NewBuffer(postBody)

	request, err := http.NewRequest(http.MethodPost, BASE_URL+url, responseBody)
	if err != nil {
		return -1, err
	}

	client.addCredHeaders(request)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

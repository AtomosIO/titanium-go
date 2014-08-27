package titanium

import (
	"fmt"
	"github.com/atomosio/common"
	"io"
	"net/http"
	neturl "net/url"
)

type HttpClient struct {
	endpoint string
	token    string
	client   *http.Client
	log      bool

	//URL Related data
	scheme string
	host   string
}

type URL struct {
	neturl.URL
}

// Create a new Oxygen client using HTTP protocol. The client will use the
// specified token for all interactions with the service. An empty token string
// will cause the omission of token cookie, resulting in only public repositories
// being accesible.
func NewHttpClient(endpoint, token string) (client *HttpClient) {
	return newHttpClient(endpoint, token, false)
}

func newHttpClient(endpoint, token string, log bool) (client *HttpClient) {
	urlURL, _ := neturl.Parse(endpoint)

	return &HttpClient{
		token:    token,
		endpoint: endpoint,
		client:   &http.Client{},
		log:      log,

		scheme: urlURL.Scheme,
		host:   urlURL.Host,
	}
}

func (client *HttpClient) StartLogging() *HttpClient {
	client.log = true
	return client
}

func (client *HttpClient) Logf(format string, args ...interface{}) {
	if client.log {
		fmt.Printf(format, args...)
	}
}

func (client *HttpClient) prepHeadRequest(url *URL) (req *http.Request, err error) {
	return client.prepEmptyRequest("HEAD", url)
}

func (client *HttpClient) prepDeleteRequest(url *URL) (req *http.Request, err error) {
	return client.prepEmptyRequest("DELETE", url)
}

func (client *HttpClient) prepGetRequest(url *URL) (req *http.Request, err error) {
	return client.prepEmptyRequest("GET", url)
}

func (client *HttpClient) prepPostRequest(url *URL, body io.Reader) (req *http.Request, err error) {
	return client.prepRequest("POST", url, body)
}

func (client *HttpClient) prepPatchRequest(url *URL, body io.Reader) (req *http.Request, err error) {
	return client.prepRequest("PATCH", url, body)
}

func (client *HttpClient) prepEmptyRequest(method string, url *URL) (req *http.Request, err error) {
	return client.prepRequest(method, url, common.NewEmptyReader())
}

func (client *HttpClient) prepRequest(method string, url *URL, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", client.token)

	client.Logf("HttpClient %s -> %s\n", method, url)

	return req, nil
}

// Did we get a 2XX respond code?
func statusGood(status int) bool {
	return status >= 200 && status <= 299
}

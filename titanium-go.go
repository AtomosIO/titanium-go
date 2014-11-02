package titanium

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	//	"github.com/atomosio/common"
	"io"
	"io/ioutil"
	"net/http"
	neturl "net/url"
)

type HttpClient struct {
	endpoint string
	Token    string
	client   *http.Client
	log      bool

	//URL Related data
	scheme string
	host   string
	path   string
}

type Response struct {
	Code        int     `json:"code"`
	Description string  `json:"description"`
	Errors      []Error `json:"errors,omitempty"`
}

type URL struct {
	neturl.URL
}

const (
	InstancesEndpoint = "instances/"
	ClustersEndpoint  = "clusters/"
	ProjectsEndpoint  = "projects/"
	TokensEndpoint    = "tokens/"
)

type Error struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

// Create a new client using HTTP protocol. The client will use the
// specified token for all interactions with the service. An empty token string
// will cause the omission of token cookie, resulting in only public repositories
// being accesible.
func NewHttpClient(endpoint, token string) (client *HttpClient) {
	return newHttpClient(endpoint, token, false)
}

func newHttpClient(endpoint, token string, log bool) (client *HttpClient) {
	urlURL, _ := neturl.Parse(endpoint)

	return &HttpClient{
		Token:    token,
		endpoint: endpoint,
		client:   &http.Client{},
		log:      log,

		scheme: urlURL.Scheme,
		host:   urlURL.Host,
		path:   urlURL.Path,
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
	return client.prepRequest(method, url, nil)
}

func (client *HttpClient) prepRequest(method string, url *URL, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", client.Token)

	client.Logf("HttpClient %s -> %s\n", method, url)

	return req, nil
}

func (client *HttpClient) do(req *http.Request) (*http.Response, error) {
	return client.client.Do(req)
}

// Did we get a 2XX respond code?
func statusGood(status int) bool {
	return status >= 200 && status <= 299
}

func (client *HttpClient) get(format string, args ...interface{}) (data []byte, err error) {
	url := client.NewURL(fmt.Sprintf(format, args...))

	// Prepare request
	req, err := client.prepGetRequest(url)
	if err != nil {
		client.Logf("Failed PrepRequest: %s\n", err)
		return nil, err
	}

	// Do request
	resp, err := client.do(req)
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (client *HttpClient) doRequestAndReadResponse(req *http.Request) ([]byte, error) {
	resp, err := client.do(req)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (client *HttpClient) patch(jsonVar interface{}, format string, args ...interface{}) (data []byte, err error) {
	url := client.NewURL(fmt.Sprintf(format, args...))

	marshalledData, err := json.Marshal(jsonVar)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(marshalledData)

	// Prepare request
	req, err := client.prepPatchRequest(url, reader)
	if err != nil {
		client.Logf("Failed PrepRequest: %s\n", err)
		return nil, err
	}

	// Do request
	data, err = client.doRequestAndReadResponse(req)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (client *HttpClient) post(jsonVar interface{}, format string, args ...interface{}) (data []byte, err error) {
	url := client.NewURL(fmt.Sprintf(format, args...))

	marshalledData, err := json.Marshal(jsonVar)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(marshalledData)

	// Prepare request
	req, err := client.prepPostRequest(url, reader)
	if err != nil {
		client.Logf("Failed PrepRequest: %s\n", err)
		return nil, err
	}

	// Do request
	data, err = client.doRequestAndReadResponse(req)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (client *HttpClient) NewURL(path string) *URL {
	path = strings.Replace(path, "//", "/", -1)
	return &URL{
		URL: neturl.URL{
			Scheme: client.scheme,
			Host:   client.host,
			Path:   client.path + path,
		},
	}
}

func (client *HttpClient) getAndUnmarshal(addr string, i interface{}) error {
	data, err := client.get(addr)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, i)
	if err != nil {
		return err
	}

	return nil
}

func (client *HttpClient) postAndUnmarshal(addr string, jsonVar interface{}, i interface{}) error {
	data, err := client.post(jsonVar, addr)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, i)
	if err != nil {
		return err
	}

	return nil
}

func (client *HttpClient) patchAndUnmarshal(addr string, jsonVar interface{}, i interface{}) error {
	data, err := client.patch(jsonVar, addr)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, i)
	if err != nil {
		return err
	}

	return nil
}

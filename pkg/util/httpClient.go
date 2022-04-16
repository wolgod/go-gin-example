package util

import (
	"errors"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	httpClient *http.Client
)

const (
	ContentType     = "Content-Type"
	ErrThirdNetwork = "third netword err"
)

// init HTTPClient
func init() {
	httpClient = createHTTPClient()
}

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{}
	return client
}
func SendHttpRequest(method, path, contentType string, reader io.Reader) ([]byte, error) {

	req, err := http.NewRequest(method, path, reader)
	if err != nil {
		logging.Errorf("build http request  : %s", err.Error())
		return nil, errors.New(ErrThirdNetwork)
	}
	if contentType != "" {
		req.Header.Set(ContentType, contentType)
	}
	resp, err := httpClient.Do(req)

	//确保返回的err不为nil,如果resp也不会nil的情况下 body也能关闭
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		logging.Errorf("http request error : %s", err.Error())
		return nil, errors.New(ErrThirdNetwork)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logging.Errorf("http request error when read body: %s", err.Error())
		return nil, errors.New(ErrThirdNetwork)
	}
	if resp.StatusCode == http.StatusOK {
		return body, nil
	} else {
		logging.Errorf("nacos request error : %s,response status code :%d", string(body), resp.StatusCode)
		return nil, errors.New(" handle fail from nacos")
	}
}

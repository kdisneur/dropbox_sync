package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

func POSTWithDataHeadersAndBinary(url string, token string, data map[string]interface{}, content []byte) ([]byte, error) {
	arguments, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "can't encode POST body request")
	}

	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	header.Set("Content-Type", "application/octet-stream")
	header.Set("Dropbox-API-Arg", string(arguments))

	return doPOSTRequestWithBinary(url, header, content)
}

func POSTWithDataHeaders(url string, token string, data map[string]interface{}) ([]byte, error) {
	arguments, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "can't encode POST body request")
	}

	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	header.Set("Dropbox-API-Arg", string(arguments))

	return doPOSTRequestWithBinary(url, header, nil)
}

// POSTWithBody posts data and read the response back. It returns an error when status code is
// greater than or equal to 400
func POSTWithBody(url string, token string, data map[string]interface{}) ([]byte, error) {
	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	header.Set("Content-Type", "application/json")

	return doPOSTRequestWithJSON(url, header, data)
}

// UnuathenticatedPOSTWithBody posts data and read the response back. It returns an error when status code is
// greater than or equal to 400
func UnuathenticatedPOSTWithBody(url string, data map[string]interface{}) ([]byte, error) {
	header := http.Header{}
	header.Set("Content-Type", "application/json")

	return doPOSTRequestWithJSON(url, header, data)
}

func doPOSTRequestWithJSON(url string, headers http.Header, data map[string]interface{}) ([]byte, error) {
	var body []byte
	var err error

	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return nil, errors.Wrap(err, "can't encode POST body request")
		}
	}

	return doPOSTRequestWithBinary(url, headers, body)
}

func doPOSTRequestWithBinary(url string, headers http.Header, data []byte) ([]byte, error) {
	var bodyReader io.Reader

	if data != nil {
		bodyReader = bytes.NewReader(data)
	}

	request, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "can't create new POST request")
	}

	request.Header = headers

	var client http.Client
	response, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "can't execute new POST request")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't read POST response")
	}

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("error when POSTing the request. detail: %s", body)
	}

	return body, nil
}

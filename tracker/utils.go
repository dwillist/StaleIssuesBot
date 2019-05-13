package tracker

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const ReponseError = "error calling %s, got status code %d"

type Callout struct {
}

func (c Callout) Call(endpoint string) (string, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(ReponseError, endpoint, resp.StatusCode)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

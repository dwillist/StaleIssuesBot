package tracker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const ReponseError = "error calling %s, got status code %d"

type Callout struct {
}

func (c Callout) Get(endpoint string) (string, error) {
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

func (c Callout) Post(endpoint, data string) (string, error) {
	return "", nil
}

type Time struct {
}

func (t Time) Time() time.Time {
	return time.Now()
}

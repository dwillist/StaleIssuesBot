package tracker

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const ReponseError = "error calling %s, got status code %d"
const TrackerTokenHeader = "X-TrackerToken"

type Callout struct{}

func (c Callout) Get(endpoint string) ([]byte, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf(ReponseError, endpoint, resp.StatusCode)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return buf, nil
}

func (c Callout) Post(endpoint string, data []byte) ([]byte, error) {
	client := &http.Client{}

	trackerToken := os.Getenv("TRACKER_TOKEN")

	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return []byte{}, err
	}

	request.Header[TrackerTokenHeader] = []string{trackerToken}

	response, err := client.Do(request)
	if err != nil {
		return []byte{}, err
	}

	responseBody := make([]byte, 0)
	if _, err := response.Body.Read(responseBody); err != nil {
		return []byte{}, err
	}

	return responseBody, nil
}

type Time struct {
}

func (t Time) Time() time.Time {
	return time.Now()
}

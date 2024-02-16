package commands

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

var testingToken = ""

const reportEndpoint = "https://go.manytask.org/api/report"

func reportTestResults(token string, task string, userID string, failed bool) error {
	form := url.Values{}
	form.Set("token", token)
	form.Set("task", task)
	form.Set("user_id", userID)

	if failed {
		form.Set("failed", "1")
	}

	var rsp *http.Response
	var err error

	for i := 0; i < 3; i++ {
		rsp, err = http.PostForm(reportEndpoint, form)
		if err != nil {
			log.Printf("retrying report: %v", err)
			continue
		}

		if rsp.StatusCode != 200 {
			err = fmt.Errorf("server returned status %d", rsp.StatusCode)
			log.Printf("retrying report: %v", err)
			continue
		}

		return nil
	}

	return err
}

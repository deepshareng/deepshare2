package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/MISingularity/deepshare2/api"
)

func PostJson(apiPath string, reqData interface{}, ua, dscookie string, respObj interface{}) (time.Duration, error) {
	return requestWithJsonBody("POST", apiPath, reqData, ua, dscookie, respObj)
}

func PutJson(apiPath string, reqData interface{}, ua, dscookie string, respObj interface{}) (time.Duration, error) {
	return requestWithJsonBody("PUT", apiPath, reqData, ua, dscookie, respObj)
}

func requestWithJsonBody(method string, apiPath string, reqData interface{}, ua, dscookie string, respObj interface{}) (time.Duration, error) {
	cli := &http.Client{}
	b, err := json.Marshal(reqData)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("[%s] %s %s\n", method, apiPath, string(b))
	req, err := http.NewRequest(method, apiPath, bytes.NewReader(b))
	if err != nil {
		log.Panic(err)
	}
	if dscookie != "" {
		req.AddCookie(&http.Cookie{Name: api.CookieName, Value: dscookie})
	}
	if ua != "" {
		req.Header.Add("User-Agent", ua)
	}
	start := time.Now()
	resp, err := cli.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Panic(fmt.Sprintln("invalid response:", resp.Status, resp.ContentLength))
	}

	if respObj != nil && len(body) > 0 {
		err := json.Unmarshal(body, respObj)
		if err != nil {
			return 0, err
		}
	}
	return time.Since(start), nil
}

func DeleteRequest(apiPath string) time.Duration {
	cli := &http.Client{}

	req, err := http.NewRequest("DELETE", apiPath, nil)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("[%s] %s\n", "DELETE", apiPath)
	start := time.Now()
	resp, err := cli.Do(req)
	if err != nil {
		log.Panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Panic(fmt.Sprintln("invalid response:", resp.Status, resp.ContentLength))
	}
	return time.Since(start)
}

func GetJsonResponse(apiPath string, ua string, respObj interface{}) (time.Duration, error) {
	cli := &http.Client{}
	req, err := http.NewRequest("GET", apiPath, nil)
	if err != nil {
		log.Panic(err)
	}
	if ua != "" {
		req.Header.Add("User-Agent", ua)
	}

	log.Printf("[%s] %s\n", "GET", apiPath)
	start := time.Now()
	resp, err := cli.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Panic(fmt.Sprintln("invalid response:", resp.Status, resp.ContentLength))
	}

	duration := time.Since(start)
	if respObj != nil && len(body) > 0 {
		err := json.Unmarshal(body, respObj)
		if err != nil {
			return duration, err
		}
	}

	return duration, nil
}

func GetHttp(apiPath string, ua string) (header http.Header, body []byte, duration time.Duration) {
	cli := &http.Client{}

	req, err := http.NewRequest("GET", apiPath, nil)
	if err != nil {
		log.Panic(err)
	}
	if ua != "" {
		req.Header.Add("User-Agent", ua)
	}

	start := time.Now()
	resp, err := cli.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Panic(fmt.Sprintln("invalid response:", resp.Status, resp.ContentLength))
	}
	return resp.Header, body, time.Since(start)
}

func GetHttpWithDSCookie(t *testing.T, apiPath string, ua string, dsCookie string) (header http.Header, body []byte, duration time.Duration) {
	cli := &http.Client{}
	log.Printf("[%s] %s\n", "GET", apiPath)
	req, err := http.NewRequest("GET", apiPath, nil)
	if err != nil {
		log.Panic(err)
	}
	if dsCookie != "" {
		req.AddCookie(&http.Cookie{Name: api.CookieName, Value: dsCookie})
	}
	if ua != "" {
		req.Header.Add("User-Agent", ua)
	}
	start := time.Now()
	resp, err := cli.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Panic(fmt.Sprintln("invalid response:", resp.Status, resp.ContentLength))
	}
	return resp.Header, body, time.Since(start)
}

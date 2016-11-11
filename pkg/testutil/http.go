package testutil

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"io/ioutil"

	"github.com/MISingularity/deepshare2/pkg/log"
)

// Shorthand for creating http request, asking handler to handle it,
// and returning a ResponseRecorder for checking results.
func HandleWithBody(hl http.Handler, method, url, body string) *httptest.ResponseRecorder {
	req := mustNewHTTPRequest(method, url, strings.NewReader(body))
	w := httptest.NewRecorder()
	hl.ServeHTTP(w, req)
	return w
}

func MustHandleWithBodyOK(hl http.Handler, method, url, body string) *httptest.ResponseRecorder {
	w := HandleWithBody(hl, method, url, body)
	if w.Code != http.StatusOK {
		log.Errorf("handle returns non OK code! HTTP code: %d, body: %s", w.Code, string(w.Body.Bytes()))
		panic(w.Code)
	}
	return w
}

// Shorthand for creating http request with header, asking handler to handle it,
// and returning a ResponseRecorder for checking results.
func HandleWithRequestInfo(hl http.Handler, method, url, body string, header map[string]string, remoteAddr string) *httptest.ResponseRecorder {
	req := mustNewHTTPRequest(method, url, strings.NewReader(body))
	req.RemoteAddr = remoteAddr
	for k, v := range header {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	hl.ServeHTTP(w, req)
	return w
}

func MustHandleWithRequestInfoOK(hl http.Handler, method, url, body string, header map[string]string, remoteAddr string) *httptest.ResponseRecorder {
	w := HandleWithRequestInfo(hl, method, url, body, header, remoteAddr)
	if w.Code != http.StatusOK {
		log.Errorf("handle returns non OK code! HTTP code: %d, body: %s", w.Code, string(w.Body.Bytes()))
		panic(w.Code)
	}
	return w
}

func mustNewHTTPRequest(method, urlStr string, body io.Reader) *http.Request {
	r, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		log.Errorf("http.NewRequest failed!\nmethod: %s, url: %s, body: %s, err: %v", method, urlStr, body, err)
		panic(err)
	}
	return r
}

func MockResponse(code map[string]int, body map[string]string) (mockServer *httptest.Server, mockClient *http.Client, requestUrlHistory *(map[string]string), requestBodyHistory *(map[string]string)) {
	index := 0
	requestUrlHis := make(map[string]string)
	requestBodyHis := make(map[string]string)
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqPath := r.URL.Path
		if reqPath == "" {
			reqPath = "*"
		}
		requestUrlHis[reqPath] = r.URL.String()
		if reqBody, err := ioutil.ReadAll(r.Body); err != nil {
			log.Errorf("InAppDataHandler; Read Request body error: %s", err)
			fmt.Fprintln(w, "")
		} else {
			requestBodyHis[reqPath] = string(reqBody)
			respCode := code["*"]
			if v, ok := code[reqPath]; ok {
				respCode = v
			}
			w.WriteHeader(respCode)
			w.Header().Set("Content-Type", "application/json")
			respBody := body["*"]
			if v, ok := body[reqPath]; ok {
				respBody = v
			}
			log.Debugf("Response of mock response: %s", respBody)
			fmt.Fprintln(w, respBody)
			index++
		}
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(mockServer.URL)
		},
	}

	mockClient = &http.Client{Transport: transport}
	requestUrlHistory = &requestUrlHis
	requestBodyHistory = &requestBodyHis
	return
}

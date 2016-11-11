package token

import (
	"log"
	"net/http"
	"regexp"
	"testing"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/pkg/testutil"
)

func TestTokenGET(t *testing.T) {
	reg := regexp.MustCompile(`{\"token\":\"\w{10}\"}`)
	log.Println(reg.Match([]byte(`{"token":"4YPVCilFRe"}`)))
	tests := []struct {
		reqPath      string
		wcode        int
		wbodyPattern string
	}{
		{"/v2/tokens/cookie", http.StatusOK, `{\"token\":\"\w{10}\"}`},
		{"/v2/tokens/", http.StatusNotFound, `{"code":102,"message":"Failed to find resource at given path"}`},
	}

	handler := NewTokenTestHandler(api.TokenPrefix)
	for i, tt := range tests {
		path := "http://example.com" + tt.reqPath
		log.Println(path)
		w := testutil.HandleWithBody(handler, "GET", path, "")
		log.Println(i, w.Code, w.Body.String())
		if w.Code != tt.wcode {
			t.Errorf("#%d token handler responsed wrong code = %d, want = %d\n", i, w.Code, tt.wcode)
		}
		reg := regexp.MustCompile(tt.wbodyPattern)
		if !reg.Match(w.Body.Bytes()) {
			t.Errorf("#%d token handler responsed wrong body = %s, want pattern = %s\n", i, w.Body.String(), tt.wbodyPattern)
		}
	}
}

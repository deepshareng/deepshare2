package attribpush

import (
	"net/http"
	"path"
	"testing"

	"github.com/MISingularity/deepshare2/pkg/storage"
	"github.com/MISingularity/deepshare2/pkg/testutil"
)

func TestAppUrlHandlerPUT(t *testing.T) {
	tests := []struct {
		body  string
		wcode int
	}{
		{`{"url":"testurl1"}` + "\n", http.StatusOK},
	}
	handler := newAppUrlHandler(newSimpleAppToUrl(storage.NewInMemSimpleKV()), "/app/url/")
	for i, tt := range tests {
		url := "http://" + path.Join("example.com", "/app/url/", "testAppID")
		w := testutil.HandleWithBody(handler, "PUT", url, tt.body)
		if w.Code != tt.wcode {
			t.Errorf("#%d failed, code = %d, want = %d\n", i, w.Code, tt.wcode)
		}
	}
}

func TestAppUrlHandlerGET(t *testing.T) {
	tests := []struct {
		putBody string

		wcode int
		wbody string
	}{
		{
			`{"url":"testurl1"}` + "\n",
			200,
			`{"url":"testurl1"}` + "\n",
		},
		{
			`{"url":"testurl2"}` + "\n",
			200,
			`{"url":"testurl2"}` + "\n",
		},
		{
			`{"url":""}` + "\n",
			200,
			`{"url":""}` + "\n",
		},
	}

	handler := newAppUrlHandler(newSimpleAppToUrl(storage.NewInMemSimpleKV()), "/app/url/")
	for i, tt := range tests {
		url := "http://" + path.Join("example.com", "/app/url/", "testAppID")
		testutil.MustHandleWithBodyOK(handler, "PUT", url, tt.putBody)

		w := testutil.HandleWithBody(handler, "GET", url, "")
		if w.Code != tt.wcode {
			t.Errorf("#%d failed, code = %d, want = %d\n", i, w.Code, tt.wcode)
		}
		if w.Body.String() != tt.wbody {
			t.Errorf("#%d failed, body = %s, want = %s\n", i, w.Body, tt.wbody)
		}
	}
}

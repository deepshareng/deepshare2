package dsusage

import (
	"net/http"
	"testing"

	"path"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"github.com/MISingularity/deepshare2/pkg/testutil"
)

func TestUsageDELETEAndGET(t *testing.T) {
	tests := []struct {
		method string
		wcode  int
		wbody  string
	}{
		{
			"DELETE",
			http.StatusOK,
			"",
		},
		{
			"GET",
			http.StatusOK,
			`{"new_install":0,"new_open":0}` + "\n",
		},
	}

	db := storage.NewInMemSimpleKV()
	handler := newUsageHandler(db, api.DSUsagesPrefix)
	urlStr := "http://" + path.Join("example.com", api.DSUsagesPrefix, "testAppID", "testSenderID")
	for i, tt := range tests {
		w := testutil.HandleWithBody(handler, tt.method, urlStr, "")
		if w.Code != tt.wcode {
			t.Errorf("#%d: HTTP status code = %d, want = %d\n", i, w.Code, tt.wcode)
		}
		if string(w.Body.Bytes()) != tt.wbody {
			t.Errorf("#%d: response body = %s, want = %s\n", i, w.Body.Bytes(), tt.wbody)
		}
	}

}

package device_stat

import (
	"net/http"
	"path"
	"testing"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"github.com/MISingularity/deepshare2/pkg/testutil"
)

func TestDeviceGet(t *testing.T) {
	skv := storage.NewInMemSimpleKV()
	ds := NewGeneralDeviceStat(skv, "test")
	prepareDevice(ds, t)

	tests := []struct {
		requestPath string

		wcode int
		wbody string
	}{
		{ //good request, no filters
			api.DeviceStatPrefix + "/all",

			http.StatusOK,
			`{"Value":6}` + "\n",
		},
		{ //good request, filter in "install"
			api.DeviceStatPrefix + "/ios",

			http.StatusOK,
			`{"Value":2}` + "\n",
		},
		{ //wrong path request
			api.DeviceStatPrefix + "/",

			http.StatusNotFound,
			api.ErrPathNotFound.Error() + "\n",
		},
	}

	handler := newDeviceStatHandler(ds)
	// Do some GET requests with different channel IDs and event filters.
	for i, tt := range tests {
		var url string

		url = "http://" + path.Join("exmaple.com", tt.requestPath)
		w := testutil.HandleWithBody(handler, "GET", url, "")

		if w.Code != tt.wcode {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.wcode)
		}

		if string(w.Body.Bytes()) != tt.wbody {
			t.Errorf("#%d: HTTP response body = %q, want = %q", i, string(w.Body.Bytes()), tt.wbody)
		}
	}

}

func prepareDevice(ds DeviceStatService, t *testing.T) {
	testmessages := []messaging.Event{
		messaging.Event{
			AppID:     "app1",
			EventType: "sharelink:/d/",
			Channels:  []string{""},
			SenderID:  "a2",
			UniqueID:  "a1",
			UAInfo:    uainfo.UAInfo{Os: "android"},
			TimeStamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
		},
		messaging.Event{
			AppID:     "app1",
			EventType: "sharelink:/d/",
			Channels:  []string{""},
			SenderID:  "a2",
			UniqueID:  "a3",
			UAInfo:    uainfo.UAInfo{Os: "android"},
			TimeStamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
		},
		messaging.Event{
			AppID:     "app1",
			EventType: "sharelink:/d/",
			Channels:  []string{""},
			SenderID:  "a2",
			UniqueID:  "w1",
			UAInfo:    uainfo.UAInfo{Os: "ios"},
			TimeStamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
		},
		messaging.Event{
			AppID:     "app1",
			EventType: "sharelink:/d/",
			Channels:  []string{""},
			SenderID:  "",
			UniqueID:  "w1",
			UAInfo:    uainfo.UAInfo{},
			TimeStamp: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local).Unix(),
		},
	}
	for _, v := range testmessages {
		err := ds.Insert(&v)
		if err != nil {
			t.Fatal(err)
		}
	}
}

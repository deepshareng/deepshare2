package attribpush

import (
	"github.com/MISingularity/deepshare2/pkg/storage"
	"testing"
)

func TestSimpleAppToUrlSetGet(t *testing.T) {
	testAppID := "appid1"
	sau := newSimpleAppToUrl(storage.NewInMemSimpleKV())

	if u, _ := sau.GetUrl(testAppID); u != "" {
		t.Error("url should be empty before set")
	}

	if err := sau.SetUrl(testAppID, "url1"); err != nil {
		t.Fatal(err)
	}
	u, err := sau.GetUrl(testAppID)
	if err != nil {
		t.Fatal(err)
	}
	if u != "url1" {
		t.Errorf("want = %s, get = %s\n", "url1", u)
	}

	if err := sau.SetUrl(testAppID, "url2"); err != nil {
		t.Fatal(err)
	}
	u, err = sau.GetUrl(testAppID)
	if err != nil {
		t.Fatal(err)
	}
	if u != "url2" {
		t.Errorf("want = %s, get = %s\n", "url2", u)
	}

	if err := sau.SetUrl(testAppID, ""); err != nil {
		t.Fatal(err)
	}
	if u, _ := sau.GetUrl(testAppID); u != "" {
		t.Errorf("want = %s, get = %s\n", "", u)
	}
}

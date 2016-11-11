package appinfo

import (
	"testing"

	"github.com/MISingularity/deepshare2/pkg/storage"
)

func TestAppIDShortID(t *testing.T) {
	tests := []struct {
		shortID string
		appID   string
	}{
		{"shortid", "longlonglongid"},
	}
	db := storage.NewInMemSimpleKV()

	for i, tt := range tests {
		if err := SetAppIDShortIDPair(db, tt.shortID, tt.appID); err != nil {
			t.Fatalf("#%d Failed to call SetAppIDShortIDPair, err: %v", i, err)
		}
		appID, err := GetAppID(db, tt.shortID)
		if err != nil {
			t.Fatalf("#%d Failed to call GetAppID, err: %v", i, err)
		}
		if appID != tt.appID {
			t.Errorf("#%d appID got = %s, want = %s\n", i, appID, tt.appID)
		}

		shortID, err := GetShortID(db, tt.appID)
		if err != nil {
			t.Fatalf("#%d Failed to call GetShortID, err: %v", i, err)
		}
		if shortID != tt.shortID {
			t.Errorf("#%d shortID got = %s, want = %s\n", i, shortID, tt.shortID)
		}
	}

}

package attribpush

import (
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/storage"
)

type AppToUrl interface {
	SetUrl(appID string, url string) error
	GetUrl(appID string) (string, error)
}

type simpleAppToUrl struct {
	db storage.SimpleKV
}

func newSimpleAppToUrl(db storage.SimpleKV) *simpleAppToUrl {
	return &simpleAppToUrl{
		db: db,
	}
}

func (au *simpleAppToUrl) SetUrl(appID string, url string) error {
	log.Info("SetUrl:", appID, "->", url)
	if url == "" {
		if err := au.db.Delete([]byte(appID)); err != nil {
			return err
		}
		return nil
	}
	return au.db.Set([]byte(appID), []byte(url))
}

func (au *simpleAppToUrl) GetUrl(appID string) (string, error) {
	b, err := au.db.Get([]byte(appID))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

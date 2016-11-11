/*
Package deepshare implements the business logic and provides functions to support REST API.
*/
package shorturl

import (
	"net/url"
	"path"
	"time"

	"strings"

	"net/http"

	"errors"

	"github.com/MISingularity/deepshare2/deepshared/appinfo"
	"github.com/MISingularity/deepshare2/deepshared/token"
	in "github.com/MISingularity/deepshare2/pkg/instrumentation"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/storage"
)

const (
	ShortURLLifeTimeDefault = time.Duration(24) * time.Hour
)

var ErrShortSegNotFound = errors.New("shortseg not found")

// Server provides the functions to support REST API.
type UrlShortener interface {
	// GenerateShortURL generate a token, saves the rawUrl's query and then returns a short url for future access.
	// The token is unique under specific AppID.
	// The host will remain unchanged
	ToShortURL(rawUrl *url.URL, namespace string, isPermanent bool, shortURLLifeTime time.Duration, useShortID bool) (*url.URL, error)
	// GetShortURL gets the raw url given namespace and raw url.
	ToRawURL(shortUrl *url.URL, namespace string) (*url.URL, error)
}

func NewUrlShortener(cli *http.Client, skv storage.SimpleKV, specificTokenURL string) UrlShortener {
	return &server{
		skv:              skv,
		specificTokenURL: specificTokenURL,
		cli:              cli,
	}
}

func IsLegalShortFormat(url *url.URL) (isShort, isLegal bool) {
	lenPath := len(strings.Split(url.Path, "/"))
	if lenPath == 3 {
		isShort = false
		isLegal = true
	} else if lenPath == 4 {
		isShort = true
		isLegal = true
	} else {
		log.Errorf("UrlShortener; url is in wrong format; url: %s", url.String())
		isLegal = false
	}
	return
}

type server struct {
	skv              storage.SimpleKV
	specificTokenURL string
	cli              *http.Client
}

func (s *server) ToShortURL(rawUrl *url.URL, namespace string, isPermanent bool, shortURLLifeTime time.Duration, useShortID bool) (*url.URL, error) {
	v := rawUrl.RawQuery
	token, err := token.GetNewToken(s.cli, s.specificTokenURL, namespace)
	if err != nil {
		return nil, err
	}
	k := append([]byte(namespace), token...)

	start := time.Now()
	if isPermanent {
		if err := s.skv.Set(k, []byte(v)); err != nil {
			return nil, err
		}
	} else {
		if err := s.skv.SetEx(k, []byte(v), shortURLLifeTime); err != nil {
			return nil, err
		}
	}

	in.PrometheusForShorturl.StorageSaveDuration(start)
	shortUrl := new(url.URL)
	shortUrl.Scheme = rawUrl.Scheme
	shortUrl.Host = rawUrl.Host
	shortUrl.Path = path.Join(rawUrl.Path, token)

	//handle short app id
	if useShortID {
		appID := ""
		shortAppID := ""
		parts := strings.Split(rawUrl.Path, "/")
		if len(parts) > 1 {
			appID = parts[len(parts)-1]
		}
		if appID != "" {
			shortID, err := appinfo.GetShortID(s.skv, appID)
			if err == nil && shortID != "" {
				shortAppID = shortID
			}
		}
		if shortAppID != "" {
			log.Debug("ToShortURL; appID:", appID, "shortID:", shortAppID)
			shortUrl.Path = strings.Replace(shortUrl.Path, "/"+appID+"/", "/"+shortAppID+"/", -1)
		}
	}

	log.Debugf("UrlShortener; short url is generated: %s, with original raw url %s", shortUrl.String(), rawUrl.String())
	return shortUrl, nil
}

func (s *server) ToRawURL(shortUrl *url.URL, namespace string) (*url.URL, error) {
	tokenStr := path.Base(shortUrl.Path)
	start := time.Now()

	//handle short app id
	shortAppID := ""
	appID := ""
	parts := strings.Split(shortUrl.Path, "/")
	if len(parts) > 1 {
		shortAppID = parts[len(parts)-2]
	}
	if shortAppID != "" && len(shortAppID) < 16 {
		id, err := appinfo.GetAppID(s.skv, shortAppID)
		if err == nil && id != "" {
			appID = id
		}
	}
	log.Debug("ToRawURL; shortID:", shortAppID, "appID:", appID)
	if appID != "" {
		namespace = appID
	}

	k := append([]byte(namespace), tokenStr...)
	v, err := s.skv.Get(k)
	if err != nil {
		return nil, err
	}

	in.PrometheusForShorturl.StorageGetDuration(start)
	rawUrl := new(url.URL)
	rawUrl.Scheme = shortUrl.Scheme
	rawUrl.Host = shortUrl.Host
	rawUrl.Path = path.Dir(shortUrl.Path)

	if appID != "" {
		rawUrl.Path = strings.Replace(shortUrl.Path, "/"+shortAppID+"/", "/"+appID+"/", -1)
	}

	if v == nil {
		log.Errorf("UrlShortener; raw url is not found with short url %s", shortUrl.String())
		//If not found, just return the url that not contain query data.
		return rawUrl, ErrShortSegNotFound
	} else {
		rawUrl.RawQuery = string(v)
		log.Debugf("UrlShortener; raw url is found: %s, with short url %s", rawUrl.String(), shortUrl.String())
	}

	return rawUrl, nil
}

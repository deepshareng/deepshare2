package shorturl

import (
	"reflect"
	"testing"

	"net/url"
	"path"

	"net/http"
	"net/http/httptest"

	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/token"
	"github.com/MISingularity/deepshare2/pkg/storage"
)

// TestGenerateAndGetURL tests pair functioning of GenerateShortURL() and GetRawURL().
func TestGenerateAndGetURL(t *testing.T) {
	testData := []struct {
		longUrl string
	}{
		{
			"http://fds.so/u/7713337217A6E150?download_title=aaa&download_msg=bbb&redirect_url=ccc&inapp_data=y",
		},
	}
	tokenServer := httptest.NewServer(token.NewTokenTestHandler(api.TokenPrefix))
	s := NewUrlShortener(http.DefaultClient, storage.NewInMemSimpleKV(), tokenServer.URL+api.TokenPrefix)
	for i, strRawUrl := range testData {
		testRawUrl, err := url.Parse(strRawUrl.longUrl)
		if err != nil {
			t.Fatalf("#%d Parse rawUrl failed: %v, please check test cases", i, err)
		}
		shortUrl, err := s.ToShortURL(testRawUrl, "7713337217A6E150", false, ShortURLLifeTimeDefault, false)
		if err != nil {
			t.Fatalf("#%d SaveShortURL failed: %v", i, err)
		}
		rawUrl, err := s.ToRawURL(shortUrl, "7713337217A6E150")
		if err != nil {
			t.Fatalf("#%d GetShortURL failed: %v", i, err)
		}
		if !reflect.DeepEqual(testRawUrl, rawUrl) {
			t.Errorf("#%d LongURL = %#v, want = %#v", i, rawUrl, testRawUrl)
		}
	}
}

// TestGenerateShortURL tests generate url, check each possible fields
func TestGenerateShortURL(t *testing.T) {
	testData := []string{
		"http://fds.so/d/7713337217A6E150?download_title=aaa&download_msg=bbb&redirect_url=ccc&inapp_data=y",
	}
	tokenServer := httptest.NewServer(token.NewTokenTestHandler(api.TokenPrefix))
	s := NewUrlShortener(http.DefaultClient, storage.NewInMemSimpleKV(), tokenServer.URL+api.TokenPrefix)
	for i, strRawUrl := range testData {
		testRawUrl, err := url.Parse(strRawUrl)
		if err != nil {
			t.Fatalf("#%d Parse rawUrl failed: %v, please check test cases", i, err)
		}
		shortUrl, err := s.ToShortURL(testRawUrl, "7713337217A6E150", false, ShortURLLifeTimeDefault, false)
		if err != nil {
			t.Fatalf("#%d SaveShortURL failed: %v", i, err)
		}
		if shortUrl.Scheme != testRawUrl.Scheme {
			t.Errorf("#%d Scheme = %#v, want = %#v", i, shortUrl.Scheme, testRawUrl.Scheme)
		}
		if shortUrl.Host != testRawUrl.Host {
			t.Errorf("#%d Host = %#v, want = %#v", i, shortUrl.Host, testRawUrl.Host)
		}
		if path.Dir(shortUrl.Path) != testRawUrl.Path {
			t.Errorf("#%d Path = %#v, want = %#v", i, path.Dir(shortUrl.Path), testRawUrl.Path)
		}
	}
}

// TestGetRawURL tests get raw url
func TestGetNotExistRawURL(t *testing.T) {
	testData := []struct {
		shortUrl string
		rawUrl   string
		err      error
	}{
		{
			"http://fds.so/u/7713337217A6E150/aabbcc",
			"http://fds.so/u/7713337217A6E150",
			ErrShortSegNotFound,
		},
	}

	tokenServer := httptest.NewServer(token.NewTokenTestHandler(api.TokenPrefix))
	s := NewUrlShortener(http.DefaultClient, storage.NewInMemSimpleKV(), tokenServer.URL+api.TokenPrefix)
	for i, test := range testData {
		testShortUrl, err := url.Parse(test.shortUrl)
		if err != nil {
			t.Fatalf("#%d Parse rawUrl failed: %v, please check test cases", i, err)
		}
		rawUrl, err := s.ToRawURL(testShortUrl, "7713337217A6E150")
		if err != ErrShortSegNotFound {
			t.Fatalf("#%d SaveShortURL failed: %v", i, err)
		}
		if rawUrl.String() != test.rawUrl {
			t.Errorf("#%d RawUrl = %#v, want %s", i, rawUrl, test.rawUrl)
		}
	}
}

func TestURLFormat(t *testing.T) {
	testData := []struct {
		url     string
		isShort bool
		isLegal bool
	}{
		{
			"http://fds.so/u/7713337217A6E150/MA==",
			true,
			true,
		},
		{
			"http://fds.so/u/7713337217A6E150/aabbcc",
			true,
			true,
		},
		{
			"http://fds.so/u/7713337217A6E150",
			false,
			true,
		},
		{
			"http://fds.so/u/7713337217A6E150?download_title=aaa&download_msg=bbb&redirect_url=ccc&inapp_data=y",
			false,
			true,
		},
		{
			"http://fds.so/u",
			false,
			false,
		},
		{
			"http://fds.so/u/7713337217A6E150/aabbcc/aabbcc",
			false,
			false,
		},
		{
			"http://fds.so/u/7713337217A6E150/aabbcc?download_title=aaa",
			true,
			true,
		},
	}
	for i, urlData := range testData {
		testUrl, err := url.Parse(urlData.url)
		if err != nil {
			t.Fatalf("#%d Parse rawUrl failed: %v, please check test cases", i, err)
		}
		isShort, isLegal := IsLegalShortFormat(testUrl)
		if isLegal != urlData.isLegal {
			t.Errorf("URL #%d = %s the Url is legal: %t, want: %t ", i, urlData.url, isLegal, urlData.isLegal)
		}
		if !isShort && urlData.isShort {
			t.Errorf("URL #%d = %s Should be in short format", i, urlData.url)
		}
		if isShort && !urlData.isShort {
			t.Errorf("URL #%d = %s Should not be in short format", i, urlData.url)
		}
	}
}

func TestToShortURLWithLifeTime(t *testing.T) {
	tests := []struct {
		appID       string
		rawURL      string
		lifeTime    time.Duration
		emptyRawURL string
	}{
		{
			"7713337217A6E150",
			"http://fds.so/u/7713337217A6E150?download_title=aaa&download_msg=bbb&redirect_url=ccc&inapp_data=y",
			time.Duration(100) * time.Millisecond,
			"http://fds.so/u/7713337217A6E150",
		},
		{
			"7713337217A6E150",
			"http://fds.so/u/7713337217A6E150?download_title=aaa&download_msg=bbb&redirect_url=ccc&inapp_data=y",
			time.Duration(200) * time.Millisecond,
			"http://fds.so/u/7713337217A6E150",
		},
	}
	tokenServer := httptest.NewServer(token.NewTokenTestHandler(api.TokenPrefix))
	s := NewUrlShortener(http.DefaultClient, storage.NewInMemSimpleKV(), tokenServer.URL+api.TokenPrefix)

	for i, tt := range tests {
		testRawURL, err := url.Parse(tt.rawURL)
		if err != nil {
			t.Fatalf("#%d Parse rawUrl failed: %v, please check test cases", i, err)
		}
		u, err := s.ToShortURL(testRawURL, tt.appID, false, tt.lifeTime, false)
		if err != nil {
			t.Fatalf("#%d ToShortURL failed, err:%v", i, err)
		}

		// ToRawURL should get the same rawURL
		rawURL, err := s.ToRawURL(u, tt.appID)
		if err != nil {
			t.Fatalf("#%d ToRawURL failed, err:%v", i, err)
		}
		if rawURL.String() != tt.rawURL {
			t.Errorf("#%d restore rawURL failed, rawURL=%s, want=%s\n", i, rawURL.String(), tt.rawURL)
		}

		// shortURL should expire after lifetime
		time.Sleep(tt.lifeTime)
		rawURL, err = s.ToRawURL(u, tt.appID)
		if err != ErrShortSegNotFound {
			t.Fatalf("#%d ToRawURL failed, err:%v", i, err)
		}
		if rawURL.String() != tt.emptyRawURL {
			t.Errorf("#%d shortURL should expire after lifetime! rawURL=%s, want=%s\n", i, rawURL.String(), tt.emptyRawURL)
		}
	}

}

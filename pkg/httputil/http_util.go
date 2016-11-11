package httputil

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/MISingularity/deepshare2/pkg/log"
)

// allowMethod verifies that the given method is one of the allowed methods,
// and if not, it writes an error to w.  A boolean is returned indicating
// whether or not the method is allowed.
func AllowMethod(w http.ResponseWriter, m string, ms ...string) bool {
	for _, meth := range ms {
		if m == meth {
			return true
		}
	}
	w.Header().Set("Allow", strings.Join(ms, ","))
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	return false
}

func WriteHTTPError(w http.ResponseWriter, err HTTPError) {
	http.Error(w, err.Error(), err.StatusCode)
}

func RequestLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: use a logger pkg and change this to debug level
		log.Infof("[%s] %s, remote: %s", r.Method, r.RequestURI, r.RemoteAddr)
		handler.ServeHTTP(w, r)
	})
}

func AppendPath(urlOri string, newPathSeg string) (newUrl string, err error) {
	urlNew, err := url.Parse(urlOri)
	if err != nil {
		return "", err
	}
	pathUrl := urlNew.Path
	newPath := path.Join(pathUrl, newPathSeg)
	urlNew.Path = newPath
	return urlNew.String(), nil
}

func GetNewClient() *http.Client {
	//TODO: we need to discuss how to define this transport
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{RootCAs: nil},
		DisableCompression: true,
		DisableKeepAlives:  false,
	}
	client := &http.Client{Transport: tr}
	return client
}

func ParseClientIP(r *http.Request) string {
	xForwardFor := r.Header.Get("X-Forwarded-For")
	ip := ""
	if xForwardFor != "" {
		ip = strings.Split(xForwardFor, ", ")[0]
		log.Debugf("HttpUtil; xForwardFor exsit; ip: %s; X-Forwarded-For: %s", ip, xForwardFor)
	} else {
		var err error
		ip, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Errorf("HttpUtil; SplitHostPort error; RemoteAddr: %s", r.RemoteAddr)
			return ""
		} else {
			log.Debugf("HttpUtil; xForwardFor does not exsit; ip: %s", r.RemoteAddr)
		}
	}

	return ip
}

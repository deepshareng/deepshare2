package token

import (
	"net/http"
	"time"

	"strings"

	"encoding/json"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/pkg/httputil"
	in "github.com/MISingularity/deepshare2/pkg/instrumentation"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/tokenutil"
)

type tokenHandler struct {
	tg       tokenutil.TokenGenerator
	endpoint string
}

func newTokenHandler(workerID, dataCenterID int64, endpoint string) http.Handler {
	tg, err := tokenutil.NewSnowflakeTokenGenerator(workerID, dataCenterID)
	if err != nil {
		log.Fatal("Failed to new snowflake token generator, err:", err)
	}
	return &tokenHandler{
		tg:       tg,
		endpoint: endpoint,
	}
}

func NewTokenTestHandler(endpoint string) http.Handler {
	return newTokenHandler(0, 0, endpoint)
}

func AddHandler(mux *http.ServeMux, workerID, dataCenterID int64, endpoint string) {
	mux.Handle(endpoint, newTokenHandler(workerID, dataCenterID, endpoint))
}

func (th *tokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !httputil.AllowMethod(w, r.Method, "GET") {
		return
	}

	start := time.Now()
	switch r.Method {
	case "GET":
		defer in.PrometheusForToken.HTTPGetDuration(start)
		fields := strings.Split(r.URL.Path[len(th.endpoint):], "/")
		if len(fields) != 1 {
			log.Error("Token; Invalid path:", r.URL.Path)
			httputil.WriteHTTPError(w, api.ErrPathNotFound)
			return
		}

		namespace := fields[0]
		if namespace == "" {
			log.Error("Token; Invalid path:", r.URL.Path, ", namespace should not be empty")
			httputil.WriteHTTPError(w, api.ErrPathNotFound)
			return
		}
		token, err := th.tg.Generate(namespace)
		if err != nil {
			log.Error("Token; generate token failed, err:", err)
			httputil.WriteHTTPError(w, api.ErrInternalServer)
			return
		}
		resp := TokenResponse{
			Token: token,
		}
		en := json.NewEncoder(w)
		if err := en.Encode(resp); err != nil {
			log.Error("Token; encode json err:", err)
			httputil.WriteHTTPError(w, api.ErrInternalServer)
		}
	}
}

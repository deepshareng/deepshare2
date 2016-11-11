package attribpush

import (
	"encoding/json"
	"net/http"

	"github.com/MISingularity/deepshare2/pkg/httputil"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/storage"
)

type appUrlHandler struct {
	endpoint string
	au       AppToUrl
}

func AddHandler(mux *http.ServeMux, endpoint string, db storage.SimpleKV) AppToUrl {
	au := newSimpleAppToUrl(db)
	mux.Handle(endpoint, newAppUrlHandler(au, endpoint))
	return au
}
func newAppUrlHandler(au AppToUrl, endpoint string) *appUrlHandler {
	return &appUrlHandler{
		endpoint: endpoint,
		au:       au,
	}
}

type AppUrl struct {
	Url string `json:"url"`
}

func (h *appUrlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info("simpleAppToUrl is called", r.URL)
	if !httputil.AllowMethod(w, r.Method, "PUT", "GET") {
		return
	}
	appID := r.URL.Path[len(h.endpoint):]
	if appID == "" {
		log.Error("need appid")
		return
	}
	switch r.Method {
	case "PUT":
		decoder := json.NewDecoder(r.Body)
		req := new(AppUrl)
		if err := decoder.Decode(req); err != nil {
			log.Error("PUT app url error:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.au.SetUrl(appID, req.Url); err != nil {
			log.Error("Set url error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "GET":
		url, err := h.au.GetUrl(appID)
		if err != nil {
			log.Error("Get url error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp := AppUrl{Url: url}
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(&resp); err != nil {
			log.Error("GET app url error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

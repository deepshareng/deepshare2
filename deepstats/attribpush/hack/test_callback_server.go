package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/MISingularity/deepshare2/deepstats/attribpush"
	"github.com/MISingularity/deepshare2/pkg/log"
)

const (
	serverAddr = "127.0.0.1:8082"
)

// http://127.0.0.1:8082/callback
func main() {
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("/callback is called")
		defer r.Body.Close()
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		var req []attribpush.AttributionPushInfo
		if err := json.Unmarshal(b, &req); err != nil {
			panic(err)
		}
		log.Infof("%#v\n", req)
	})
	log.Debug("Listen on", serverAddr)
	http.ListenAndServe(serverAddr, nil)
}

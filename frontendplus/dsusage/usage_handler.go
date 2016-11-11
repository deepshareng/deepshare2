package dsusage

import (
	"fmt"
	"net/http"

	"strings"

	"encoding/json"

	"strconv"

	"time"

	usage "github.com/MISingularity/deepshare2/deepstats/dsusage"
	"github.com/MISingularity/deepshare2/pkg/httputil"
	in "github.com/MISingularity/deepshare2/pkg/instrumentation"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/storage"
)

type usageHandler struct {
	endpoint string
	db       storage.SimpleKV
}

func AddHandler(mux *http.ServeMux, endpoint string, db storage.SimpleKV) {
	mux.Handle(endpoint, newUsageHandler(db, endpoint))
}

func newUsageHandler(db storage.SimpleKV, endpoint string) *usageHandler {
	return &usageHandler{
		endpoint: endpoint,
		db:       db,
	}
}

func (uh *usageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("usageHandler is called, request:", r.Method, r.URL.Path)
	if !httputil.AllowMethod(w, r.Method, "GET", "DELETE") {
		return
	}

	start := time.Now()
	var appID, senderID string
	seg := r.URL.Path[len(uh.endpoint):]
	if parts := strings.Split(seg, "/"); len(parts) >= 2 {
		appID = parts[0]
		senderID = parts[1]
	}

	if appID == "" || senderID == "" {
		log.Errorf("usageHandler, appID and senderID should not be empty, appID: %s, senderID: %s\n", appID, senderID)
		http.Error(w, "appID and senderID should not be empty", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		defer in.PrometheusForDSUsage.HTTPGetDuration(start)
		installs, opens := uh.GetUsage(appID, senderID)
		log.Debug("installs:", installs, "opens:", opens, "appID:", appID, "senderID:", senderID)
		resp := &ResponseGetNewUsageObj{
			Installs: installs,
			Opens:    opens,
		}
		en := json.NewEncoder(w)
		if err := en.Encode(resp); err != nil {
			log.Error("usageHandler, json encode error:", err)
			http.Error(w, "json encode error", http.StatusInternalServerError)
			return
		}
	case "DELETE":
		defer in.PrometheusForDSUsage.HTTPDeleteDuration(start)
		uh.ClearUsage(appID, senderID)
	}
}

func (uh *usageHandler) GetUsage(appID string, senderID string) (installs int, opens int) {
	k := []byte(fmt.Sprintf(usage.RedisKeyUsageFmt, appID, senderID))
	installs = uh.getIntFromDB(k, usage.RedisHKeyInstall)
	opens = uh.getIntFromDB(k, usage.RedisHKeyOpen)
	return installs, opens
}

func (uh *usageHandler) ClearUsage(appID string, senderID string) {
	k := []byte(fmt.Sprintf(usage.RedisKeyUsageFmt, appID, senderID))
	var installs, opens int
	installs = uh.getIntFromDB(k, usage.RedisHKeyInstall)
	opens = uh.getIntFromDB(k, usage.RedisHKeyOpen)
	//don't waste to do the next operations if usages are 0
	if installs == 0 && opens == 0 {
		return
	}
	start := time.Now()
	if err := uh.db.HDel(k, usage.RedisHKeyInstall); err != nil {
		log.Error(err)
		panic(err)
	}
	in.PrometheusForDSUsage.StorageDeleteDuration(start)
	start = time.Now()
	if err := uh.db.HDel(k, usage.RedisHKeyOpen); err != nil {
		log.Error(err)
		panic(err)
	}
	in.PrometheusForDSUsage.StorageDeleteDuration(start)

	start = time.Now()
	if err := uh.db.HIncrBy(k, usage.RedisHKeyInstall+"_total", installs); err != nil {
		log.Error(err)
	}
	in.PrometheusForDSUsage.StorageIncDuration(start)

	start = time.Now()
	if err := uh.db.HIncrBy(k, usage.RedisHKeyOpen+"_total", opens); err != nil {
		log.Error(err)
	}
	in.PrometheusForDSUsage.StorageIncDuration(start)
}

func (uh *usageHandler) getIntFromDB(k []byte, hk string) int {
	start := time.Now()
	if b, err := uh.db.HGet(k, hk); err != nil {
		log.Error(err)
		panic(err)
	} else {
		in.PrometheusForDSUsage.StorageGetDuration(start)
		if n, err := strconv.ParseInt(string(b), 10, 64); err != nil {
			log.Debug("[Warn]GetUsage ParseInt err:", err, "k:", string(k), "hk:", hk, "v:", string(b))
		} else {
			return int(n)
		}
	}
	return 0
}

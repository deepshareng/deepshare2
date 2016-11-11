package appinfo

import (
	"encoding/json"
	"net/http"

	"github.com/MISingularity/deepshare2/pkg/httputil"
	"github.com/MISingularity/deepshare2/pkg/log"
)

func GetAppInfoByUrl(client *http.Client, appID, appInfoUrl string) (appInfo *AppInfo, err error) {
	appInfoUrlStr, err := httputil.AppendPath(appInfoUrl, appID)
	if err != nil {
		log.Errorf("shareLinkHandler; appInfoUrl %s is constructed by us, should not in wrong format: %v", appInfoUrl, err)
		panic(err)
	}

	log.Debugf("Share link; Request App Info URL string is %s", appInfoUrlStr)
	req, err := http.NewRequest("GET", appInfoUrlStr, nil)
	if err != nil {
		log.Error("Share link; Setup New Request to App Info service failed", err)
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Share link; Request to App Info service failed", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Error("Share link; app info response code is not OK:", resp.StatusCode)
		return nil, nil
	}

	decoder := json.NewDecoder(resp.Body)
	appInfoObj := AppInfo{}
	err = decoder.Decode(&appInfoObj)
	if err != nil {
		log.Error("Share link; Decode response from App Info service failed", err)
		return nil, err
	}
	return &appInfoObj, nil
}

package token

import (
	"errors"
	"net/http"
	"strconv"

	"encoding/json"

	"github.com/MISingularity/deepshare2/pkg/httputil"
	"github.com/MISingularity/deepshare2/pkg/log"
)

// GetNewToken is a helper function to request a new token from token service
func GetNewToken(cli *http.Client, specificTokenURL string, namespace string) (string, error) {
	tokenUrl, err := httputil.AppendPath(specificTokenURL, namespace)
	if err != nil {
		log.Error("GetNewToken; failed to append path, err:", err)
		return "", err
	}
	req, err := http.NewRequest("GET", tokenUrl, nil)
	if err != nil {
		log.Error("GetNewToken; failed to construct request, err:", err)
		return "", err
	}
	resp, err := cli.Do(req)
	if err != nil {
		log.Error("GetNewToken; failed to get response from token server, err:", err)
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		errInfo := "GetNewToken; got wrong http response code:" + strconv.Itoa(resp.StatusCode)
		log.Error(errInfo)
		return "", errors.New(errInfo)
	}
	defer resp.Body.Close()
	de := json.NewDecoder(resp.Body)
	token := &TokenResponse{}
	if err := de.Decode(token); err != nil {
		log.Error("GetNewToken; decode from http response err:", err)
		return "", err
	}
	return token.Token, nil
}

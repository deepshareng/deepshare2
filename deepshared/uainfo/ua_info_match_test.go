package uainfo

import (
	"io/ioutil"
	"testing"

	"github.com/MISingularity/deepshare2/pkg/path"
	"gopkg.in/yaml.v2"
)

func TestUainfoMatching(t *testing.T) {
	curdir, err := path.Getcurdir()
	if err != nil {
		t.Errorf("Get current path failed! Err Msg=%v\n", err)
		return
	}
	data, err := ioutil.ReadFile(curdir + "/match_test.yaml")
	if err != nil {
		t.Errorf("Read file failed! Err Msg=%v\n", err)
		return
	}
	m := make(map[string][]map[string]string)
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		t.Errorf("Unmarshal yaml file failed! Err Msg=%v\n", err)
		return
	}
	for i := range m["testcases"] {
		d1 := ExtractUAInfoFromUAString("", m["testcases"][i]["app"])
		d2 := ExtractUAInfoFromUAString("", m["testcases"][i]["browser"])

		if d1.OsVersion != d2.OsVersion || d1.Brand != d2.Brand {
			t.Errorf("#%d : App/browser UAInfo matching error when parse \n %v, \n %v \n get browser os=%s, browser brand=%s,\n get app os=%s, app brand=%s",
				i,
				m["testcases"][i]["browser"],
				m["testcases"][i]["app"],
				d2.Os, d2.Brand,
				d1.Os, d1.Brand)
		}

	}
}

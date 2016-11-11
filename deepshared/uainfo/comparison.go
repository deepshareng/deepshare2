package uainfo

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/MISingularity/deepshare2/pkg/path"
	"gopkg.in/yaml.v2"
)

var uainfoFather map[string]string

func NewBrandFather() error {
	curdir, _ := path.Getcurdir()
	m := make(map[string][]map[string]string)
	data, err := ioutil.ReadFile(curdir + "/comparison_brand_fathers.yaml")
	if err != nil {
		return fmt.Errorf("Read file failed! Err Msg=%v\n", err)
	}
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("Unmarshal yaml file failed! Err Msg=%v\n", err)
	}
	uainfoFather = make(map[string]string)
	for _, v := range m["brand_fathers"] {
		uainfoFather[v["son"]] = v["father"]
	}
	return nil
}

func uainfoFathersComparison(fbrand, cbrand string) bool {
	for fbrand != "" {
		if strings.Contains(fbrand, cbrand) {
			return true
		}
		fbrand = uainfoFather[fbrand]
	}
	return false
}

func BrandComparison(brand1, brand2 string) (bool, error) {
	if len(uainfoFather) == 0 {
		err := NewBrandFather()
		if err != nil {
			return false, err
		}
	}
	if uainfoFathersComparison(brand1, brand2) || uainfoFathersComparison(brand2, brand1) {
		return true, nil
	}
	return false, nil
}

func OsComparison(d1, d2 *UAInfo) bool {
	if d1.Os != d2.Os {
		return false
	}
	if !strings.HasPrefix(d1.OsVersion, d2.OsVersion) && !strings.HasPrefix(d2.OsVersion, d1.OsVersion) {
		return false
	}
	return true
}

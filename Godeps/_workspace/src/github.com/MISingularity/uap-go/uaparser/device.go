package uaparser

import (
	"regexp"
	"strings"
)

type Device struct {
	Family string
	Brand  string
	Model  string
}

type DevicePattern struct {
	Regexp            *regexp.Regexp
	Regex             string
	RegexFlag         string
	BrandReplacement  string
	DeviceReplacement string
	ModelReplacement  string
}

var RegexPattern = ""

func (dvcPattern *DevicePattern) Match(line string, dvc *Device) {
	matches := dvcPattern.Regexp.FindStringSubmatch(line)
	if len(matches) == 0 {
		return
	}
	// fmt.Println(dvcPattern)
	RegexPattern = dvcPattern.Regex
	groupCount := dvcPattern.Regexp.NumSubexp()

	if len(dvcPattern.DeviceReplacement) > 0 {
		dvc.Family = allMatchesReplacement(dvcPattern.DeviceReplacement, matches)
	} else if groupCount >= 1 {
		dvc.Family = matches[1]
	}
	dvc.Family = strings.TrimSpace(dvc.Family)
	if strings.Contains(dvc.Family, "SonyEricsson") {
		dvc.Family = strings.Replace(dvc.Family, "SonyEricsson", "Sony", -1)
	}

	if len(dvcPattern.BrandReplacement) > 0 {
		dvc.Brand = allMatchesReplacement(dvcPattern.BrandReplacement, matches)
	} else if groupCount >= 2 {
		dvc.Brand = matches[1]
	}
	dvc.Brand = strings.TrimSpace(dvc.Brand)
	if dvc.Brand == "SonyEricsson" {
		dvc.Brand = "Sony"
	}

	if len(dvcPattern.ModelReplacement) > 0 {
		dvc.Model = allMatchesReplacement(dvcPattern.ModelReplacement, matches)
	} else if groupCount >= 2 {
		dvc.Model = matches[1]
	}
	dvc.Model = strings.TrimSpace(dvc.Model)
}

func (dvc *Device) ToString() string {
	if strings.HasPrefix(dvc.Family, "HW-") && strings.ToUpper(dvc.Brand) == "HUAWEI" {
		dvc.Family = dvc.Family[3:]
	}
	return strings.ToUpper(dvc.Family)
}

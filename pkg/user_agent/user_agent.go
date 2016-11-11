//
// Copyright (C) 2015 MISingularity.io
// author: Flora Wang <flora.wang@misingularity.io>
// (copied from "github.com/mssola/user_agent" and then patched based on our own demands)

// Package user_agent implements an HTTP User Agent string parser. It defines
// the type UserAgent that contains all the information from the parsed string.
// It also implements the Parse function and getters for all the relevant
// information that has been extracted from a parsed User Agent string.
package user_agent

import (
	"regexp"
	"strconv"
	"strings"
)

// A section contains the name of the product, its version and
// an optional comment.
type section struct {
	name    string
	version string
	comment []string
}

// The UserAgent struct contains all the info that can be extracted
// from the User-Agent string.
type UserAgent struct {
	ua           string
	mozilla      string
	platform     string
	os           string
	localization string
	browser      Browser
	bot          bool
	mobile       bool
	undecided    bool
	wechat       string
	qq           string
}

// Read from the given string until the given delimiter or the
// end of the string have been reached.
//
// The first argument is the user agent string being parsed. The second
// argument is a reference pointing to the current index of the user agent
// string. The delimiter argument specifies which character is the delimiter
// and the cat argument determines whether nested '(' should be ignored or not.
//
// Returns an array of bytes containing what has been read.
func readUntil(ua string, index *int, delimiter byte, cat bool) []byte {
	var buffer []byte

	i := *index
	catalan := 0
	for ; i < len(ua); i = i + 1 {
		if ua[i] == delimiter {
			if catalan == 0 {
				*index = i + 1
				return buffer
			}
			catalan--
		} else if cat && ua[i] == '(' {
			catalan++
		}
		buffer = append(buffer, ua[i])
	}
	*index = i + 1
	return buffer
}

// Parse the given product, that is, just a name or a string
// formatted as Name/Version.
//
// It returns two strings. The first string is the name of the product and the
// second string contains the version of the product.
func parseProduct(product []byte) (string, string) {
	prod := strings.SplitN(string(product), "/", 2)
	if len(prod) == 2 {
		return prod[0], prod[1]
	}
	return string(product), ""
}

// Parse a section. A section is typically formatted as follows
// "Name/Version (comment)". Both, the comment and the version are optional.
//
// The first argument is the user agent string being parsed. The second
// argument is a reference pointing to the current index of the user agent
// string.
//
// Returns a section containing the information that we could extract
// from the last parsed section.
func parseSection(ua string, index *int) (s section) {
	buffer := readUntil(ua, index, ' ', false)

	s.name, s.version = parseProduct(buffer)
	if *index < len(ua) && ua[*index] == '(' {
		*index++
		buffer = readUntil(ua, index, ')', true)
		s.comment = strings.Split(string(buffer), "; ")
		*index++
	}
	return s
}

// Initialize the parser.
func (p *UserAgent) initialize() {
	p.ua = ""
	p.mozilla = ""
	p.platform = ""
	p.os = ""
	p.localization = ""
	p.browser.Engine = ""
	p.browser.EngineVersion = ""
	p.browser.Name = ""
	p.browser.Version = ""
	p.bot = false
	p.mobile = false
	p.undecided = false
	p.wechat = ""
}

// Parse the given User-Agent string and get the resulting UserAgent object.
//
// Returns an UserAgent object that has been initialized after parsing
// the given User-Agent string.
func New(ua string) *UserAgent {
	o := &UserAgent{}
	o.Parse(ua)
	return o
}

// Parse the given User-Agent string. After calling this function, the
// receiver will be setted up with all the information that we've extracted.
func (p *UserAgent) Parse(ua string) {
	var sections []section

	p.initialize()
	p.ua = ua
	for index, limit := 0, len(ua); index < limit; {
		s := parseSection(ua, &index)
		if !p.mobile && s.name == "Mobile" {
			p.mobile = true
		}
		if s.name == "MicroMessenger" {
			p.wechat = s.version
		}
		if s.name == "QQ" {
			p.qq = s.version
		}
		//patch for android uas that os is not in the first section
		if strings.ToLower(s.name) == "android" {
			p.os = s.name + " " + s.version
		}

		sections = append(sections, s)
	}

	if len(sections) > 0 {
		p.mozilla = sections[0].version

		p.detectBrowser(sections)
		if p.os == "" {
			p.detectOS(sections[0])
		}

		if p.undecided {
			p.checkBot(sections)
		}
	}
}

// Returns the mozilla version (it's how the User Agent string begins:
// "Mozilla/5.0 ...", unless we're dealing with Opera, of course).
func (p *UserAgent) Mozilla() string {
	return p.mozilla
}

// Returns true if it's a bot, false otherwise.
func (p *UserAgent) Bot() bool {
	return p.bot
}

// Returns true if it's a mobile device, false otherwise.
func (p *UserAgent) Mobile() bool {
	return p.mobile
}

func (p *UserAgent) IsApple() bool {
	os, _ := p.OsNameVersion()
	return os == "ios"
}

func (p *UserAgent) IsAndroid() bool {
	os, _ := p.OsNameVersion()
	return os == "android"
}

func (p *UserAgent) IsWechat() bool {
	return p.wechat != ""
}

func (p *UserAgent) IsQQ() bool {
	return p.qq != ""
}

func (p *UserAgent) IsWeibo() bool {
	ua := strings.ToLower(p.ua)
	return strings.Contains(ua, "weibo")
}

func (p *UserAgent) IsFirefox() bool {
	ua := strings.ToLower(p.ua)
	return strings.Contains(ua, "firefox")
}

// OsNameVersion returns unified os & version parsed from ua string
// os could be "ios" or "android"
// version is unified to the format of a.b.c
// TODO should extract os & osversion when p *UserAgent is created and save them as fields for
//      performance
// TODO split this function to p.Os() & p.OsVersion()
func (p *UserAgent) OsNameVersion() (os, version string) {
	if p.bot {
		return "", ""
	}
	osinfo := ""
	if p.os != "" {
		osinfo = p.os + " " + p.platform
	} else {
		osinfo = p.ua
	}

	osinfo = strings.ToLower(osinfo)
	if strings.Contains(osinfo, "iphone") || strings.Contains(osinfo, "ipad") || strings.Contains(osinfo, "ipod") || strings.Contains(osinfo, "ios") {
		os = "ios"
		words := strings.Split(osinfo, " ")
		for _, word := range words {
			if len(word) > 0 && word[0] >= '0' && word[0] <= '9' {
				//in some browser, version is in format of a_b_c, we should unify to a.b.c
				// TODO (sean.wu) Do we need to do this for android?
				version = strings.Replace(word, "_", ".", -1)
				return os, version
			}
		}
		return os, ""
	} else if strings.Contains(osinfo, "android") {
		os = "android"
		words := strings.Split(osinfo, " ")
		for _, word := range words {
			word = strings.TrimSuffix(word, ";")
			word = strings.TrimSuffix(word, ",")
			if len(word) > 0 && word[0] >= '0' && word[0] <= '9' {
				return os, word
			}
		}
		return os, ""
	}
	return "", ""
}

// Brand returns the brand of the mobile device.
// return brand only if it is within the supportedBrands.
// It returns "-", if it cannot find brand.
// For android the ua pattern is usually like:
// 		<Browser Name>/<Browser version> (<platform>; <OS>; [<locale>;] [otherinfo;] <brand> <model> [<build>]) <Rendering engine> / <Rendering engine version> [otherinfo]
func (p *UserAgent) Brand() string {
	s := p.ua
	idx := strings.LastIndex(s, ";")
	if idx < 0 {
		return "-"
	}
	segment := s[(idx + 1):]
	brand := parseBrand(segment)

	//to fix the following UA could not find brand:
	// Mozilla/5.0 (Linux; Android 6.0; Nexus 6 Build/MRA58N; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/45.0.2454.95 Mobile Safari/537.36 MicroMessenger/6.3.7.51_rbb7fa12.660 NetType/WIFI Language/zh_CN
	if brand == "-" {
		segment := s[:idx]
		parts := strings.Split(segment, ";")
		segment = parts[len(parts)-1]
		brand = parseBrand(segment)
	}
	return brand
}

func parseBrand(segment string) string {
	segment = strings.TrimPrefix(segment, " ")
	parts := strings.Split(segment, " ")

	brand := parts[0]
	brand = strings.ToLower(brand)
	supportedBrands := []string{"coolpad", "hm", "htc", "huawei", "lenovo", "mi", "mi-one", "mx4", "nexus", "tcl", "vivo", "zte"}
	pattern := strings.Join(supportedBrands, "$|^")
	pattern = "^" + pattern + "$"
	reg := regexp.MustCompile(pattern)
	if !reg.MatchString(brand) {
		brand = "-"
	}
	return brand
}

// ChromeMajorVersion returns the major version code of chrome
// returns 0 if it's not a chrome UA
func (p *UserAgent) ChromeMajorVersion() int {
	ua := strings.ToLower(p.ua)

	version := ""
	for _, s := range []string{"chrome", "crios", "crmo", "chromium"} {
		idx := strings.Index(ua, s)
		if idx != -1 {
			pattern := s
			remaining := ua[idx:]
			i := strings.Index(remaining, " ")
			chrome := remaining
			if i != -1 {
				chrome = remaining[:i]
			}
			version = strings.TrimPrefix(chrome, pattern)
			version = strings.TrimPrefix(version, "/")
			break
		}
	}

	return majorVersion(version)
}

// TODO: this method is convoluted.
// IosMajorVersion returns the major version code of ios.
// returns 0 when it's not ios UA
func (p *UserAgent) IosMajorVersion() int {
	os, osv := p.OsNameVersion()
	if os != "ios" {
		return 0
	}
	return majorVersion(osv)
}

// Take a version string a.b.c return int(a).
func majorVersion(osv string) int {
	if osv == "" {
		return 0
	}

	idx := strings.Index(osv, ".")
	majorStr := osv
	if idx != -1 {
		majorStr = osv[:idx]
	}
	n, err := strconv.Atoi(majorStr)
	if err == nil && n > 0 {
		return n
	}
	return 0
}

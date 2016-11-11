package uaparser

import (
	"regexp"
	"strconv"
	"strings"
)

type Channel struct {
	Params map[string]bool
}

type ChannelPattern struct {
	Model             string
	Regexp            *regexp.Regexp
	Regex             string
	ExactMatch        string
	UserAgentFamily   string
	LeastMajorVersion string
}

func (channelPattern *ChannelPattern) leastMajorVersionJudgement(cli *Client, line string) bool {
	if channelPattern.LeastMajorVersion == "" {
		return true
	}

	cm, err := strconv.Atoi(channelPattern.LeastMajorVersion)
	if err != nil {
		return false
	}
	um, err := strconv.Atoi(cli.UserAgent.Major)
	if err != nil {
		return false
	}
	if um >= cm {
		return true
	}
	return false
}

func (channelPattern *ChannelPattern) userAgentJudgement(cli *Client, line string) bool {
	if channelPattern.UserAgentFamily != "" && strings.Contains(cli.UserAgent.Family, channelPattern.UserAgentFamily) {
		return true
	}
	return false
}

func (channelPattern *ChannelPattern) Match(cli *Client, line string) bool {
	if channelPattern.ExactMatch != "" && strings.Contains(line, channelPattern.ExactMatch) {
		return true
	}
	if channelPattern.Regex != "" {
		matches := channelPattern.Regexp.FindStringSubmatch(line)
		if len(matches) != 0 {
			return true
		}
	}
	if channelPattern.userAgentJudgement(cli, line) && channelPattern.leastMajorVersionJudgement(cli, line) {
		return true
	}
	return false
}

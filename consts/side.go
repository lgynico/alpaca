package consts

import "regexp"

var (
	clientReg = regexp.MustCompile("[cC]")
	serverReg = regexp.MustCompile("[sS]")
)

type Side func(string) bool

var (
	SideClient Side = func(s string) bool { return len(s) == 0 || clientReg.MatchString(s) }
	SideServer Side = func(s string) bool { return len(s) == 0 || serverReg.MatchString(s) }
)

const (
	OutputServer = "server"
	OutputClient = "client"
)

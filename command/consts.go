package command

import "regexp"

const (
	FlagInput  = "input"
	FlagOutput = "output"
	FlagServer = "server"
	FlagClient = "client"
)

var (
	regexpGolang = regexp.MustCompile(`[Gg][Oo]([Ll][Aa][Nn][Gg])*`)
	regexpCSharp = regexp.MustCompile(`[Cc](#|[Ss][Hh][Aa][Rr][Pp])`)
	//regexpJava   = regexp.MustCompile(`[Jj][Aa][Vv][Aa]`)
)

type CodeType int32

const (
	CodeUnknown CodeType = iota
	CodeGolang
	CodeCSharp
)

func codeType(str string) CodeType {
	if len(str) == 0 {
		return CodeUnknown
	}

	if regexpGolang.MatchString(str) {
		return CodeGolang
	}

	if regexpCSharp.MatchString(str) {
		return CodeCSharp
	}

	return CodeUnknown
}

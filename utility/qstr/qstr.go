package qstr

import "github.com/gogf/gf/v2/text/gstr"

func ReplaceN(origin string) string {
	return gstr.ReplaceByMap(origin, map[string]string{
		"     ": "",
		"\n":    " ",
	})
}

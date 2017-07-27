package main

import (
	"strings"
	"testing"
)

func mapping(ss []string) map[string]struct{} {
	ret := map[string]struct{}{}

	for _, s := range ss {
		ret[s] = struct{}{}
	}

	return ret
}

func Test_findKeys(t *testing.T) {

	buf := strings.NewReader(`colin,*,aaa
# xxxxxx
colin,yyy|zzz|xxx,bbb
colin,*.yyy,ccc
nobody,*,xxx
`)

	users, err := readConfig(buf)
	if err != nil {
		t.Error(err)
	}

	keys := mapping(users.findKeys("colin", "abc.com"))
	if _, finded := keys["aaa"]; !finded || len(keys) != 1 {
		t.Error("find error", keys)
	}

	keys = mapping(users.findKeys("colin", "abc.yyy"))
	if len(keys) != 2 {
		t.Error("find error", keys)
	}
	if _, finded := keys["aaa"]; !finded {
		t.Error("find error", keys)
	}
	if _, finded := keys["ccc"]; !finded {
		t.Error("find error", keys)
	}

}

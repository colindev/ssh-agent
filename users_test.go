package main

import (
	"strings"
	"testing"
)

func Test_findKeys(t *testing.T) {

	buf := strings.NewReader(`colin,*,aaa
# xxxxxx
colin,yyy,bbb
colin,*.yyy,ccc
nobody,*,xxx
`)

	users, err := readConfig(buf)
	if err != nil {
		t.Error(err)
	}

	var keys []string

	keys = users.findKeys("colin", "abc.com")
	if len(keys) != 1 || keys[0] != "aaa" {
		t.Error("find error", keys)
	}

	keys = users.findKeys("colin", "abc.yyy")
	if len(keys) != 2 || keys[0] != "aaa" || keys[1] != "ccc" {
		t.Error("find error", keys)
	}

}

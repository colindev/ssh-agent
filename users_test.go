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

func Test_readConfig(t *testing.T) {
	buf := strings.NewReader(`colin,2,3
1,2,3
1,2,3
`)

	_, err := readConfig(buf)
	if err != nil {
		t.Error(err)
	}
}

func TestUsers_findKeys(t *testing.T) {

	buf := strings.NewReader(`
colin,*,aaa
# xxxxxx
colin,yyy|zzz|xxx,bbb
colin,*.yyy,ccc
nobody,*,xxx
`)

	users, err := readConfig(buf)
	if err != nil {
		t.Error(err)
		t.Skip()
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

func TestUsers_findUsers(t *testing.T) {

	buf := strings.NewReader(`
user_1,*.project-1,aaa
user_2,*.project-1,bbb
user_3,*.project-2,ccc
`)

	users, err := readConfig(buf)
	if err != nil {
		t.Error(err)
		t.Skip()
	}

	ret := mapping(users.findUsers("project-1"))
	if len(ret) != 2 {
		t.Error("find error", ret)
	}
	if _, finded := ret["user_1"]; !finded {
		t.Error("find error", ret)
	}
	if _, finded := ret["user_2"]; !finded {
		t.Error("find error", ret)
	}
}

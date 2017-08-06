package main

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"regexp"
	"strings"
	"sync"
)

func readConfig(r io.Reader) (*Users, error) {

	var line int
	users := newUsers()

	conf := csv.NewReader(r)
	conf.Comment = '#'
	for {
		line++
		row, err := conf.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("line=", line, err)
		}

		if err := users.add(row); err != nil {
			return nil, err
		}
	}

	return users, nil
}

// Users ...
type Users struct {
	sync.RWMutex
	keys map[string]map[string]map[string]*regexp.Regexp
}

func newUsers() *Users {
	return &Users{
		keys: make(map[string]map[string]map[string]*regexp.Regexp),
	}
}

func (u *Users) swap(users *Users) {
	u.Lock()
	defer u.Unlock()
	u.keys = users.keys
}

func (u *Users) add(row []string) error {
	if len(row) != 3 {
		return errors.New("config format must be [user,key] in one record")
	}
	u.Lock()
	defer u.Unlock()

	user := row[0]
	key := row[2]
	tags := strings.Split(row[1], "|")
	if len(tags) == 0 {
		tags = []string{"*"}
	}

	_, exists := u.keys[user]
	if !exists {
		u.keys[user] = map[string]map[string]*regexp.Regexp{}
	}
	_, exists = u.keys[user][key]
	if !exists {
		u.keys[user][key] = map[string]*regexp.Regexp{}
	}

	for _, tag := range tags {

		if tag == "" || tag == "-" {
			tag = "*"
		}

		if _, exists := u.keys[user][key][tag]; !exists {
			tag = strings.Replace(tag, ".", "\\.", -1)
			tag = strings.Replace(tag, "*", ".*", -1)
			re := regexp.MustCompile("^" + tag + "$")

			u.keys[user][key][tag] = re
		}
	}

	return nil
}

// only for GCP
func (u *Users) findUsers(project string) []string {

	ret := []string{}

	u.RLock()
	defer u.RUnlock()

	for user, mps := range u.keys {
		for _, fps := range mps {
			for fp := range fps {
				if strings.HasSuffix(fp, "."+project) {
					ret = append(ret, user)
				}
			}
		}
	}

	return ret
}

func (u *Users) findKeys(username, hostname string) []string {

	ret := []string{}

	u.RLock()
	defer u.RUnlock()

	keys := u.keys[username]
	hn := []byte(hostname)
	for key, res := range keys {
		for _, re := range res {
			if re.Match(hn) {
				ret = append(ret, key)
			}
		}
	}

	return ret
}

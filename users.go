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

	users := newUsers()

	conf := csv.NewReader(r)
	conf.Comment = '#'
	for {
		row, err := conf.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
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
	rowLen := len(row)
	if rowLen < 2 {
		return errors.New("config format must be [user,key] in one record")
	}
	u.Lock()
	defer u.Unlock()

	user := row[0]
	key := row[rowLen-1]
	tags := row[1 : rowLen-1]
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

		if _, exists := u.keys[user][key][tag]; !exists {
			tag = strings.Replace(tag, ".", "\\.", -1)
			tag = strings.Replace(tag, "*", ".*", -1)
			re := regexp.MustCompile("^" + tag + "$")

			u.keys[user][key][tag] = re
		}
	}

	return nil
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

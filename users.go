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
	keys map[string]map[*regexp.Regexp]string
}

func newUsers() *Users {
	return &Users{
		keys: make(map[string]map[*regexp.Regexp]string),
	}
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
		u.keys[user] = map[*regexp.Regexp]string{}
	}

	for _, tag := range tags {
		tag = strings.Replace(tag, ".", "\\.", -1)
		tag = strings.Replace(tag, "*", ".*", -1)
		re := regexp.MustCompile("^" + tag + "$")

		u.keys[user][re] = key
	}

	return nil
}

func (u *Users) findKeys(username, hostname string) []string {

	ret := []string{}

	u.RLock()
	defer u.RUnlock()

	keys := u.keys[username]
	hn := []byte(hostname)
	for re, key := range keys {
		if re.Match(hn) {
			ret = append(ret, key)
		}
	}

	return ret
}

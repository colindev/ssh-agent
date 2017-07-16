package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/julienschmidt/httprouter"
)

var env struct {
	Addr       string
	ConfigFile string
}

func init() {

	flag.StringVar(&env.Addr, "addr", ":6666", "http listen address")
	flag.StringVar(&env.ConfigFile, "conf", "/etc/ssh-agent-server.conf", "config file for authorized keys (use csv format)")
}

func main() {

	flag.Parse()

	users, err := readConfig(env.ConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()
	router.GET("/:user/keys", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		w.Write([]byte(users.findKey(p.ByName("user"), r.URL.Query().Get("fingerprint"))))
	})

	log.Println(http.ListenAndServe(":6666", router))
}

func readConfig(filename string) (*Users, error) {

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	users := newUsers()

	conf := csv.NewReader(f)
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

func (u *Users) findKey(username, fingerprint string) string {

	u.RLock()
	defer u.RUnlock()

	keys := u.keys[username]
	fp := []byte(fingerprint)
	for re, key := range keys {
		if re.Match(fp) {
			return key
		}
	}

	return ""
}

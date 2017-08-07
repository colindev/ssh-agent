package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/julienschmidt/httprouter"
)

const (
	// LogDebug specific more detail for process
	LogDebug = 1 << iota
)

var (
	version string
	env     struct {
		Addr       string
		ConfigFile string
	}
	logLevel    int
	showVersion bool
)

func init() {

	flag.IntVar(&logLevel, "log", 0, "log level (0-1)")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.StringVar(&env.Addr, "addr", ":6666", "http listen address")
	flag.StringVar(&env.ConfigFile, "conf", "/etc/ssh-agent-server/config", "config file for authorized keys (use csv format)")
}

func main() {

	flag.Parse()

	fmt.Println("Version:", version)
	if showVersion {
		os.Exit(0)
	}

	var (
		users *Users
	)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)
	go func() {

		for {
			f, err := os.Open(env.ConfigFile)
			if err != nil {
				log.Fatal(err)
			}
			u, err := readConfig(f)
			f.Close()
			if err != nil {
				log.Fatal(err)
			}
			if users != nil {
				users.swap(u)
			} else {
				users = u
			}
			<-c
			log.Println("reload:", env.ConfigFile)
		}
	}()

	scriptsDir := path.Dir(env.ConfigFile)
	log.Printf("scripts dir [%s]\n", scriptsDir)

	router := httprouter.New()
	router.GET("/users/:user/keys", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		user := p.ByName("user")
		fp := r.URL.Query().Get("fingerprint")
		keys := users.findKeys(user, fp)
		log.Printf("search [%s] keys: %s find(%d)\n", user, fp, len(keys))
		buf := bytes.NewBuffer(nil)
		for _, key := range keys {
			// for debug
			buf.WriteString(key + "\n")
		}
		w.Write(buf.Bytes())
		if logLevel&LogDebug > 0 {
			log.Println(buf.String())
		}
	})

	router.HandlerFunc(http.MethodGet, "/users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, user := range users.findUsers(r.URL.Query().Get("project")) {
			w.Write([]byte(user + "\n"))
		}
	}))

	router.HandlerFunc(http.MethodGet, "/scripts/installer.sh", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadFile(scriptsDir + "/scripts/installer.sh")
		if err != nil {
			log.Println("read script error:", err)
			http.NotFound(w, r)
			return
		}

		selfLink := r.Host
		resBody := strings.Replace(string(b), "{{selfLink}}", selfLink, -1)

		w.Write([]byte(resBody))
	}))

	log.Println(http.ListenAndServe(env.Addr, router))
}

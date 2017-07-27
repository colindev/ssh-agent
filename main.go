package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/julienschmidt/httprouter"
)

var (
	version string
	env     struct {
		Addr       string
		ConfigFile string
	}
	showVersion bool
)

func init() {

	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.StringVar(&env.Addr, "addr", ":6666", "http listen address")
	flag.StringVar(&env.ConfigFile, "conf", "/etc/ssh-agent-server.conf", "config file for authorized keys (use csv format)")
}

func main() {

	flag.Parse()

	fmt.Println(version)
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

	router := httprouter.New()
	router.GET("/:user/keys", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
		log.Println(buf.String())
	})

	log.Println(http.ListenAndServe(env.Addr, router))
}

package main

import (
	"flag"
	"log"
	"net/http"
	"os"

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

	f, err := os.Open(env.ConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	users, err := readConfig(f)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()
	router.GET("/:user/keys", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		for _, key := range users.findKeys(p.ByName("user"), r.URL.Query().Get("fingerprint")) {
			w.Write([]byte(key + "\n"))
		}
	})

	log.Println(http.ListenAndServe(env.Addr, router))
}

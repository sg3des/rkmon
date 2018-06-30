package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/badoux/checkmail"
	"github.com/gorilla/schema"
	"github.com/sg3des/bytetree"
)

var configfile = flag.String("config", "rkn.conf", "path to configuration file")
var emails = flag.Bool("emails", true, "enable/disable sending alert emails")

func main() {
	fmt.Println("RKN monitor, version: 1.0.1-20180630")

	flag.Parse()

	conf, err := NewConfig(*configfile)
	if err != nil {
		log.Fatal("failed parse configuration file", *configfile, err)
	}

	if conf.Debug {
		log.SetFlags(log.Ltime | log.Lshortfile)
	} else {
		log.SetFlags(log.Ltime)
	}

	rep, err := NewRepository(conf.GitRepo)
	if err != nil {
		log.Fatal("faield initialize repository,", err)
	}

	api := &API{
		rep:         rep,
		reg:         NewRegister(),
		repofile:    conf.GitRepoFilename,
		debug:       conf.Debug,
		subscribers: bytetree.NewTree(),
		smtp:        conf.SMTP,
	}

	for _, subs := range conf.Subscribers {
		api.AddSubscriber(subs)
	}

	go api.Daemon(conf.GitInterval)

	err = api.WebServer(conf.HTTPaddr)
	if err != nil {
		log.Fatal(err)
	}
}

type API struct {
	rep  *Repository
	reg  *Register
	smtp SMTP

	repofile   string
	CommitTime time.Time

	subscribers *bytetree.Tree

	debug bool
}

//
// Daemon reload list of blocked ip addresses
//   check subscrides addrs and send emails if this blocked
//

func (api *API) Daemon(interval time.Duration) {
	for {
		err := api.Update()
		if err != nil {
			log.Println(err)
		}

		time.Sleep(interval)
	}
}

func (api *API) Update() error {
	commitTime, err := api.rep.Get()
	if err != nil {
		return err
	}

	if commitTime.Equal(api.CommitTime) {
		return nil
	}
	if commitTime.Before(api.CommitTime) {
		return errors.New("local commit newer than remote commit, it discourages")
	}

	start := time.Now()
	if api.debug {
		log.Println("update storage...")
	}

	f, err := api.rep.OpenFile(api.repofile)
	if err != nil {
		return err
	}

	api.reg.Load(f)

	if api.debug {
		log.Println("storage updated, it took", time.Now().Sub(start))
	}

	api.CommitTime = commitTime

	api.CheckSubscribers()

	return nil
}

//
// WEB SERVER
//

func (api *API) WebServer(addr string) error {
	http.HandleFunc("/", api.Index)
	http.HandleFunc("/subscribe/", api.Subscribe)
	http.HandleFunc("/unsubscribe/", api.Unsubscribe)
	http.HandleFunc("/assets/", api.Assets)

	fmt.Println("start web-server on", addr)

	return http.ListenAndServe(addr, nil)
}

//Index urlpath /?ip=127.0.0.1
func (api *API) Index(w http.ResponseWriter, r *http.Request) {
	if ip := r.URL.Query().Get("ip"); len(ip) != 0 {
		api.ipQuery(ip, w, r)
		return
	}

	api.Template(w)
}

func (api *API) ipQuery(sip string, w http.ResponseWriter, r *http.Request) {
	ip := net.ParseIP(sip)
	if ip.IsUnspecified() {
		err := "invalid request, failed parse ip " + sip
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	if api.reg.LookupIP(ip) {
		http.Error(w, fmt.Sprintf("IP address %s is blocked", ip), http.StatusUnavailableForLegalReasons)
		return
	}

	fmt.Fprintf(w, "IP address %s is not blocked", ip)
}

//Subscribe GET request /subscribe/?email=email@email.com&ip=127.0.0.1&ip=...
func (api *API) Subscribe(w http.ResponseWriter, r *http.Request) {
	var req Subscriber

	err := schema.NewDecoder().Decode(&req, r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := checkmail.ValidateFormat(req.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	api.AddSubscriber(&req)
	fmt.Fprintf(w, "successfull subscribe %s for %s ip address", req.Email, req.IP)
}

//Unsubscribe GET request /unsubscribe/?email=email@email.com
func (api *API) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		err := "invalid request, required variable 'email' not found"
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	if err := checkmail.ValidateFormat(email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	api.RemoveSubscriber(email)
	fmt.Fprintf(w, "successfull unsubscribe %s", email)
}

//Assets serve static files stored to go code how bind-data
func (api *API) Assets(w http.ResponseWriter, r *http.Request) {
	assetname := r.URL.Path[1:]

	//lookup assets file
	fi, err := AssetInfo(assetname)
	if err != nil {
		http.Error(w, fmt.Sprintf("File %s not found", r.URL.Path), 404)
		return
	}

	//check modified date
	modSince, err := time.Parse(time.RFC1123, r.Header.Get("If-Modified-Since"))
	if err == nil && fi.ModTime().Before(modSince) {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	//restore assets from bind data
	data, err := Asset(assetname)
	if err != nil {
		http.Error(w, fmt.Sprintf("File %s not found", r.URL.Path), http.StatusNotFound)
		return
	}

	//server file
	http.ServeContent(w, r, r.URL.Path, time.Now(), bytes.NewReader(data))
}

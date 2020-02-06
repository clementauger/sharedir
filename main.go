package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/99designs/basicauth-go"
	"github.com/clementauger/tor-prebuilt/embedded"
	"github.com/cretz/bine/tor"
	"github.com/cretz/bine/torutil"
	tued25519 "github.com/cretz/bine/torutil/ed25519"
	"github.com/gorilla/handlers"
)

func main() {

	var c string
	var dir string
	var user string
	var pwd string
	var pkpath string
	var realm string
	flag.StringVar(&c, "c", "config.json", "the configuration file path to load")
	flag.StringVar(&dir, "d", "", "directory path to serve")
	flag.StringVar(&user, "u", "", "username")
	flag.StringVar(&pwd, "p", "", "password")
	flag.StringVar(&pkpath, "pk", "onion.pk", "ed25519 pem encoded privatekey file path")
	flag.StringVar(&realm, "realm", "", "realm to show in login box")
	flag.Parse()

	config, err := loadConfig(c)
	if err != nil {
		log.Fatal(err)
	}
	if dir != "" {
		config.Path = dir
	}
	if realm != "" {
		config.Realm = realm
	}
	if user != "" {
		config.Users[user] = []string{pwd}
	}

	templatePath := "assets/dir.tpl"

	var tpl TplExecer
	if build == "dev" {
		tpl, err = fileTemplate(templatePath)
	} else {
		tpl, err = assetTemplate(templatePath)
	}
	if err != nil {
		log.Fatal(err)
	}

	fs := FileServer(http.Dir(config.Path), DirList(tpl))
	auth := basicauth.New(config.Realm, config.Users)
	h := handlers.LoggingHandler(os.Stdout, auth(fs))

	var server serverListener
	if build == "dev" {
		server = &http.Server{
			Addr:    ":9090",
			Handler: h,
		}
	} else {
		server = &torServer{
			PrivateKey: pkpath,
			Handler:    h,
		}
	}

	errc := make(chan error)
	go func() {
		errc <- server.ListenAndServe()
	}()

	sc := make(chan os.Signal)
	signal.Notify(sc)
	select {
	case err := <-errc:
		log.Fatal(err)
	case s := <-sc:
		log.Printf("got signal %v\n", s)
	}
}

func getOrCreatePK(fpath string) (ed25519.PrivateKey, error) {
	var privateKey ed25519.PrivateKey
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		_, privateKey, err = ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		x509Encoded, err := x509.MarshalPKCS8PrivateKey(privateKey)
		if err != nil {
			return nil, err
		}
		pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "ED25519 PRIVATE KEY", Bytes: x509Encoded})
		ioutil.WriteFile(fpath, pemEncoded, os.ModePerm)
	} else {
		d, _ := ioutil.ReadFile(fpath)
		block, _ := pem.Decode(d)
		x509Encoded := block.Bytes
		tPk, err := x509.ParsePKCS8PrivateKey(x509Encoded)
		if err != nil {
			return nil, err
		}
		if x, ok := tPk.(ed25519.PrivateKey); ok {
			privateKey = x
		} else {
			return nil, fmt.Errorf("invalid key type %T wanted ed25519.PrivateKey", tPk)
		}
	}
	return privateKey, nil
}

type serverListener interface {
	ListenAndServe() error
}

type torServer struct {
	Handler http.Handler
	// PrivateKey path to a pem encoded ed25519 private key
	PrivateKey string
}

func onion(pk ed25519.PrivateKey) string {
	return torutil.OnionServiceIDFromV3PublicKey(tued25519.PublicKey([]byte(pk.Public().(ed25519.PublicKey))))
}

func (ts *torServer) ListenAndServe() error {

	pk, err := getOrCreatePK(ts.PrivateKey)
	if err != nil {
		return err
	}

	d, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}

	// Start tor with default config (can set start conf's DebugWriter to os.Stdout for debug logs)
	fmt.Println("Starting and registering onion service, please wait a couple of minutes...")
	t, err := tor.Start(nil, &tor.StartConf{TempDataDirBase: d, ProcessCreator: embedded.NewCreator(), NoHush: true})
	if err != nil {
		return fmt.Errorf("unable to start Tor: %v", err)
	}
	defer t.Close()

	// Wait at most a few minutes to publish the service
	listenCtx, listenCancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer listenCancel()
	// Create a v3 onion service to listen on any port but show as 80
	onion, err := t.Listen(listenCtx, &tor.ListenConf{Key: pk, Version3: true, RemotePorts: []int{80}})
	if err != nil {
		return fmt.Errorf("unable to create onion service: %v", err)
	}
	defer onion.Close()

	fmt.Printf("server listening at http://%v.onion\n", onion.ID)

	return http.Serve(onion, ts.Handler)
}

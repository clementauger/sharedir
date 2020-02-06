# sharedir

zero configuration WAN sharing folder. You basically run `./sharedir` and it begins serving the current directory over TOR. It automatically creates and persists an ed25519 pem encoded private key for you, or use the one you will configure.

You can use a json configuration file or the command line arguments to setup options,

```go
$ go run . -h
Usage of /tmp/go-build281710412/b001/exe/sharedir:
  -c string
    	the configuration file path to load (default "config.json")
  -d string
    	directory path to serve
  -p string
    	password
  -pk string
    	ed25519 pem encoded privatekey file path (default "onion.pk")
  -realm string
    	realm to show in login box
  -u string
    	username
```

Check the `config.json` and `config.go` files to learn more about the configuration options.

upon starting it shows

```sh
févr. 06 12:16:00 Host-001 sharedir[141917]: Starting and registering onion service, please wait a couple of minutes...
févr. 06 12:16:00 Host-001 sharedir[141917]: 2020/02/06 12:16:00 /tmp/803903500/tor [--DataDirectory /tmp/619647320/data-dir-921963223 --CookieAuthentication 1 --DisableNetwork 1 --SocksPort auto -f /tmp/619647320/data-dir-921963223/torrc-586636362 --ControlPort auto --ControlPortWriteToFile /tmp/619647320/data-dir-921963223/control-port-873716001]
...
févr. 06 12:16:05 Host-001 sharedir[141917]: Feb 06 12:16:05.000 [notice] Bootstrapped 100% (done): Done
févr. 06 12:16:08 Host-001 sharedir[141917]: server listening at http://xxxx.onion
```

you can browse to the given onion address (http://xxxx.onion)

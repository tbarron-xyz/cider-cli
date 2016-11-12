# cider-cli
Simple command line interface for [cider](http://github.com/tbarron-xyz/cider). Connects to an already-running instance of cider.

#Usage
`./cider-cli -port 1234 -url google.com` will attempt to connect to a cider instance running at `google.com:1234`. The default port is `6969` (the default port for cider) and the default url is localhost.

`cider cli
Connecting... Connected.
> SET key value
< {
    "response": null,
    "status": "success"
}
> GET key
< {
    "response": "value",
    "status": "success"
}`

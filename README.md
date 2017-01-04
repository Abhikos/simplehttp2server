`simplehttp2server` serves the current directory on an HTTP/2.0 capable server.
This server is for development purposes only.

# Download

## Binaries 
`simplehttp2server` is `go get`-able:

```
$ go get github.com/GoogleChrome/simplehttp2server
```

Precompiled binaries can be found in the [release section](https://github.com/GoogleChrome/simplehttp2server/releases).

## Brew
You can also install `simplehttp2server` using brew if you are on macOS:

```
$ brew tap GoogleChrome/simplehttp2server https://github.com/GoogleChrome/simplehttp2server
$ brew install simplehttp2server
```

## Docker
If you have Docker set up, you can serve the current directory via `simplehttp2server` using the following command:

```
$ docker run -p 5000:5000 -v $PWD:/data surma/simplehttp2server
```

# Push Manifest

`simplehttp2server` supports the [push manifest](https://www.npmjs.com/package/http2-push-manifest).
All requests will be looked up in a file named `push.json`. If there is a key
for the request path, all resources under that key will be pushed.

Example `push.json`:

```JS
{
  "/": {
    "/css/app.css": {
      "type": "style",
      "weight": 1
    },
    // ...
  },
  "/page.html": {
    "/css/page.css": {
      "type": "style",
      "weight": 1
    },
    // ...
  }
}
```

Support for `weight` and `type` is not implemented yet. Pushes cannot trigger additional pushes.

# TLS Certificate

Since HTTP/2 requires TLS, `simplehttp2server` checks if `cert.pem` and
`key.pem` are present. If not, a self-signed certificate will be generated.

# Delays

`simplehttp2server` can add artificial delays to responses to emulate processing
time. The command line flags `-mindelay` and `-maxdelay` allow you to delay
responses with a random delay form the interval `[minDelay, maxDelay]` in milliseconds.

If a request has a `delay` query parameter (like `GET /index.html?delay=4000`),
that delay will take precedence.

# Other features

* Support for serving Single Page Applications (SPAs) using the `-spa` flag
* Support for throttling network throughput *per request* using the `-throttle` flag

# License

Apache 2.

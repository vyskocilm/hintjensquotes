# Mastodon quoting bot

[Pieter Hintjens](http://hintjens.com) edition.

```html
<a rel="me" href="https://botsin.space/@hintjensquotes">Mastodon</a>
```

## Run

Run daily through Github actions

## Contributing

Use [C4.1](http://rfc.zeromq.org/spec:22) process, like zeromq does. Fell free
to propose more quotes or more powerful posting code.

## Developer run

Use provided docker-compose.yml to start the local Mastodon instance for a testing.

 * user: user@example.com
 * password: bitnami1

```sh
docker-compose up -d
```

In a `http://localhost/settings` click on a Development and create a new
application. You need `write:statuses` scope. Copy the generated access token
and export it as `MASTO_ACCESS_TOKEN` variable.

## Troubleshooting

To record the HTTP requests and responses, use `github.com/motemen/go-loghttp`

```go
	logtr := loghttp.DefaultTransport
	logtr.LogRequest = func(req *http.Request) {
		log.Printf("--> %s %s", req.Method, req.URL)
		for key, values := range req.Header {
			log.Printf("    %s: %s", key, values)
		}
		var body bytes.Buffer
		_, err := io.Copy(&body, req.Body)
		if err != nil {
			log.Fatalf("io.Copy: %s", err)
		}
		req.Body.Close()
		req.Body = io.NopCloser(&body)
		log.Printf("    %s\n", body.String())
	}
	logtr.LogResponse = func(resp *http.Response) {
		ctx := resp.Request.Context()
		if start, ok := ctx.Value(loghttp.ContextKeyRequestStart).(time.Time); ok {
			log.Printf("<-- %d %s (%s)", resp.StatusCode, resp.Request.URL, roundtime.Duration(time.Now().Sub(start), 2))
		} else {
			log.Printf("<-- %d %s", resp.StatusCode, resp.Request.URL)
		}
		for key, values := range resp.Header {
			log.Printf("    %s: %s", key, values)
		}
		var body bytes.Buffer
		_, err := io.Copy(&body, resp.Body)
		if err != nil {
			log.Fatalf("io.Copy: %s", err)
		}
		resp.Body.Close()
		resp.Body = io.NopCloser(&body)
		log.Printf("    %s\n", body.String())
	}
	http.DefaultTransport = logtr
```



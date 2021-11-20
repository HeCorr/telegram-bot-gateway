# Telegram Bot Gateway
A simple but fast webhook reverse proxy written in Go that allows you to run multiple Telegram bots in the same HTTP port.

## Usage

(Pre-built binaries coming soon)

### Building from source
1. Make sure you have [Go](https://golang.org/doc/install#download) installed (at least v1.16)
2. Clone the repo
3. Run `go build -o gateway.exe ./src`

### Running it
1. Create the `bots.yaml` file in the same folder as the executable (.yml also works):
```yaml
bots:
  - name: Bot One    # the name is optional. if omitted one will be generated.
    endpoint: /bot1
    path: http://localhost:9800/bot
  - name: Bot Two
    endpoint: /bot2
    path: http://localhost:9801/bot
```
2. Run it: 
`./gateway.exe -c <certFile> -k <keyFile>`

### Available arguments
```
-f string
      Use the specified .yaml file (default "bots.yaml")
-l string
      Listen address (default "localhost:8443")
-c string
      Certificate file for HTTPS (required)
-k string
      Private key file for HTTPS (required)
-s bool
      Strict mode - blocks requests not coming from Telegram (default false)
```

Note: Don't enable Strict mode when running this behind a proxy and try disabling it if you're having issues receiving updates.

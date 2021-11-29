# Telegram Bot Gateway
A simple but fast webhook reverse proxy written in Go that allows you to run multiple Telegram bots in the same HTTP port.

## Usage

### Running from the pre-built binaries
- Download the [latest](https://github.com/HeCorr/telegram-bot-gateway/releases/latest) binary version for your system and skip to the **Running it** section below.

### Building from source
1. Make sure you have [Go](https://golang.org/doc/install#download) installed (at least v1.16)
2. Clone the repo
3. Run `go build -o gateway.exe ./src`

### Running it
1. Generate a config file (`bots.yaml`) by running `./gateway.exe -i`. The file will look like this:
```yaml
config:
  # Private key file
  key: bot.key
  # Certificate file
  cert: bot.crt
  # Listen address (ip:port)
  listen: localhost:8443
  # Strict mode, ignore requests not coming from Telegram
  strict: false

bots:
    # Bot name for reference (optional)
  - name: Bot One
    # Gateway endpoint for receiving updates on
    endpoint: /bot1
    # Path that requests will be forwarded to (the bot's webhook URL)
    path: https://localhost:3000/bot
  - name: Bot Two
    endpoint: /bot2
    path: https://localhost:3001/bot
```
2. Run it: `./gateway.exe`

### Available arguments
```
-f string
      Use the specified .yaml file (default "bots.yaml")
-l string
      Listen address (default "localhost:8443")
-c string
      Certificate file for HTTPS
-k string
      Private key file for HTTPS
-i
      Initialize (create) bots.yaml file
-s
      Strict mode - blocks requests not coming from Telegram
```

Note: Running this behind a proxy or a NAT with Strict mode is not recommended. Try disabling it if you're having issues receiving updates.

### TODO
- [ ] Improve logging
- [x] Add more comments to the code
- [x] Allow specifying options in the `bots.yaml` file
- [x] Improve HTTP client (set proper timeouts, ignore certs, etc)
- [x] Implement `-i` (init) flag for generating default `bots.yaml` file
- [ ] Improve non-Telegram IP blocking middleware
- [ ] Support JSON config files
- [ ] Auto reload..?

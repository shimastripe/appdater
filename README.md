# Appdater

AppStore update checker in Go

- App Store
- Google Play
- Kindle Store

## Feature

- Single executable binary (Embedded config file using [statik](https://github.com/rakyll/statik))
- Multiple application register
- Slack notification

## Config

See [appdater/config.toml.sample](https://github.com/shimastripe/appdater/blob/master/assets/config.toml.sample)

## Run

```
# make assets/config.toml
$ statik -src assets
$ go run cmd/appdater/main.go
```

## Make executable binary

```
# make assets/config.toml
$ statik -src assets
$ go build -o appdater cmd/appdater/main.go
```
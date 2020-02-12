package appdater

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/ashwanthkumar/slack-go-webhook"
	_ "github.com/shimastripe/appdater/statik"

	"github.com/BurntSushi/toml"
	"github.com/rakyll/statik/fs"
)

type CLI struct {
	OutStream, ErrStream io.Writer
}

func (c *CLI) Run(args []string) int {
	var conf Config

	statikFs, err := fs.New()
	if err != nil {
		fmt.Fprintf(c.ErrStream, "Statik cannot initialize. %v", err)
		return 1
	}

	fp, err := statikFs.Open("/config.toml")
	if err != nil {
		fmt.Fprintf(c.ErrStream, "File open error. %v", err)
		return 1
	}

	if _, err := toml.DecodeReader(fp, &conf); err != nil {
		fmt.Fprintf(c.ErrStream, "Cannot decode toml file. %v", err)
		return 1
	}

	sleep := time.Duration(conf.SleepTime*60) * time.Second

	for {

		for i, app := range conf.Android {
			latestVersion, err := ScrapeLatestVersion(app)
			if err != nil {
				fmt.Fprintf(c.ErrStream, "Cannot scrape a version: %v. %v", app.CreateAppPageUrl(), err)
				continue
			}

			if err := Validate(latestVersion); err != nil {
				message := fmt.Sprintf("Semver is missing. Perhaps DOM is changed: %v. %v\n", app.CreateAppPageUrl(), err)
				fmt.Fprint(c.ErrStream, message)
				continue
			}

			if app.LastVersion != "" && app.LastVersion != latestVersion {
				log.Printf("Update app! %v\n", latestVersion)

				if err := slack.Send(app.WebhookUrl, "", CreatePayload(app.Name, latestVersion, app.CreateAppPageUrl(), app.Emoji)); err != nil {
					fmt.Fprintf(c.ErrStream, "Cannot send slack: %v\n", err)
				}

			} else {
				log.Printf("No update!\n")
			}

			// range accessing is shallow copy
			conf.Android[i].LastVersion = latestVersion
		}

		for i, app := range conf.Ios {
			latestVersion, err := ScrapeLatestVersion(app)
			if err != nil {
				fmt.Fprintf(c.ErrStream, "Cannot scrape a version: %v. %v", app.CreateAppPageUrl(), err)
				continue
			}

			if err := Validate(latestVersion); err != nil {
				message := fmt.Sprintf("Semver is missing. Perhaps DOM is changed: %v. %v\n", app.CreateAppPageUrl(), err)
				fmt.Fprint(c.ErrStream, message)
				continue
			}

			if app.LastVersion != "" && app.LastVersion != latestVersion {
				log.Printf("Update app! %v\n", latestVersion)

				if err := slack.Send(app.WebhookUrl, "", CreatePayload(app.Name, latestVersion, app.CreateAppPageUrl(), app.Emoji)); err != nil {
					fmt.Fprintf(c.ErrStream, "Cannot send slack: %v\n", err)
				}

			} else {
				log.Printf("No update!\n")
			}

			// range accessing is shallow copy
			conf.Ios[i].LastVersion = latestVersion
		}

		for i, app := range conf.Kindle {
			latestVersion, err := ScrapeLatestVersion(app)
			if err != nil {
				fmt.Fprintf(c.ErrStream, "Cannot scrape a version: %v. %v", app.CreateAppPageUrl(), err)
				continue
			}

			if err := Validate(latestVersion); err != nil {
				message := fmt.Sprintf("Semver is missing. Perhaps DOM is changed: %v. %v\n", app.CreateAppPageUrl(), err)
				fmt.Fprint(c.ErrStream, message)
				continue
			}

			if app.LastVersion != "" && app.LastVersion != latestVersion {
				log.Printf("Update app! %v\n", latestVersion)

				if err := slack.Send(app.WebhookUrl, "", CreatePayload(app.Name, latestVersion, app.CreateAppPageUrl(), app.Emoji)); err != nil {
					fmt.Fprintf(c.ErrStream, "Cannot send slack: %v\n", err)
				}

			} else {
				log.Printf("No update!\n")
			}

			// range accessing is shallow copy
			conf.Kindle[i].LastVersion = latestVersion
		}

		time.Sleep(sleep)
	}
}

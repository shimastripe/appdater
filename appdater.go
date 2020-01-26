package appdater

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/blang/semver"
)

func Validate(version string) error {
	_, err := semver.Parse(version)
	return err
}

func CreatePayload(name string, version string, website string, emoji string) slack.Payload {
	att := slack.Attachment{}
	att.AddField(slack.Field{Title: "ストアの" + name + "が更新されました", Value: version})
	att.AddAction(slack.Action{Type: "button", Text: "Store Page", Url: website, Style: "primary"})

	return slack.Payload{
		Username:    "Store Checker",
		IconEmoji:   emoji,
		Attachments: []slack.Attachment{att},
	}
}

func CreateErrorPayload(message string, emoji string) slack.Payload {
	return slack.Payload{
		Username:  "Store Checker",
		IconEmoji: emoji,
		Text:      message,
	}
}

type App interface {
	CreateAppPageUrl() string
	GetQuery() string
	CleansingDomValue(value string) string
}

func ScrapeLatestVersion(app App) (string, error) {
	doc, err := goquery.NewDocument(app.CreateAppPageUrl())
	if err != nil {
		return "", err
	}

	domValue := app.CleansingDomValue(doc.Find(app.GetQuery()).First().Text())
	log.Print(domValue)
	return domValue, nil
}

type Android struct {
	Name        string `toml:"name"`
	Package     string `toml:"package"`
	WebhookUrl  string `toml:"webhook_url"`
	Emoji       string `toml:"emoji"`
	LastVersion string
}

func (a Android) CreateAppPageUrl() string {
	BASE_URL := "https://play.google.com/store/apps/details?id="
	googlePlayURL := BASE_URL + a.Package
	log.Print(googlePlayURL)
	return googlePlayURL
}

func (a Android) GetQuery() string {
	return "#fcxH9b > div.WpDbMd > c-wiz > div > div.ZfcPIb > div > div.JNury.Ekdcne > div > c-wiz:nth-child(4) > div.W4P4ne > div.JHTxhe.IQ1z0d > div > div:nth-child(4) > span > div > span"
}

func (a Android) CleansingDomValue(value string) string {
	return value
}

type Ios struct {
	Name        string `toml:"name"`
	Country     string `toml:"country"`
	AppID       string `toml:"app_id"`
	WebhookUrl  string `toml:"webhook_url"`
	Emoji       string `toml:"emoji"`
	LastVersion string
}

func (i Ios) CreateAppPageUrl() string {
	BASE_URL := "https://itunes.apple.com/{{country}}/app/{{appId}}"
	replaceCountryURL := strings.Replace(BASE_URL, "{{country}}", i.Country, 1)
	appStoreURL := strings.Replace(replaceCountryURL, "{{appId}}", i.AppID, 1)
	log.Print(appStoreURL)
	return appStoreURL
}

func (i Ios) GetQuery() string {
	return "p.whats-new__latest__version"
}

func (i Ios) CleansingDomValue(value string) string {
	return strings.Replace(value, "バージョン ", "", -1)
}

type Kindle struct {
	Name        string `toml:"name"`
	Asin        string `toml:"asin"`
	WebhookUrl  string `toml:"webhook_url"`
	Emoji       string `toml:"emoji"`
	LastVersion string
}

func (k Kindle) CreateAppPageUrl() string {
	BASE_URL := "https://www.amazon.co.jp/gp/product/"
	kindleStoreURL := BASE_URL + k.Asin
	log.Print(kindleStoreURL)
	return kindleStoreURL
}

func (k Kindle) GetQuery() string {
	return "#mas-technical-details div.a-spacing-none:nth-child(2)"
}

func (k Kindle) CleansingDomValue(value string) string {
	return strings.Replace(value, "バージョン: ", "", -1)
}

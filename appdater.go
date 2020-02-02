package appdater

import (
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/blang/semver"
)

var UA_LIST = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X)",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.106 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:24.0) Gecko/20100101 Firefox/24.0",
	"Mozilla/5.0 (X11; U; Linux x86_64; en-US; rv:1.9.1.1) Gecko/20090716 Linux Mint/7 (Gloria) Firefox/3.5.1",
	"Mozilla/5.0 (X11; U; Linux i686; fr; rv:1.9.0.10) Gecko/2009042708 Fedora/3.0.10-1.fc10 Firefox/3.0.10",
	"Mozilla/5.0 (X11; U; Linux x86_64; de; rv:1.9.0.11) Gecko/2009070611 Gentoo Firefox/3.0.11",
	"Mozilla/5.0 (X11; Arch Linux; Linux x86_64; rv:55.0) Gecko/20100101 Firefox/55.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3864.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:62.0) Gecko/20100101 Firefox/62.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:67.0) Gecko/20100101 Firefox/67.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:68.0) Gecko/20100101 Firefox/68.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:61.0) Gecko/20100101 Firefox/61.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Safari/605.1.15",
}

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
	req, err := http.NewRequest("GET", app.CreateAppPageUrl(), nil)
	if err != nil {
		return "", err
	}

	rand.Seed(time.Now().UnixNano())
	req.Header.Add("User-Agent", UA_LIST[rand.Intn(len(UA_LIST))])

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
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

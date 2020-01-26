package appdater

type Config struct {
	Ios       []Ios     `toml:"ios"`
	Android   []Android `toml:"android"`
	Kindle    []Kindle  `toml:"kindle"`
	SleepTime int       `toml:"sleeptime"`
}

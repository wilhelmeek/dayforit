package dayforit

import (
	"fmt"
	"math"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"

	owm "github.com/briandowns/openweathermap"
)

func Check(w http.ResponseWriter, r *http.Request) {
	err := check()
	if err != nil {
		http.Error(
			w,
			errors.Wrap(err, "running check").Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func check() error {
	wtr, err := owm.NewCurrent("C", "EN", os.Getenv("OWM_KEY"))
	if err != nil {
		return errors.Wrap(err, "getting owm new current")
	}

	wtr.CurrentByZip(2017, "AU")
	if err != nil {
		return errors.Wrap(err, "getting current by zip")
	}

	windSpeedKPH := int(math.Trunc(wtr.Wind.Speed * 60 * 60 / 1000))
	windDirection := directionFromDegrees(wtr.Wind.Deg)
	primaryWeather := wtr.Weather[0]

	var title, subtitle string
	switch true {
	case primaryWeather.ID == 800:
		title = "It is an absolute day for it! :sunny:"
		subtitle = "The weather today is *perfect*. Good day to get on the bikes!"
	case primaryWeather.ID > 800:
		title = "No reason not to cycle. :cloud:"
		subtitle = "The weather today is perfectly cycleable. Get on the bikes!"
	case primaryWeather.ID >= 300 && primaryWeather.ID < 400:
		title = "Worth a shot! :rain_cloud:"
		subtitle = "It's a bit drizzly today, but gotta risk it for the biscuit."
	default:
		title = "Nope :no_good:"
		subtitle = "Not today."
	}

	err = slack.PostWebhook(
		os.Getenv("SLACK_WEBHOOK_URL"),
		&slack.WebhookMessage{
			Text: title,
			Blocks: &slack.Blocks{
				BlockSet: []slack.Block{
					slack.NewHeaderBlock(
						slack.NewTextBlockObject(
							slack.PlainTextType,
							title,
							false,
							false,
						),
					),
					slack.NewSectionBlock(
						slack.NewTextBlockObject(
							slack.MarkdownType,
							subtitle,
							false,
							false,
						),
						nil,
						nil,
					),
					slack.NewDividerBlock(),
					slack.NewSectionBlock(
						nil,
						[]*slack.TextBlockObject{
							slack.NewTextBlockObject(
								slack.MarkdownType,
								fmt.Sprintf("*Currently*: %v??c", wtr.Main.Temp),
								false,
								false,
							),
							slack.NewTextBlockObject(
								slack.MarkdownType,
								fmt.Sprintf(
									"*Range*: %v??c - %v??c",
									wtr.Main.TempMin,
									wtr.Main.TempMax,
								),
								false,
								false,
							),
							slack.NewTextBlockObject(
								slack.MarkdownType,
								fmt.Sprintf("*Rain*: %vmm", wtr.Rain.ThreeH),
								false,
								false,
							),
							slack.NewTextBlockObject(
								slack.MarkdownType,
								fmt.Sprintf(
									"*Windspeed*: %vkm/h %v",
									windSpeedKPH,
									windDirection,
								),
								false,
								false,
							),
						},
						nil,
					),
				},
			},
		},
	)
	if err != nil {
		return errors.Wrap(err, "posting to slack")
	}

	return nil
}

func directionFromDegrees(deg float64) string {
	idx := int(math.Round(math.Mod(deg, 360) / 22.5))

	return []string{
		"N",
		"NNE",
		"NE",
		"ENE",
		"E",
		"ESE",
		"SE",
		"SSE",
		"S",
		"SSW",
		"SW",
		"WSW",
		"W",
		"WNW",
		"NW",
		"NNW",
	}[idx]
}

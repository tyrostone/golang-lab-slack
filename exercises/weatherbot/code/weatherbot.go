package main

import (
	"strings"

	"golang.org/x/net/context"

	slackbot "github.com/BeepBoopHQ/go-slackbot"
	"github.com/briandowns/openweathermap"
	"github.com/nlopes/slack"
	"log"
	"os"
	"fmt"
)

func main() {
	ListenForHelloOrWeather()
}

func ListenForHelloOrWeather() {
	bot := slackbot.New(os.Getenv("SLACK_API_TOKEN"))
	directMessage := bot.Messages(slackbot.DirectMessage, slackbot.DirectMention).Subrouter()
	directMessage.Hear("(?i)hi|hello|howdy.*").MessageHandler(HelloHandler)
	bot.Hear("(?i)weather (.*)").MessageHandler(WeatherHandler)
	bot.Hear("(?i)weather").MessageHandler(WeatherHelpHandler)
	bot.Run()

}

func HelloHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	bot.Reply(evt, "Hi there!", true)
}

func WeatherHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	place := strings.Split(evt.Msg.Text, " ")[1]
	w, err := GetWeather(place)
	if err != nil {
		log.Fatal("Failed to get weather: ", err)
	}
	bot.Reply(evt,
		fmt.Sprintf("The current temperature for %s is %.0f degrees fahrenheit (%s).",
			w.Name, w.Main.Temp, w.Weather[0].Description), true)
}

func GetWeather(place string) (weather *openweathermap.CurrentWeatherData, error error){
	weather, err := openweathermap.NewCurrent("F", "EN")
	if err != nil {
		log.Fatal("Failed to create weather client: ", err)
	}
	err = weather.CurrentByName(place)
	return weather, err
}

func WeatherHelpHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	bot.Reply(evt, "To get the current weather for an area, type 'weather <ZIPCODE>'", true)
}

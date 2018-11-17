// Copyright © 2017 Alessandro Sanino <saninoale@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package examples

import (
	"fmt"
	"log"
	"time"

	"github.com/nlopes/slack"
	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/saniales/golang-crypto-trading-bot/exchanges"
	"github.com/saniales/golang-crypto-trading-bot/strategies"
	"github.com/shomali11/slacker"
	"github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	Watch5Sec                  = watch5Sec()
	SlackIntegrationExample    = slackIntegrationExample()
	TelegramIntegrationExample = telegramIntegrationExample()
)

// Watch5Sec prints out the info of the market every 5 seconds.
func watch5Sec() strategies.Strategy {
	return strategies.IntervalStrategy{
		Model: strategies.StrategyModel{
			Name: "Watch5Sec",
			Setup: func(wrappers []exchanges.ExchangeWrapper, markets []*environment.Market) error {
				fmt.Println("Watch5Sec starting")
				return nil
			},
			OnUpdate: func(wrappers []exchanges.ExchangeWrapper, markets []*environment.Market) error {
				_, err := wrappers[0].GetMarketSummary(markets[0])
				if err != nil {
					return err
				}
				logrus.Info(markets)
				logrus.Info(wrappers)
				return nil
			},
			OnError: func(err error) {
				fmt.Println(err)
			},
			TearDown: func(wrappers []exchanges.ExchangeWrapper, markets []*environment.Market) error {
				fmt.Println("Watch5Sec exited")
				return nil
			},
		},
		Interval: time.Second * 5,
	}
}

var slackBot *slacker.Slacker

// SlackIntegrationExample send messages as a strategy.
// RTM not supported (and usually not requested when trading, this is an automated slackBot).
func slackIntegrationExample() strategies.Strategy {
	return strategies.IntervalStrategy{
		Model: strategies.StrategyModel{
			Name: "SlackIntegrationExample",
			Setup: func([]exchanges.ExchangeWrapper, []*environment.Market) error {
				// connect slack token
				slackBot = slacker.NewClient("YOUR-TOKEN-HERE")
				slackBot.Init(func() {
					log.Println("Slack BOT Connected")
				})
				slackBot.Err(func(err string) {
					log.Println("Error during slack slackBot connection: ", err)
				})
				go func() {
					err := slackBot.Listen()
					if err != nil {
						log.Fatal(err)
					}
				}()
				return nil
			},
			OnUpdate: func([]exchanges.ExchangeWrapper, []*environment.Market) error {
				//if updates has requirements
				_, _, err := slackBot.Client.PostMessage("DESIRED-CHANNEL", "OMG something happening!!!!!", slack.PostMessageParameters{})
				return err
			},
			OnError: func(err error) {
				logrus.Errorf("I Got an error %s", err)
			},
		},
		Interval: time.Second * 10,
	}
}

var telegramBot *tb.Bot

func telegramIntegrationExample() strategies.Strategy {
	return strategies.IntervalStrategy{
		Model: strategies.StrategyModel{
			Name: "TelegramIntegrationExample",
			Setup: func([]exchanges.ExchangeWrapper, []*environment.Market) error {
				telegramBot, err := tb.NewBot(tb.Settings{
					Token:  "692317936:AAGgC-IFG5M5PBZquTJzl4a83uWUF46eUj8",
					Poller: &tb.LongPoller{Timeout: 10 * time.Second},
				})

				if err != nil {
					return err
				}

				telegramBot.Start()
				return nil
			},
			OnUpdate: func([]exchanges.ExchangeWrapper, []*environment.Market) error {
				telegramBot.Send(&tb.User{
					Username: "YOUR-USERNAME-GROUP-OR-USER",
				}, "OMG SOMETHING HAPPENING!!!!!", tb.SendOptions{})

				/*
					// Optionally it can have options
					telegramBot.Send(tb.User{
						Username: "YOUR-JOINED-GROUP-USERNAME",
					}, "OMG SOMETHING HAPPENING!!!!!", tb.SendOptions{})
				*/
				return nil
			},
			OnError: func(err error) {
				logrus.Errorf("I Got an error %s", err)
				telegramBot.Stop()
			},
			TearDown: func([]exchanges.ExchangeWrapper, []*environment.Market) error {
				telegramBot.Stop()
				return nil
			},
		},
	}
}

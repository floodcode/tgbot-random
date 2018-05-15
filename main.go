package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/floodcode/tgbot"
)

var (
	bot     tgbot.TelegramBot
	botUser tgbot.User
)

type botConfig struct {
	Token string `json:"token"`
}

func main() {
	rand.Seed(time.Now().Unix())

	configData, err := ioutil.ReadFile("config.json")
	checkError(err)

	var config botConfig
	err = json.Unmarshal(configData, &config)
	checkError(err)

	bot, err = tgbot.New(config.Token)
	checkError(err)

	botUser, err = bot.GetMe()
	checkError(err)

	bot.Poll(tgbot.PollConfig{
		Delay:    100,
		Callback: updatesCallback,
	})
}

func updatesCallback(updates []tgbot.Update) {
	for _, update := range updates {
		if update.Message == nil || len(update.Message.Text) == 0 {
			continue
		}

		processTextMessage(update.Message)
	}
}

func processTextMessage(msg *tgbot.Message) {
	cmdMatch, _ := regexp.Compile(`^\/([a-zA-Z_]+)(?:@` + botUser.Username + `)?(?:\s(.+)|)$`)
	match := cmdMatch.FindStringSubmatch(msg.Text)

	if match == nil {
		return
	}

	command := strings.ToLower(match[1])

	if command == "start" || command == "help" {
		sendExamples(msg)
		return
	}

	if command == "random" {
		if len(match[2]) == 0 {
			sendExamples(msg)
			return
		}

		sendRandom(msg, match[2])
	}
}

func sendExamples(msg *tgbot.Message) {
	text := strings.Join([]string{
		"_Usage examples:_",
		"`/random 1-10`",
		"`/random apple|pear|lemon`",
	}, "\n")

	bot.SendMessage(tgbot.SendMessageConfig{
		ChatID:           tgbot.ChatID(msg.Chat.ID),
		Text:             text,
		ReplyToMessageID: msg.MessageID,
		ParseMode:        tgbot.ParseModeMarkdown(),
	})
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func generateRandom(pattern string) (string, bool) {
	rangeMatch, _ := regexp.Compile(`^([0-9]+)-([0-9]+)$`)
	if match := rangeMatch.FindStringSubmatch(pattern); match != nil {
		min, _ := strconv.Atoi(match[1])
		max, _ := strconv.Atoi(match[2])

		result := fmt.Sprintf(strings.Join([]string{
			"_Random number betwen %d and %d:_",
			"*%d*",
		}, "\n"), min, max, random(min, max))

		return result, true
	}

	parts := strings.Split(pattern, "|")
	if len(parts) > 1 {
		tokens := []string{}

		for _, part := range parts {
			token := strings.TrimSpace(part)
			if len(token) > 0 {
				tokens = append(tokens, token)
			}
		}

		if len(tokens) <= 1 {
			return "", false
		}

		randomToken := tokens[random(0, len(tokens))]

		result := fmt.Sprintf(strings.Join([]string{
			"_Random item:_",
			"*%s*",
		}, "\n"), randomToken)

		return result, true
	}

	return "", false
}

func sendRandom(msg *tgbot.Message, pattern string) {
	var messageText string
	result, ok := generateRandom(pattern)
	if ok {
		messageText = result
	} else {
		messageText = "Invalid pattern"
		defer sendExamples(msg)
	}

	bot.SendMessage(tgbot.SendMessageConfig{
		ChatID:           tgbot.ChatID(msg.Chat.ID),
		Text:             messageText,
		ReplyToMessageID: msg.MessageID,
		ParseMode:        tgbot.ParseModeMarkdown(),
	})
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

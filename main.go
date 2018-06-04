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

	"github.com/floodcode/tbf"
)

type botConfig struct {
	Token string `json:"token"`
	Delay int    `json:"delay"`
}

func main() {
	rand.Seed(time.Now().Unix())

	configData, err := ioutil.ReadFile("config.json")
	checkError(err)

	var config botConfig
	err = json.Unmarshal(configData, &config)
	checkError(err)

	bot, err := tbf.New(config.Token)
	checkError(err)

	bot.AddRoute("start", helpAction)
	bot.AddRoute("help", helpAction)
	bot.AddRoute("random", randomAction)

	bot.Poll(tbf.PollConfig{
		Delay: config.Delay,
	})
}

func helpAction(req tbf.Request) {
	req.QuickMessageMD(strings.Join([]string{
		"Usage examples:",
		"/random `1-10` or `apple|pear|lemon`",
	}, "\n"))
}

func randomAction(req tbf.Request) {
	pattern := req.Args
	if len(pattern) == 0 {
		req.QuickMessage("Enter the pattern:")
		pattern = req.WaitNext().Message.Text
	}

	result, ok := generateRandom(pattern)
	if ok {
		req.QuickMessageMD(result)
		return
	}

	req.QuickMessage("Invalid pattern")
	helpAction(req)
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

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

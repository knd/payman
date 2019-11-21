package twitter

import (
	"fmt"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Bot is a wrapper for a twitter session to post payout results to
type Bot struct {
	prefix  string
	session *twitter.Client
}

// NewBot creates a new twitter bot.
func NewBot(prefix, consumerKey, consumerSecret, accessKey, accessSecret string) (*Bot, error) {
	bot := Bot{prefix: prefix}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessKey, accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	bot.session = client
	return &bot, nil
}

// Post posts a tezblock link to the operation hash
func (bot *Bot) Post(ophash string, cycle int) error {
	ophash = ophash[1 : len(ophash)-2]
	link := fmt.Sprintf("https://mvp.tezblock.io/transaction/%s", ophash)

	var sb strings.Builder
	sb.WriteString(bot.prefix)
	sb.WriteString(fmt.Sprintf("Payout for Cycle %d: %s", cycle, link))

	_, _, err := bot.session.Statuses.Update(sb.String(), nil)
	if err != nil {
		return err
	}

	return nil
}

package bitbot

import (
	"github.com/whyrusleeping/hellabot"
)

var MessageCounterTrigger = NamedTrigger{
	ID: "messageCounter",
	Condition: func(irc *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG"

	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		b.counters["messageCounter"].Inc()
		return true
	},
}
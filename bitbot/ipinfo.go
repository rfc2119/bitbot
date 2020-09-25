package bitbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/whyrusleeping/hellabot"
)

type GeoData struct {
	IP       string
	Hostname string
	City     string
	Region   string
	Country  string
	Loc      string
	Org      string
	Postal   string
	Timezone string
	Readme   string
}

var IPinfoTrigger = NamedTrigger{ //nolint:gochecknoglobals,golint
	ID:   "ipinfo",
	Help: "!ipinfo <valid IP>",
	Condition: func(irc *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, "!ipinfo")
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		var resp string
		cmd := strings.Split(m.Content, " ")
		if len(cmd) > 1 {
			resp = query(cmd[1])
		} else {
			resp = "please provide an ip...ya twatsicle"
		}
		irc.Reply(m, resp)
		return true
	},
}

func decodeJSON(encodedJSON []byte) string {
	var (
		ipinfo GeoData
		reply  string
	)

	err := json.Unmarshal(encodedJSON, &ipinfo)
	if err != nil {
		b.Config.Logger.Warn("IPinfo trigger, couldn't decode JSON", "error", err)
	}

	if ipinfo.IP == "" {
		reply = "either the IP was not valid or we are being rate limited"
	} else {
		reply = fmt.Sprintf("ip: %s\nhostname: %s\ncity: %s\nregion: %s\ncountry: %s\ncoords: %s\norg: %s\npostal: %s\ntimezone: %s",
			ipinfo.IP,
			ipinfo.Hostname,
			ipinfo.City,
			ipinfo.Region,
			ipinfo.Country,
			ipinfo.Loc,
			ipinfo.Org,
			ipinfo.Postal,
			ipinfo.Timezone)
	}
	return reply
}

func query(ip string) string {
	url := "http://ipinfo.io/" + ip
	res, err := http.Get(url)
	if err != nil {
		b.Config.Logger.Warn("IPinfo trigger, couldn't query ipinfo.io", "error", err)
	}
	jsonData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		b.Config.Logger.Warn("IPinfo trigger, couldn't read ipinfo.io answer", "error", err)
	}

	res.Body.Close() //nolint:errcheck,gosec
	return decodeJSON(jsonData)
}

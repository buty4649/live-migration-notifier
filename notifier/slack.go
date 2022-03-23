package notifier

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type block map[string]interface{}

func section(text block) block {
	return block{"type": "section", "text": text}
}

func plainText(text string, emoji bool) block {
	return block{"type": "plain_text", "text": text, "emoji": emoji}
}

func markdown(text string) block {
	return block{"type": "mrkdwn", "text": text}
}

func divider() block {
	return block{"type": "divider"}
}

func context(elements []block) block {
	return block{"type": "context", "elements": elements}
}

func postSlack(data *LiveMigrationData, webhook_url string) error {
	payload := map[string][]block{
		"blocks": {
			section(markdown(data.Message)),
			divider(),
			section(plainText(fmt.Sprintf(":desktop_computer: %s", data.Hostname), true)),
			context([]block{plainText(data.UUID, false)}),
			section(plainText(fmt.Sprintf(":package: %s → %s", data.Src, data.Dest), true)),
			divider(),
			context([]block{
				plainText(data.RequestId, false),
				plainText(fmt.Sprintf("by %s", data.User), false),
			}),
		},
	}

	p, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.PostForm(webhook_url, url.Values{"payload": {string(p)}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

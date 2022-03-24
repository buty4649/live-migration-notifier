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

func markdown(text string, verbatim bool) block {
	return block{"type": "mrkdwn", "text": text, "verbatim": verbatim}
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
			section(plainText(data.Message, true)),
			divider(),
			context([]block{
				markdown(
					fmt.Sprintf(":desktop_computer: *%s*\n:pencil2: %s", data.Hostname, data.UUID),
					true,
				),
			}),
			context([]block{
				plainText(fmt.Sprintf(":outbox_tray: %s", data.Src), true),
				plainText(fmt.Sprintf(":inbox_tray: %s", data.Dest), true),
			}),
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

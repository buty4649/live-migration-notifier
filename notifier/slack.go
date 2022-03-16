package notifier

import (
	"encoding/json"
	"net/http"
	"net/url"
)

func postSlack(data *LiveMigrationData, webhook_url string) error {
	payload := map[string]interface{}{
		"text":  data.Message,
		"color": data.Color,
		"fields": []map[string]interface{}{
			{
				"title": "instance",
				"value": data.Hostname,
				"short": true,
			},
			{
				"title": "user",
				"value": data.User,
				"short": true,
			},
			{
				"title": "src",
				"value": data.Src,
				"short": true,
			},
			{
				"title": "dest",
				"value": data.Dest,
				"short": true,
			},
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

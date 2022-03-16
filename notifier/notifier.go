package notifier

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

type Notifier struct {
	uri         string
	webhook_url string
}

type OsloMessage struct {
	Message string `json:"oslo.message"`
	Version string `json:"oslo.version"`
}

type Payload struct {
	Args             PayloadArgs `json:"args"`
	Method           string      `json:"method"`
	Version          string      `json:"version"`
	ContextUserName  string      `json:"_context_user_name"`
	ContextRequestId string      `json:"_context_request_id"`
}

type PayloadArgs struct {
	Dest     string     `json:"dest"`
	Instance NovaObject `json:"instance"`
}

type NovaObject struct {
	Changes   []string               `json:"nova_object.changes"`
	Data      map[string]interface{} `json:"nova_object.data"`
	Name      string                 `json:"nova_object.name"`
	Namespace string                 `json:"nova_object.namespace"`
	Version   string                 `json:"nova_object.version"`
}

type LiveMigrationData struct {
	Message   string
	Color     string
	Src       string
	Dest      string
	Hostname  string
	User      string
	RequestId string
	UUID      string
}

func Init(uri, webhook_url string) *Notifier {
	return &Notifier{
		uri:         uri,
		webhook_url: webhook_url,
	}
}

func (n *Notifier) Start() error {
	conn, err := amqp.Dial(n.uri)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare("live-migration-notifier", false, false, true, false, nil)
	if err != nil {
		return err
	}

	err = ch.QueueBind(queue.Name, "compute.*", "nova", false, nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			payload, err := extractPayload(d.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				continue
			}

			dest := strings.TrimPrefix(d.RoutingKey, "compute.")
			data := extractMessage(dest, payload)
			if data != nil {
				postSlack(data, n.webhook_url)
			}
		}
	}()

	<-forever
	return nil
}

func extractPayload(data []byte) (*Payload, error) {
	var omsg OsloMessage
	err := json.Unmarshal(data, &omsg)
	if err != nil {
		return nil, err
	}

	var payload Payload
	err = json.Unmarshal([]byte(omsg.Message), &payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}

func extractMessage(dest string, payload *Payload) *LiveMigrationData {
	var data LiveMigrationData

	switch payload.Method {
	case "live_migration":
		data.Message = "live-migrationが開始しました"
		data.Color = "warning"

	case "post_live_migration_at_destination":
		data.Message = "live-migrationが完了しました"
		data.Color = "good"

	case "rollback_live_migration_at_destination":
		data.Message = "live-migrationが失敗しました"
		data.Color = "danger"

	default:
		return nil
	}

	data.Src = payload.Args.Instance.Data["host"].(string)
	if payload.Args.Dest != "" {
		data.Dest = payload.Args.Dest
	} else {
		data.Dest = dest
	}
	data.Hostname = payload.Args.Instance.Data["display_name"].(string)
	data.User = payload.ContextUserName
	data.RequestId = payload.ContextRequestId
	data.UUID = payload.Args.Instance.Data["uuid"].(string)

	return &data
}

package cmdutil

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestReadConfigFile(t *testing.T) {
	file, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(file.Name())

	_, err = file.WriteString(`rabbitmq_host: example.com
rabbitmq_port: 5672
rabbitmq_user: test
rabbitmq_password: test
slack_webhook_url: https://example.com/webhook
`)
	if err != nil {
		t.Error(err)
	}
	file.Close()

	got, err := ReadConfigFile(file.Name())
	if err != nil {
		t.Error(err)
	}

	expect := "amqp://test:test@example.com:5672/"
	if got.Uri() != expect {
		t.Errorf("got: %v; expect: %s", got.Uri(), expect)
	}

	expect = "https://example.com/webhook"
	if got.SlackWebHookUrl != expect {
		t.Errorf("got: %v; expect: %s", got.SlackWebHookUrl, expect)
	}
}

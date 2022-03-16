package notifier

import "testing"

func TestExtractPayload(t *testing.T) {
	data := `{
		"oslo.message": "{\"method\": \"test\"}",
		"oslo.version": "2.0"
}`

	got, err := extractPayload([]byte(data))
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if got.Method != "test" {
		t.Errorf("got: %v; expect: test", got)
	}
}

func TestExtractMessage(t *testing.T) {
	payload := Payload{
		Args: PayloadArgs{
			Instance: NovaObject{
				Data: map[string]interface{}{
					"host":         "test",
					"display_name": "testtest",
					"uuid":         "33720cb6-f6f0-40fa-93fc-b38306988798",
				},
			},
		},
		ContextUserName:  "test",
		ContextRequestId: "test",
	}

	payload.Method = "live_migration"
	data := extractMessage("test", &payload)

	got := data.Message
	expect := "live-migrationが開始しました"
	if got != expect {
		t.Errorf("got: %v; expect: %s", got, expect)
	}

	got = data.Color
	expect = "warning"
	if got != expect {
		t.Errorf("got: %v; expect: %s", got, expect)
	}

	payload.Method = "post_live_migration_at_destination"
	data = extractMessage("test", &payload)

	got = data.Message
	expect = "live-migrationが完了しました"
	if got != expect {
		t.Errorf("got: %v; expect: %s", got, expect)
	}

	got = data.Color
	expect = "good"
	if got != expect {
		t.Errorf("got: %v; expect: %s", got, expect)
	}

	payload.Method = "rollback_live_migration_at_destination"
	data = extractMessage("test", &payload)

	got = data.Message
	expect = "live-migrationが失敗しました"
	if got != expect {
		t.Errorf("got: %v; expect: %s", got, expect)
	}

	got = data.Color
	expect = "danger"
	if got != expect {
		t.Errorf("got: %v; expect: %s", got, expect)
	}
}

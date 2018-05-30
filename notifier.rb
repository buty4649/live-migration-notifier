#!/usr/bin/env ruby

require 'bunny'
require 'json'
require 'yaml'
require 'slack/incoming/webhooks'

config = YAML.load_file("config.yml")

def rabbitmq_init
  host = config['RABBITMQ_HOST']
  port = config['RABBITMQ_PORT'] || 5672
  user = config['RABBITMQ_USER'] || 'guest'
  pass = config['RABBITMQ_PASSWORD'] || 'guest'

  Bunny.new(hostname: host, port: port, user: user, pass: pass)
end

webhook_url = config["SLACK_WEBHOOK_URL"]
slack = Slack::Incoming::Webhooks.new(webhook_url)

connection = rabbitmq_init
connection.start

ch = connection.create_channel
q = ch.queue("live-migration_notifier", exclusive: true)
x = ch.topic("nova")
q.bind(x, routing_key: "compute.*")

begin
  q.subscribe(manual_ack: true, block: true) do |delivery_info, properties, msg|
    payload = JSON.parse(JSON.parse(msg)["oslo.message"])
    method = payload["method"]

    case method
    when "pre_live_migration"
      instance = payload["args"]["instance"]["nova_object.data"]

      src = instance["host"]
      dst = delivery_info.routing_key.gsub(/^compute\./, "")

      hostname = instance["display_name"]
      user = payload["_context_user_name"]

      attachment = [{
        text: "live-migrationが開始されました",
        fallback: "start: #{hostname} (#{src} -> #{dst})",
        color: "warning",
        fields: [
          {
            title: "instance",
            value: hostname,
            short: true,
          },
          {
            title: "user",
            value: user,
            short: true,
          },
          {
            title: "src",
            value: src,
            short: true
          },
          {
            title: "dst",
            value: dst,
            short: true
          },
        ]
      }]
    when "post_live_migration_at_destination"
      instance = payload["args"]["instance"]["nova_object.data"]

      src = instance["host"]
      dst = delivery_info.routing_key.gsub(/^compute\./, "")

      hostname = instance["display_name"]
      user = payload["_context_user_name"]

      attachment = [{
        text: "live-migrationが完了しました",
        fallback: "end: #{hostname} (#{dst})",
        color: "good",
        fields: [
          {
            title: "instance",
            value: hostname,
            short: true,
          },
          {
            title: "user",
            value: user,
            short: true,
          },
          {
            title: "src",
            value: src,
            short: true
          },
          {
            title: "dst",
            value: dst,
            short: true
          },
        ]
      }]
    when "rollback_live_migration_at_destination"
      instance = payload["args"]["instance"]["nova_object.data"]

      src = instance["host"]
      dst = delivery_info.routing_key.gsub(/^compute\./, "")

      hostname = instance["display_name"]
      user = payload["_context_user_name"]

      attachment = [{
        text: "live-migrationが失敗しました",
        fallback: "error: #{hostname} (#{dst})",
        color: "danger",
        fields: [
          {
            title: "instance",
            value: hostname,
            short: true,
          },
          {
            title: "user",
            value: user,
            short: true,
          },
          {
            title: "src",
            value: src,
            short: true
          },
          {
            title: "dst",
            value: dst,
            short: true
          },
        ]
      }]
    end
    slack.post "", attachments: attachment
  end
rescue Interrupt
  connection.close
end

#!/usr/bin/env ruby

require 'bunny'
require 'json'
require 'yaml'
require 'slack/incoming/webhooks'

config_file = ARGV[0] || File.dirname(__FILE__) + '/config.yml'
config = YAML.load_file(config_file)

def rabbitmq_init(config)
  host = config['rabbitmq_host']
  port = config['rabbitmq_port'] || 5672
  user = config['rabbitmq_user'] || 'guest'
  pass = config['rabbitmq_password'] || 'guest'

  Bunny.new(hostname: host, port: port, user: user, pass: pass)
end

webhook_url = config["slack_webhook_url"]
slack = Slack::Incoming::Webhooks.new(webhook_url)

connection = rabbitmq_init(config)
connection.start

ch = connection.create_channel
q = ch.queue("live-migration_notifier", exclusive: true)
x = ch.topic("nova")
q.bind(x, routing_key: "compute.*")

begin
  q.subscribe(manual_ack: false, block: true) do |delivery_info, properties, msg|
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
      instance_uuid = instance["uuid"]
      user = payload["_context_user_name"]
      req_id = payload["_context_request_id"]

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
          {
            title: "request-id",
            value: req_id,
            short: false,
          },
          {
            title: "uuid",
            value: instance_uuid,
            short: false,
          }
        ]
      }]
    end
    slack.post "", attachments: attachment if attachment
  end
rescue Interrupt
  connection.close
end

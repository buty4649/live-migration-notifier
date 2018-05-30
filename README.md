# live-migration-notifier
OpenStackのLiveMigrationの通知するくん

## 使い方

```
$ export RABBITMQ_HOST=<rabbitmq host>
$ export RABBITMQ_PORT=<rabbitmq port>
$ export RABBITMQ_USER=<rabbitmq user>
$ export RABBITMQ_PASSWORD=<rabbitmq password>
$ export SLACK_WEBHOOK_URL=<webhook url>

$ bundle install
$ bundle exec ./notifier.rb
```

systemdで動かす場合には以下のようなUnitファイルを作ります

```
[Unit]
Description = live-migration notifier

[Service]
User=ubuntu
Group=ubuntu
Environment = RABBITMQ_HOST=<rabbitmq host>
Environment = RABBITMQ_PORT=<rabbitmq port>
Environment = RABBITMQ_USER=<rabbitmq user>
Environment = RABBITMQ_PASSWORD=<rabbitmq password>
Environment = SLACK_WEBHOOK_URL=<webhook url>
ExecStart = /usr/bin/ruby /path/to/notifier.rb
Restart = always
Type = simple

[Install]
WantedBy = multi-user.target
```

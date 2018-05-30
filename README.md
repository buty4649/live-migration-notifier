# live-migration-notifier

OpenStackのLiveMigrationの通知するくん

## 使い方

```sh
# config.yml の情報を埋めます
$ cp config.yml.sample config.yml

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
ExecStart = /usr/bin/ruby /path/to/notifier.rb
Restart = always
Type = simple

[Install]
WantedBy = multi-user.target
```

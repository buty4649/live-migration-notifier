
release:
	git push origin master --tags
	docker run --rm -v $(pwd):/go/src/github.com/buty4649/live-migration-notifier -w /go/src/github.com/buty4649/live-migration-notifier -e GITHUB_TOKEN goreleaser/goreleaser:latest release --rm-dist

DOKCER_IMAGE := goreleaser/goreleaser:latest
DOCKER_COMMAND := docker run --rm -v $(shell pwd):/go/src/github.com/buty4649/live-migration-notifier -w /go/src/github.com/buty4649/live-migration-notifier -e GITHUB_TOKEN $(DOKCER_IMAGE)

release: update_image
	git push origin master --tags
	$(DOCKER_COMMAND) release --rm-dist

release_test: update_image
	$(DOCKER_COMMAND) release --rm-dist --skip-validate --skip-publish

update_image:
	docker pull $(DOKCER_IMAGE)

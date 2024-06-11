snapshot:
	goreleaser release --snapshot --skip-publish --rm-dist

release:
	goreleaser release --rm-dist
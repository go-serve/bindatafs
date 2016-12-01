#
# This file is only used to standardize testing
# environment. It doesn't build any binary. Nor
# does it needed in the installation process.
#
# For installation details, please read README.md
#

test: timestamp
	@echo
	@echo "== Run tests"
	go test -v -cover ./...

generate: timestamp
	@echo
	@echo "== Generate assets.go"
	go generate ./examples/...

timestamp:
	@echo
	@echo "== Ensure timestamp of local assets"
	TZ=Asia/Hong_Kong find ./examples/. -type f -exec touch -t 201611210125.30 "{}" \;

.PHONY: test generate timestamp

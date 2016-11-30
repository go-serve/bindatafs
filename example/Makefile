all: timestamp build
	@echo

build:
	@echo
	@echo "* Generate assets.go"
	go generate

timestamp:
	@echo
	@echo "* Ensure timestamp of local assets"
	TZ=Asia/Hong_Kong find -type f -exec touch -t 201611210125.30 "{}" \;

.PHONY: all build timestamp

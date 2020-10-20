PLUGIN_NAME=archive

all: protos build

protos:
	@echo ""
	@echo "Build Protos"

	protoc -I . --go_opt=paths=source_relative --go_out=. ./builder/output.proto

.PHONY: build
build:
	@echo ""
	@echo "Compile Plugin"

	go build -o ./bin/waypoint-plugin-${PLUGIN_NAME} ./main.go

.PHONY: install
install: build
	@echo ""
	@echo "Installing Plugin"

	cp ./bin/waypoint-plugin-${PLUGIN_NAME} ${HOME}/.config/waypoint/plugins/
	# For MacOS Big Sur
	cp ./bin/waypoint-plugin-${PLUGIN_NAME} /Users/${USER}/Library/Preferences/waypoint/plugins/

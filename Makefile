DIR ?= "$$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PI_IP ?= 192.168.1.11
BIN ?= ghost-detector

.PHONY: build.arm
build.arm:
	GOOS=linux GOARCH=arm GOARM=5 go build -o build/$(BIN)-arm `go list ./cmd/main/`

.PHONY: build.x86
build.x86:
	go build -o build/$(BIN) `go list ./cmd/main/`

sync: 
	@echo "Syncing $(DIR) to $(PI_IP)..."
	@rsync -avr --progress $(DIR)/build/*-arm pi@$(PI_IP):/home/pi/.
	@rsync -avr --progress $(DIR)/words pi@$(PI_IP):/home/pi/.
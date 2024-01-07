INSTALL_PATH := $(HOME)/.local/bin

build: cmd/NastySecrets/main.go
	@echo "Building --> bin/NastySecrets"
	@go build -o bin/NastySecrets cmd/NastySecrets/main.go

perm: bin/NastySecrets
	@echo "Setting permission --> 700 (-rwx------)"
	@chmod 700 bin/NastySecrets

install: bin/NastySecrets
	@echo "Copying binary --> $(INSTALL_PATH)"
	@cp bin/NastySecrets $(INSTALL_PATH)

done: $(INSTALL_PATH)/NastySecrets
	@echo "Done"

all: build perm install done

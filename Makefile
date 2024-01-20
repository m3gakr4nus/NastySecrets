INSTALL_PATH := $(HOME)/.local/bin
MAN_PATH := $(HOME)/.local/share/man/man1

build: cmd/NastySecrets/main.go
	@echo "[*] Building --> bin/nastysecrets"
	@go build -o bin/nastysecrets cmd/NastySecrets/main.go

perm: bin/nastysecrets
	@echo "[*] Setting permission --> 700 (-rwx------)"
	@chmod 700 bin/nastysecrets

install: bin/nastysecrets
	@echo "[*] Copying binary --> $(INSTALL_PATH)"
	@cp bin/nastysecrets $(INSTALL_PATH)

man: docs/manpage/nastysecrets
	@echo "[*] Coping manpage --> $(MAN_PATH)"
	@cp docs/manpage/nastysecrets $(MAN_PATH)/nastysecrets.1
	
done: $(INSTALL_PATH)/nastysecrets
	@echo "[+] Done"

uninstall:
	@echo "[*] Removing manpage ($(MAN_PATH)/nastysecrets.1)"
	@rm $(MAN_PATH)/nastysecrets.1
	@echo "[*] Removing from PATH ($(INSTALL_PATH)/nastysecrets)"
	@rm $(INSTALL_PATH)/nastysecrets
	@echo "[+] Done"

all: build perm install man done

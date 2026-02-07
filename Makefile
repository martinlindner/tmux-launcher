.PHONY: build install clean

BINARY = tmux-launcher
INSTALL_DIR = $(HOME)/.local/bin

build:
	go build -o $(BINARY) .

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY) $(INSTALL_DIR)/

clean:
	rm -f $(BINARY)
	go clean

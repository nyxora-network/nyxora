.PHONY: build install clean run test

BINARY=nyxora
INSTALL_PATH=/usr/local/bin

build:
	go build -ldflags="-s -w" -o $(BINARY) ./cmd/nyxora

install: build
	cp $(BINARY) $(INSTALL_PATH)/$(BINARY)
	mkdir -p /etc/nyxora/tunnels /etc/nyxora/cache /var/log/nyxora
	cp -r tunnels/*.tar.gz /etc/nyxora/tunnels/ 2>/dev/null || true

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)
	rm -rf /etc/nyxora/cache/*

test:
	go test ./...

vet:
	go vet ./...

tunnels:
	@for dir in tunnels/*/; do \
		name=$$(basename $$dir); \
		echo "packing $$name..."; \
		(cd tunnels && tar czf /etc/nyxora/tunnels/$$name.tar.gz $$name); \
	done

daemon: build
	systemctl daemon-reload
	systemctl enable nyxora
	systemctl restart nyxora
	systemctl status nyxora --no-pager

uninstall:
	systemctl stop nyxora 2>/dev/null || true
	systemctl disable nyxora 2>/dev/null || true
	rm -f $(INSTALL_PATH)/$(BINARY)
	rm -f /etc/systemd/system/nyxora.service
	rm -rf /etc/nyxora
	rm -rf /var/log/nyxora

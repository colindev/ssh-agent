
APP := ssh-agent
DES := /usr/local/bin/
AGENT ?= 127.0.0.1:6666
SERVICE ?= ssh-agent-server
AUTHORIZATION ?= ssh-authorization.sh

build:
	go get -a ./... && go build -o ./$(APP)

build-via-docker:
	docker run --rm -v `pwd`:/go/src/app -w /go/src/app golang:1.8 go build -o $(APP)

install:
	cp $(APP) $(DES) && chmod 700 $(DES)$(APP) && chown root:root $(DES)$(APP)
	touch /etc/$(SERVICE).conf
	cat ssh-agent-server.service | sed 's/{APP}/'$(APP)'/g' | sed 's/{AGENT}/'$(AGENT)'/g' > /etc/systemd/system/$(SERVICE).service
	systemctl daemon-reload
	systemctl enable $(SERVICE)
	systemctl start $(SERVICE)
	systemctl status $(SERVICE)

upgrade:
	systemctl stop $(SERVICE)
	$(MAKE) install -e AGENT=$(AGENT) -e SERVICE=$(SERVICE)

uninstall:
	systemctl stop $(SERVICE)
	systemctl disable $(SERVICE)
	rm /etc/systemd/system/$(SERVICE).service $(DES)$(APP)
	systemctl daemon-reload

set-sshd:
	cp /etc/ssh/sshd_config /etc/ssh/sshd_config.bak
	sed -i 's@#\?AuthorizedKeysCommand\s\+[^#]\+@AuthorizedKeysCommand '$(DES)$(AUTHORIZATION)'@' /etc/ssh/sshd_config
	sed -i 's/#\?AuthorizedKeysCommandUser\s\+[^#]\+/AuthorizedKeysCommandUser root/' /etc/ssh/sshd_config
	systemctl reload sshd

install-authorization: set-sshd
	cat scripts/auth.sh | sed 's/{AGENT}/'$(AGENT)'/g' > $(DES)$(AUTHORIZATION) && chmod 700 $(DES)$(AUTHORIZATION) && chown root:root $(DES)$(AUTHORIZATION)

install-authorization-gcp: set-sshd
	cat scripts/auth-gcp.sh | sed 's/{AGENT}/'$(AGENT)'/g' > $(DES)$(AUTHORIZATION) && chmod 700 $(DES)$(AUTHORIZATION) && chown root:root $(DES)$(AUTHORIZATION)

uninstall-authorization:
	rm $(DES)$(AUTHORIZATION)

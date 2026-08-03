deploy:
	cd ../SuConfig/linux/asus &&\
	make docker-su-home-okaeri docker-su-home-okaeri-monitor
deploy-home:
	cd ../SuConfig/linux/asus &&\
	make docker-su-home-okaeri

dev:
	aira

include ../SuConfig/linux/asus/config.mk
include ../SuConfig/deploy/import.mk

deploy-tool:
	$(call upload, tool/, TOOL/su-home-okaeri/)

# deploy:
build:
	CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -v -o su-home-okaeri-monitor cmd/monitor/main.go
	CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -v -o su-home-okaeri cmd/home/main.go

deploy-monitor:
	cd ../SuConfig/linux/asus &&\
	make docker-su-home-okaeri-monitor
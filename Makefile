deploy:
	cd ../SuConfig/linux/asus &&\
	make docker-su-home-okaeri

dev:
	aira

include ../SuConfig/linux/asus/config.mk
include ../SuConfig/deploy/import.mk

deploy-tool:
	$(call upload, tool/, TOOL/su-home-okaeri/)

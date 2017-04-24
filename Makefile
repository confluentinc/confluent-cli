PACKAGE_NAME ?= confluent-cli
VERSION ?= SNAPSHOT

.PHONY:

all: build

install:
ifndef CONFLUENT_HOME
	$(error Cannot install. CONFLUENT_HOME is not set)
endif
	install -m 755 bin/confluent $(CONFLUENT_HOME)/bin/confluent

build: oss

prep:
	mkdir -p bin/

oss: prep
	cp -f src/oss/confluent.sh bin/confluent
	chmod 755 bin/confluent

platform: prep
	cp -f src/platform/confluent.sh bin/confluent
	chmod 755 bin/confluent

clean:
	rm -rf bin/

distclean: clean

archive:
	git archive --prefix=$(PACKAGE_NAME)-$(VERSION)/ \
		-o $(PACKAGE_NAME)-$(VERSION).tar.gz HEAD
	git archive --prefix=$(PACKAGE_NAME)-$(VERSION)/ \
		-o $(PACKAGE_NAME)-$(VERSION).zip HEAD

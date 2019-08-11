PACKAGE_NAME ?= confluent-cli
VERSION ?= 5.4.0-beta190811193210
PLATFORM = $(shell uname -s)
INSTALL_FLAGS = -D
ifeq ($(PLATFORM),Linux)
	INSTALL_FLAGS = -D
endif
ifeq ($(PLATFORM),Darwin)
	INSTALL_FLAGS =
endif

.PHONY:

all: build

install: build
ifndef CONFLUENT_HOME
	$(error Cannot install. CONFLUENT_HOME is not set)
endif
	install $(INSTALL_FLAGS) -m 755 bin/confluent $(CONFLUENT_HOME)/bin/confluent

build: ccs

prep:
	mkdir -p bin/

ccs: prep
	cp -f src/ccs/confluent.sh bin/confluent
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

# Dependencies you'll probably need to install to compile this: make, curl, git,
# zip, unzip, patch, java7-jdk | openjdk-7-jdk, maven.

# Release specifics. Note that some of these (VERSION, DESTDIR)
# are required and passed to create_archive.sh as environment variables. That
# script can also pick up some other settings (PREFIX, SYSCONFDIR) to customize
# layout of the installation.
ifndef VERSION
# Note that this is sensitive to this package's version being the first
# <version> tag in the pom.xml
VERSION = 0.0.1-SNAPSHOT
endif

export PACKAGE_TITLE ?= confluent-cli
export FULL_PACKAGE_TITLE = $(PACKAGE_TITLE)
export PACKAGE_NAME = $(FULL_PACKAGE_TITLE)-$(VERSION)

# Defaults that are likely to vary by platform. These are cleanly separated so
# it should be easy to maintain altered values on platform-specific branches
# when the values aren't overridden by the script invoking the Makefile
DEFAULT_DESTDIR = $(CURDIR)/BUILD
DEFAULT_PREFIX = /usr
DEFAULT_SYSCONFDIR = /etc/$(PACKAGE_TITLE)

# Install directories
ifndef DESTDIR
DESTDIR = $(DEFAULT_DESTDIR)
endif
# For platform-specific packaging you'll want to override this to a normal
# PREFIX like /usr or /usr/local. Using the PACKAGE_NAME here makes the default
# zip/tgz files use a format like:
#   kafka-version-scalaversion/
#     bin/
#     etc/
#     share/kafka/
ifndef PREFIX
PREFIX = $(DEFAULT_PREFIX)
endif

ifndef SYSCONFDIR
SYSCONFDIR:= $(DEFAULT_SYSCONFDIR)
endif
SYSCONFDIR:=$(subst PREFIX,$(PREFIX),$(SYSCONFDIR))

export VERSION
export DESTDIR
export PREFIX
export SYSCONFDIR

# For this makefile to work, packaging needs first to merge the main code branch
# within this rpm branch.
CONFLUENT_HOME = $(DESTDIR)$(PREFIX)
include Common.mk

all: rpm

export RPM_VERSION=$(shell echo $(VERSION) | sed -e 's/-alpha[0-9]*//' -e 's/-beta[0-9]*//' -e 's/-rc[0-9]*//' -e 's/-SNAPSHOT//')
# Get any -alpha, -beta, -rc piece that we need to put into the Release part of
# the version since RPM versions don't support non-numeric
# characters. Ultimately, for something like 0.8.2-beta, we want to end up with
# Version=0.8.2 Release=0.X.beta
# where X is the RPM release # of 0.8.2-beta (the prefix 0. forces this to be
# considered earlier than any 0.8.2 final releases since those will start with
# Version=0.8.2 Release=1)
export RPM_RELEASE_POSTFIX=$(subst -,,$(subst $(RPM_VERSION),,$(VERSION)))
ifneq ($(RPM_RELEASE_POSTFIX),)
	export RPM_RELEASE_POSTFIX_UNDERSCORE=_$(RPM_RELEASE_POSTFIX)
endif

rpm: RPM_BUILDING/SOURCES/$(FULL_PACKAGE_TITLE)-$(RPM_VERSION).tar.gz
	echo "Building the rpm"
	RPMOPTS='define "_topdir `pwd`/RPM_BUILDING"'
	export RPMOPTS
	rpmbuild --define="_topdir `pwd`/RPM_BUILDING" -tb $<
	find RPM_BUILDING/{,S}RPMS/ -type f | xargs -n1 -IXXX mv XXX .
	echo
	echo "================================================="
	echo "The rpms have been created and can be found here:"
	@ls -laF $(FULL_PACKAGE_TITLE)*rpm
	echo "================================================="

# Unfortunately, because of version naming issues and the way rpmbuild expects
# the paths in the tar file to be named, we need to rearchive the package. So
# instead of depending on archive, this target just uses the unarchived,
# installed version to generate a new archive. Note that we always regenerate
# the symlink because the RPM_VERSION doesn't include all the version info -- it
# can leave of things like -beta, -rc1, etc.
RPM_BUILDING/SOURCES/$(FULL_PACKAGE_TITLE)-$(RPM_VERSION).tar.gz: rpm-build-area install $(FULL_PACKAGE_TITLE).spec.in RELEASE_$(RPM_VERSION)$(RPM_RELEASE_POSTFIX_UNDERSCORE)
	rm -rf $(FULL_PACKAGE_TITLE)-$(RPM_VERSION)
	mkdir $(FULL_PACKAGE_TITLE)-$(RPM_VERSION)
	cp -R $(DESTDIR)/* $(FULL_PACKAGE_TITLE)-$(RPM_VERSION)
	./create_spec.sh $(FULL_PACKAGE_TITLE).spec.in $(FULL_PACKAGE_TITLE)-$(RPM_VERSION)/$(FULL_PACKAGE_TITLE).spec
	rm -f $@ && tar -czf $@ $(FULL_PACKAGE_TITLE)-$(RPM_VERSION)
	rm -rf $(FULL_PACKAGE_TITLE)-$(RPM_VERSION)

rpm-build-area: RPM_BUILDING/BUILD RPM_BUILDING/RPMS RPM_BUILDING/SOURCES RPM_BUILDING/SPECS RPM_BUILDING/SRPMS

RPM_BUILDING/%:
	mkdir -p $@

$(FULL_PACKAGE_TITLE).spec:
	./create_spec.sh $(FULL_PACKAGE_TITLE).spec.in $(FULL_PACKAGE_TITLE).spec

RELEASE_%:
	echo 0 > $@

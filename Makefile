# Dependencies you'll probably need to install to compile this: make, curl, git,
# zip, unzip, patch, java7-jdk | openjdk-7-jdk, maven.

# Release specifics. Note that some of these (VERSION, DESTDIR)
# are required and passed to create_archive.sh as environment variables. That
# script can also pick up some other settings (PREFIX, SYSCONFDIR) to customize
# layout of the installation.
ifndef VERSION
# Note that this is sensitive to this package's version being the first
# <version> tag in the pom.xml
VERSION = SNAPSHOT
endif

export PACKAGE_TITLE ?=  confluent-cli
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
CONFLUENT_HOME = $(DESTDIR)
include Common.mk

all: install

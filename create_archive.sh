#!/usr/bin/env bash
#
# Creates an archive suitable for distribution (standard layout for binaries,
# libraries, etc.).

set -e

if [ -z ${PACKAGE_TITLE} -o -z ${DESTDIR} ]; then
    echo "PACKAGE_TITLE and DESTDIR environment variables must be set."
    exit 1
fi

BINPATH=${PREFIX}/bin

INSTALL="install -D -m 644"
INSTALL_X="install -D -m 755"

rm -rf ${DESTDIR}${PREFIX}
mkdir -p ${DESTDIR}${PREFIX}
mkdir -p ${DESTDIR}${BINPATH}

function copy_package() {
    find bin/ -type f | xargs -I XXX ${INSTALL_X} -o root -g root XXX ${DESTDIR}${PREFIX}/XXX
}

case "${PACKAGE_TITLE}" in
  "confluent-cli")
    copy_package ${PACKAGE_TITLE}
    ;;
  *)
    echo "Unexpected value for PACKAGE_TITLE environment variable found: ${PACKAGE_TITLE}. 
    Expected confluent-cli."
    exit 1
    ;;
esac

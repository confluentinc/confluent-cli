#!/usr/bin/env bash

# Copyright 2017 Confluent Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

success=false
set -x

command_name="$( basename ${BASH_SOURCE[0]} )"

confluent_bin="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

confluent_home="$( dirname "${confluent_bin}" )"

confluent_conf="${confluent_home}/etc"

# $TMPDIR includes a trailing '/' by default.
tmp_dir="${TMPDIR:-/tmp/}"

# Contains the result of functions that intend to return a value.
_retval=""

declare -a services=(
    "zookeeper"
    "kafka"
)

echo_variable() {
    local var_value="${!1}"
    echo "${1} = ${var_value}"
}

set_or_get_current() {
    if [[ -f "${tmp_dir}confluent.current" ]]; then
        export confluent_current="$( cat "${tmp_dir}confluent.current" )"
    else
        export confluent_current="$( mktemp -d -t confluent )"
        echo ${confluent_current} > "${tmp_dir}confluent.current"
    fi
}

shutdown() {
    [[ ${confluent_current} == ${tmp_dir}confluent* && ${success} == false ]] \
        && echo "Removing: ${confluent_current}" \
        && rm -rf ${confluent_current}
}

is_alive() {
    local pid="${1}"
    kill -0 "${pid}" > /dev/null 2>&1
}

wait_process() {
    local pid="${1}"
    local timeout_ms="${2}"
    # Default max wait time set to 10 minutes. That's practically infinite for this program.
    [[ -n "${timeout_ms}" ]] || timeout_ms=600000

    while is_alive "${pid}" && [[ "${timeout_ms}" -gt 0 ]]; do
        sleep 0.5
        echo "Waiting: ${timeout_ms}"
        (( timeout_ms = timeout_ms - 500 ))
    done
    is_alive "${pid}"
}

stop_and_wait_process() {
    local pid="${1}"
    local timeout_ms="${2}"
    # Default max wait time set to 10 minutes. That's practically infinite for this program.
    [[ -n "${timeout_ms}" ]] || timeout_ms=600000

    kill "${pid}"
    while kill -0 "${pid}" > /dev/null 2>&1 && [[ "${timeout_ms}" -gt 0 ]]; do
        sleep 0.5
        echo "Waiting: ${timeout_ms}"
        (( timeout_ms = timeout_ms - 500 ))
    done
    kill -9 "${pid}" > /dev/null 2>&1
}

start_zookeeper() {
    local service="zookeeper"
    echo "Starting ${service}"
    local service_dir="${confluent_current}/${service}"
    mkdir -p ${service_dir}
    config_${service}
    ${confluent_bin}/zookeeper-server-start "${service_dir}/${service}.properties" \
        2> "${service_dir}/${service}.stderr" \
        1> "${service_dir}/${service}.stdout" &
    echo $! > "${service_dir}/${service}.pid"
    wait_${service} "$( cat ${service_dir}/${service}.pid )"
    #stop_and_wait_process "$( cat ${service_dir}/${service}.pid )" 10000
}

config_zookeeper() {
    local service="zookeeper"
    echo "Configuring ${service}"
    local service_dir="${confluent_current}/${service}"
    mkdir -p "${service_dir}/data"
    sed "s@^dataDir=.*@dataDir=${service_dir}/data@g" \
        < "${confluent_conf}/kafka/${service}.properties" \
        > "${service_dir}/${service}.properties"
}

stop_zookeeper() {
    local service="zookeeper"
    local service_dir="${confluent_current}/${service}"
    local service_pid="$(cat ${service_dir}/${service}.pid )"
    echo "Stopping ${service}"

    kill ${service_pid}
    rm -f ${service_dir}/${service}.pid
}

wait_zookeeper() {
    local pid="${1}"
    zk_port=$( grep "clientPort" "${confluent_conf}/kafka/${service}.properties" \
        | cut -f 2 -d '=' \
        | xargs )

    echo ${zk_port}
    started=false
    while ${started} == false; do
        ( lsof -P -c java | grep ${zk_port} ) && started=true
        [[ ${started} == false ]] && sleep 0.5
    done
    wait_process ${pid} 1000 || echo "Zookeeper failed to start"
}

start_kafka() {
    local service="kafka"
    echo "Starting ${service}"
}

stop_kafka() {
    local service="kafka"
    echo "Stopping ${service}"
}

usage() {
    cat <<EOF
${command_name}: a command line interface to manage Confluent services

Usage: ${command_name} [<options>] <command> [<subcommand>] [<parameters>]

start a service:
    start   start_doc

stop a service:
    stop   stop_doc

'${command_name} help' list available commands See 'git help <command>' to read about a
specific subcommand.

EOF
    exit 0
}

set_or_get_current

echo_variable tmp_dir
echo_variable confluent_home
echo_variable confluent_bin
echo_variable confluent_conf
echo_variable confluent_current

# Parse command-line arguments
[[ $# -lt 1 ]] && usage
command="${1}"
shift
case "${command}" in
    help) usage;;

    start|stop)
        for service in ${services}; do
            ${command}_${service} "${@}";
        done;;
    *)  echo "Unknown command '${command}'.  Type '${command_name} help' for a list of available
    commands."
        exit 1;;
esac

echo "Hello World! I'm Confluent Platform OSS CLI!"

success=true

trap shutdown EXIT

echo "Goodbye!"

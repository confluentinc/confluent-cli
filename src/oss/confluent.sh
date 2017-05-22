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
# Uncomment to enable debugging on the console
#set -x

command_name="$( basename ${BASH_SOURCE[0]} )"

confluent_bin="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

confluent_home="$( dirname "${confluent_bin}" )"

confluent_conf="${confluent_home}/etc"

# $TMPDIR includes a trailing '/' by default.
tmp_dir="${TMPDIR:-/tmp/}"

# Contains the result of functions that intend to return a value besides their exit status.
_retval=""

declare -a services=(
    "zookeeper"
    "kafka"
)

declare -a rev_services=(
    "kafka"
    "zookeeper"
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

is_running() {
    local service="${1}"
    local service_dir="${confluent_current}/${service}"
    local service_pid="$(cat ${service_dir}/${service}.pid )"

    is_alive ${service_pid}
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

    kill "${pid}" 2> /dev/null
    while kill -0 "${pid}" > /dev/null 2>&1 && [[ "${timeout_ms}" -gt 0 ]]; do
        sleep 0.5
        echo "Waiting: ${timeout_ms}"
        (( timeout_ms = timeout_ms - 500 ))
    done
    # Will have no effect if the process stopped gracefully
    kill -9 "${pid}" > /dev/null 2>&1
}

start_zookeeper_old() {
    local service="zookeeper"
    echo "Starting ${service}"
    local service_dir="${confluent_current}/${service}"
    mkdir -p ${service_dir}
    config_${service}
    # TODO: decide whether to persist logs on stdout / stderr between runs.
    ${confluent_bin}/zookeeper-server-start "${service_dir}/${service}.properties" \
        2> "${service_dir}/${service}.stderr" \
        1> "${service_dir}/${service}.stdout" &
    echo $! > "${service_dir}/${service}.pid"
    sleep 3
    stop_and_wait_process "$( cat ${service_dir}/${service}.pid )" 10000
    sleep 1
    wait_${service} "$( cat ${service_dir}/${service}.pid )"
}

start_zookeeper() {
    start_service "zookeeper" "${confluent_bin}/zookeeper-server-start"
}

config_zookeeper_old() {
    local service="zookeeper"
    echo "Configuring ${service}"
    local service_dir="${confluent_current}/${service}"
    mkdir -p "${service_dir}/data"
    sed "s@^dataDir=.*@dataDir=${service_dir}/data@g" \
        < "${confluent_conf}/kafka/${service}.properties" \
        > "${service_dir}/${service}.properties"
}

config_zookeeper() {
    config_service "zookeeper" "kafka" "zookeeper" "dataDir"
}

stop_zookeeper() {
    stop_service "zookeeper"
}

wait_zookeeper() {
    local pid="${1}"
    local zk_port=$( grep "clientPort" "${confluent_conf}/kafka/${service}.properties" \
        | cut -f 2 -d '=' \
        | xargs )

    echo ${zk_port}
    local started=false
    local timeout_ms=1000
    while ${started} == false && [[ "${timeout_ms}" -gt 0 ]]; do
        ( lsof -P -c java | grep ${zk_port} ) && started=true
        [[ ${started} == false ]] && sleep 0.5 && (( timeout_ms = timeout_ms - 500 ))
    done
    wait_process ${pid} 1000 || echo "Zookeeper failed to start"
}

start_kafka() {
    local service="kafka"
    echo "Starting ${service}"
    is_running "zookeeper" || ( echo "Cannot start Kafka, Zookeeper is not running. Check your deployment" && exit 1 )

    start_service "kafka" "${confluent_bin}/kafka-server-start"
}

config_kafka() {
    config_service "kafka" "kafka" "server" "logs.dir"
}

stop_kafka() {
    stop_service "kafka"
}

wait_kafka() {
    local pid="${1}"
    local kafka_port=9092

    echo ${kafka_port}
    local started=false
    local timeout_ms=5000
    while ${started} == false && [[ "${timeout_ms}" -gt 0 ]]; do
        ( lsof -P -c java | grep ${zk_port} ) && started=true
        [[ ${started} == false ]] && sleep 0.5 && (( timeout_ms = timeout_ms - 500 ))
    done
    wait_process ${pid} 5000 || echo "Kafka failed to start"
}

start_service() {
    local service="${1}"
    local start_command="${2}"
    echo "Starting ${service}"
    local service_dir="${confluent_current}/${service}"
    mkdir -p ${service_dir}
    config_${service}
    # TODO: decide whether to persist logs on stdout / stderr between runs.
    ${start_command} "${service_dir}/${service}.properties" \
        2> "${service_dir}/${service}.stderr" \
        1> "${service_dir}/${service}.stdout" &
    echo $! > "${service_dir}/${service}.pid"
    wait_${service} "$( cat ${service_dir}/${service}.pid )"
}

# The first 3 args seem unavoidable right now. 4th is optional
config_service() {
    local service="${1}"
    local package="${2}"
    local property_file="${3}"
    [[ -z "${service}" ]] || [[ -z "${package}" ]] || [[ -z "${package}" ]]
    echo "Configuring ${service}"
    local service_dir="${confluent_current}/${service}"
    mkdir -p "${service_dir}/data"
    local property="${4}"
    if [[ -z "${property}" ]]; then
        config_command="sed \"s@^${property}=.*@${property}=${service_dir}/data@g\""
    else
        config_command=cat
    fi

    ${config_command} < "${confluent_conf}/${package}/${property_file}.properties" \
        > "${service_dir}/${service}.properties"
}

stop_service() {
    local service="${1}"
    local service_dir="${confluent_current}/${service}"
    # check file exists, and if not issue warning.
    local service_pid="$(cat ${service_dir}/${service}.pid )"
    echo "Stopping ${service}"

    stop_and_wait_process ${service_pid} 5000
    rm -f ${service_dir}/${service}.pid
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

#echo_variable tmp_dir
#echo_variable confluent_home
#echo_variable confluent_bin
#echo_variable confluent_conf
#echo_variable confluent_current

# Parse command-line arguments
[[ $# -lt 1 ]] && usage
command="${1}"
shift
case "${command}" in
    help) usage;;

    start)
        for service in "${services[@]}"; do
            ${command}_${service} "${@}";
        done;;

    stop)
        for service in "${rev_services[@]}"; do
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

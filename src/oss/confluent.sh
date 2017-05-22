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

# Exit with an error message.
die() {
    echo $@
    exit 1
}

# Implies zero or positive
is_integer() {
    [[ -n "${1}" ]] && [[ "${1}" =~ ^[0-9]+$ ]]
}

wheel_pos=0
wheel_freq_ms=100
spinner() {
    local wheel='-\|/'
    wheel_pos=$(( (wheel_pos + 1) % ${#wheel} ))
    printf "\r${wheel:${wheel_pos}:1}"
    sleep 0.${wheel_freq_ms}
}

set_or_get_current() {
    if [[ -f "${tmp_dir}confluent.current" ]]; then
        export confluent_current="$( cat "${tmp_dir}confluent.current" )"
    fi

    if [[ ! -d "${confluent_current}" ]]; then
        export confluent_current="$( mktemp -d -t confluent )"
        echo ${confluent_current} > "${tmp_dir}confluent.current"
    fi
}

shutdown() {
    [[ ${success} == false ]] && echo "Shutting down. Whatever that means."
}

destroy() {
    [[ ${confluent_current} == ${tmp_dir}confluent* ]] \
        && echo "Removing: ${confluent_current}" \
        && rm -rf ${confluent_current}
}

is_alive() {
    local pid="${1}"
    kill -0 "${pid}" > /dev/null 2>&1
}

is_not_alive() {
    ! is_alive "${1}"
}

is_running() {
    local service="${1}"
    local service_dir="${confluent_current}/${service}"
    local service_pid="$(cat ${service_dir}/${service}.pid )"

    is_alive ${service_pid}
}

wait_process_up() {
    wait_process "${1}" "up" "${2}"
}

wait_process_down() {
    wait_process "${1}" "down" "${2}"
}

wait_process() {
    local pid="${1}"
    local event="${2}"
    local timeout_ms="${3}"

    is_integer "${pid}" || die "Need PID to wait on a process"

    # By default wait for a process to start
    [[ -n "${event}" ]] || event="up"

    # Default max wait time set to 10 minutes. That's practically infinite for this program.
    is_integer "${timeout_ms}" || timeout_ms=600000

    local mode=is_alive
    if [[ "${event}" == "down" ]]; then
        mode=is_not_alive
    fi

    while ${mode} "${pid}" && [[ "${timeout_ms}" -gt 0 ]]; do
        #echo "Waiting: ${timeout_ms}"
        spinner
        (( timeout_ms = timeout_ms - ${wheel_freq_ms} ))
    done
    # New line after spinner
    echo ""
    ${mode} "${pid}"
}

stop_and_wait_process() {
    local pid="${1}"
    local timeout_ms="${2}"
    # Default max wait time set to 10 minutes. That's practically infinite for this program.
    [[ -n "${timeout_ms}" ]] || timeout_ms=600000

    kill "${pid}" 2> /dev/null
    while kill -0 "${pid}" > /dev/null 2>&1 && [[ "${timeout_ms}" -gt 0 ]]; do
        #echo "Waiting: ${timeout_ms}"
        spinner
        (( timeout_ms = timeout_ms - ${wheel_freq_ms} ))
    done
    # New line after spinner
    echo ""
    # Will have no effect if the process stopped gracefully
    kill -9 "${pid}" > /dev/null 2>&1
}

start_zookeeper() {
    start_service "zookeeper" "${confluent_bin}/zookeeper-server-start"
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

    local started=false
    local timeout_ms=1000
    while ${started} == false && [[ "${timeout_ms}" -gt 0 ]]; do
        ( lsof -P -c java | grep ${zk_port} ) && started=true
        [[ ${started} == false ]] && spinner && (( timeout_ms = timeout_ms - ${wheel_freq_ms} ))
    done
    wait_process_up ${pid} 1000 || echo "Zookeeper failed to start"
}

start_kafka() {
    local service="kafka"
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

    local started=false
    local timeout_ms=5000
    while ${started} == false && [[ "${timeout_ms}" -gt 0 ]]; do
        ( lsof -P -c java | grep ${zk_port} ) && started=true
        [[ ${started} == false ]] && spinner && (( timeout_ms = timeout_ms - ${wheel_freq_ms} ))
    done
    wait_process_up ${pid} 5000 || echo "Kafka failed to start"
}

start_service() {
    local service="${1}"
    local start_command="${2}"
    local service_dir="${confluent_current}/${service}"
    mkdir -p ${service_dir}
    config_${service}
    echo "Starting ${service}"
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
    #echo "Configuring ${service}"
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

service_exists() {
    local arg="${1}"
    for service in "${services[@]}"; do
        [[ ${service} == "${arg}" ]] && return 0;
    done
    return 1
}

start_subcommand() {
    local subcommand="${1}"

    [[ -n "${subcommand}" ]] && ! service_exists "${subcommand}" && die "Unknown service: ${subcommand}"

    for service in "${services[@]}"; do
        start_${service} "${@}";
        [[ "${service}" == "${subcommand}" ]] && break;
    done
}

stop_subcommand() {
    local subcommand="${1}"

    [[ -n "${subcommand}" ]] && ! service_exists "${subcommand}" && die "Unknown service: ${subcommand}"

    skip=true
    [[ -z "${subcommand}" ]] && skip=false
    for service in "${rev_services[@]}"; do
        [[ "${service}" == "${subcommand}" ]] && skip=false;
        [[ "${skip}" == false ]] && stop_${service} "${@}";
    done
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
        start_subcommand $*;;

    stop)
        stop_subcommand $*;;

    destroy)
        destroy;;

    *)  echo "Unknown command '${command}'.  Type '${command_name} help' for a list of available
    commands."
        exit 1;;
esac

echo "Hello World! I'm Confluent Platform OSS CLI!"

success=true

trap shutdown EXIT

echo "Goodbye!"

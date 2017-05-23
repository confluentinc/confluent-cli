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

Gre='\e[0;32m';
Red='\e[0;31m';
#Reset color
RC='\e[0m'

declare -a services=(
    "zookeeper"
    "kafka"
    "schema-registry"
    "kafka-rest"
    "connect"
)

declare -a rev_services=(
    "connect"
    "kafka-rest"
    "schema-registry"
    "kafka"
    "zookeeper"
)

declare -a commands=(
    "start"
    "stop"
    "status"
    "current"
    "destroy"
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

# Will pick the last integer value that will encounter
get_service_port() {
    local property=${1}
    local config_file=${2}
    local delim=${3}

    [[ -z "${property}" ]] && die "Need property key to extract service port from configuration"
    [[ -z "${delim}" ]] && delim=":"
    local property_split=( $( grep -i "^${property}" "${config_file}" | tr "${delim}" "\n" ) )

    _retval=""
    for entry in "${property_split[@]}"; do
        # trim string
        entry=$( echo ${entry} | xargs )
        is_integer "${entry}" && _retval="${entry}"
    done
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
    [[ "${success}" == false ]] \
        && echo "Unsuccessful execution. Attempting service shutdown" \
        && stop_command

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
    local print_status="${2}"
    local service_dir="${confluent_current}/${service}"
    local service_pid="$( cat ${service_dir}/${service}.pid 2> /dev/null )"

    is_alive ${service_pid}
    status=$?
    [[ "${print_status}" != false ]] \
        && ( [[ ${status} -eq 0 ]] \
            && printf "${service} is [${Gre}UP${RC}]\n" \
            || printf "${service} is [${Red}DOWN${RC}]\n")

    return ${status}
}

wait_process_up() {
    #TODO: need to add port condition here too probably
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

    local mode=is_not_alive
    if [[ "${event}" == "down" ]]; then
        mode=is_alive
    fi

    while ${mode} "${pid}" && [[ "${timeout_ms}" -gt 0 ]]; do
        spinner
        (( timeout_ms = timeout_ms - ${wheel_freq_ms} ))
    done
    # Backspace to override spinner in the next printf/echo
    printf "\b"
    ! ${mode} "${pid}"
}

stop_and_wait_process() {
    local pid="${1}"
    local timeout_ms="${2}"
    # Default max wait time set to 10 minutes. That's practically infinite for this program.
    [[ -n "${timeout_ms}" ]] || timeout_ms=600000

    kill "${pid}" 2> /dev/null
    while kill -0 "${pid}" > /dev/null 2>&1 && [[ "${timeout_ms}" -gt 0 ]]; do
        spinner
        (( timeout_ms = timeout_ms - ${wheel_freq_ms} ))
    done
    # Backspace to override spinner in the next printf/echo
    printf "\b"
    # Will have no effect if the process stopped gracefully
    # TODO: maybe should issue a warning if the process is not stopped gracefully.
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
    get_service_port "clientPort" "${confluent_conf}/kafka/zookeeper.properties" "="
    export zk_port="${_retval}"

    if [[ -n "${_retval}" ]]; then
        export zk_port="${_retval}"
    else
        export zk_port="2181"
    fi

    local started=false
    local timeout_ms=5000
    while [[ "${started}" == false && "${timeout_ms}" -gt 0 ]]; do
        ( lsof -P -c java | grep ${zk_port} > /dev/null 2>&1 ) && started=true
        spinner && (( timeout_ms = timeout_ms - ${wheel_freq_ms} ))
    done
    wait_process_up ${pid} 5000 || echo "Zookeeper failed to start"
}

start_kafka() {
    local service="kafka"
    is_running "zookeeper" "false" \
        || die "Cannot start Kafka, Zookeeper is not running. Check your deployment"
    start_service "kafka" "${confluent_bin}/kafka-server-start"
}

config_kafka() {
    config_service "kafka" "kafka" "server" "log.dirs"
}

stop_kafka() {
    stop_service "kafka"
}

wait_kafka() {
    local pid="${1}"
    get_service_port "listeners" "${confluent_conf}/kafka/server.properties"
    if [[ -n "${_retval}" ]]; then
        export kafka_port="${_retval}"
    else
        export kafka_port="9092"
    fi

    local started=false
    local timeout_ms=10000

    while [[ "${started}" == false && "${timeout_ms}" -gt 0 ]]; do
        ( lsof -P -c java | grep ${kafka_port} > /dev/null 2>&1 ) && started=true
        spinner && (( timeout_ms = timeout_ms - ${wheel_freq_ms} ))
    done
    wait_process_up ${pid} 5000 || echo "Kafka failed to start"
}

start_schema-registry() {
    local service="schema-registry"
    is_running "kafka" "false" \
        || die "Cannot start Schema Registry, Kafka Server is not running. Check your deployment"
    start_service "schema-registry" "${confluent_bin}/schema-registry-start"
}

config_schema-registry() {
    config_service "schema-registry" "schema-registry" "schema-registry"\
        "kafkastore.connection.url" "localhost:${zk_port}"
}

stop_schema-registry() {
    stop_service "schema-registry"
}

wait_schema-registry() {
    local pid="${1}"

    get_service_port "listeners" "${confluent_conf}/schema-registry/schema-registry.properties"
    if [[ -n "${_retval}" ]]; then
        export schema_registry_port="${_retval}"
    else
        export schema_registry_port="8081"
    fi

    local started=false
    local timeout_ms=10000
    while [[ "${started}" == false && "${timeout_ms}" -gt 0 ]]; do
        ( lsof -P -c java | grep ${schema_registry_port} > /dev/null 2>&1 ) && started=true
        spinner && (( timeout_ms = timeout_ms - ${wheel_freq_ms} ))
    done
    wait_process_up ${pid} 5000 || echo "Schema Registry failed to start"
}

start_kafka-rest() {
    local service="kafka-rest"
    is_running "kafka" "false" \
        || die "Cannot start Kafka Rest, Kafka Server is not running. Check your deployment"
    start_service "kafka-rest" "${confluent_bin}/kafka-rest-start"
}

config_kafka-rest() {
    config_service "kafka-rest" "kafka-rest" "kafka-rest"\
        "zookeeper.connect" "localhost:${zk_port}"

    config_service "kafka-rest" "kafka-rest" "kafka-rest"\
        "schema.registry.url" "http://localhost:${schema_registry_port}" "reapply"
}

stop_kafka-rest() {
    stop_service "kafka-rest"
}

wait_kafka-rest() {
    local pid="${1}"

    get_service_port "listeners" "${confluent_conf}/kafka-rest/kafka-rest.properties"
    if [[ -n "${_retval}" ]]; then
        export kafka_rest_port="${_retval}"
    else
        export kafka_rest_port="8082"
    fi

    local started=false
    local timeout_ms=10000
    while [[ "${started}" == false && "${timeout_ms}" -gt 0 ]]; do
        ( lsof -P -c java | grep ${kafka_rest_port} > /dev/null 2>&1 ) && started=true
        spinner && (( timeout_ms = timeout_ms - ${wheel_freq_ms} ))
    done
    wait_process_up ${pid} 5000 || echo "Kafka Rest failed to start"
}

start_connect() {
    local service="connect"
    is_running "kafka" "false" \
        || die "Cannot start Kafka Connect, Kafka Server is not running. Check your deployment"
    start_service "connect" "${confluent_bin}/connect-distributed"
}

config_connect() {
    config_service "connect" "kafka" "connect-distributed" "bootstrap.servers"\
        "localhost:${kafka_port}"
}

stop_connect() {
    stop_service "connect"
}

wait_connect() {
    local pid="${1}"

    get_service_port "rest.port" "${confluent_conf}/kafka/connect-distributed.properties" "="
    if [[ -n "${_retval}" ]]; then
        export connect_port="${_retval}"
    else
        export connect_port="8083"
    fi

    local started=false
    local timeout_ms=20000
    while [[ "${started}" == false && "${timeout_ms}" -gt 0 ]]; do
        ( lsof -P -c java | grep ${connect_port} > /dev/null 2>&1 ) && started=true
        spinner && (( timeout_ms = timeout_ms - ${wheel_freq_ms} ))
    done
    wait_process_up ${pid} 5000 || echo "Kafka Connect failed to start"
}

status_service() {
    local service="${1}"

    [[ -n "${service}" ]] \
        && ! service_exists "${service}" && die "Unknown service: ${service}"

    skip=true
    [[ -z "${service}" ]] && skip=false
    for entry in "${rev_services[@]}"; do
        [[ "${entry}" == "${service}" ]] && skip=false;
        [[ "${skip}" == false ]] && is_running "${entry}"
    done
}

start_service() {
    local service="${1}"
    local start_command="${2}"
    local service_dir="${confluent_current}/${service}"
    is_running "${service}" "false" \
        && echo "${service} is already running. Try restarting if needed"\
        && return 0

    mkdir -p ${service_dir}
    config_${service}
    echo "Starting ${service}"
    # TODO: decide whether to persist logs on stdout / stderr between runs.
    ${start_command} "${service_dir}/${service}.properties" \
        2> "${service_dir}/${service}.stderr" \
        1> "${service_dir}/${service}.stdout" &
    echo $! > "${service_dir}/${service}.pid"
    local service_pid="$( cat ${service_dir}/${service}.pid 2> /dev/null )"
    wait_${service} "${service_pid}"
    is_running "${service}"
}

# The first 3 args seem unavoidable right now. 4th is optional
config_service() {
    local service="${1}"
    local package="${2}"
    local property_file="${3}"

    ( [[ -z "${service}" ]] || [[ -z "${package}" ]] || [[ -z "${package}" ]] ) \
        && die "Missing required configuration properties for service: ${service}"

    #echo "Configuring ${service}"
    local service_dir="${confluent_current}/${service}"
    mkdir -p "${service_dir}/data"
    local property_key="${4}"
    local property_value="${5}"
    if [[ -n "${property_key}" && -z "${property_value}" ]]; then
        config_command="sed -e s@^${property_key}=.*@${property_key}=${service_dir}/data@g"
    elif [[ -n "${property_key}" && -n "${property_value}" ]]; then
        config_command="sed -e s@^${property_key}=.*@${property_key}=${property_value}@g"
    else
        config_command=cat
    fi

    local input_file="${confluent_conf}/${package}/${property_file}.properties"
    local reaplly="${6}"
    [[ -n "${reapply}" ]] && input_file="${service_dir}/${service}.properties"

    ${config_command} < "${input_file}" \
        > "${service_dir}/${service}.properties"
}

stop_service() {
    local service="${1}"
    local service_dir="${confluent_current}/${service}"
    # check file exists, and if not issue warning.
    local service_pid="$( cat ${service_dir}/${service}.pid 2> /dev/null )"
    echo "Stopping ${service}"

    stop_and_wait_process ${service_pid} 10000
    rm -f ${service_dir}/${service}.pid
    is_running "${service}"
}

service_exists() {
    local service="${1}"
    exists "${service}" services
}

command_exists() {
    local command="${1}"
    exists "${command}" commands
}

exists() {
    local arg="${1}"
    local -n list="${2}"

    for entry in "${list[@]}"; do
        [[ ${entry} == "${arg}" ]] && return 0;
    done
    return 1
}

list_command() {
    for service in "${services[@]}"; do
        echo "${service}"
    done
}

start_command() {
    start_or_stop_service "start" services "${@}"
}


stop_command() {
    start_or_stop_service "stop" rev_services "${@}"
    return 0
}

status_command() {
    set_or_get_current
    status_service "${@}"
}

start_or_stop_service() {
    set_or_get_current
    local command="${1}"
    shift
    local -n list="${1}"
    shift
    local service="${1}"
    shift

    [[ -n "${service}" ]] \
        && ! service_exists "${service}" && die "Unknown service: ${service}"

    for entry in "${list[@]}"; do
        ${command}_${entry} "${@}";
        [[ "${entry}" == "${service}" ]] && break;
    done
}

print_current() {
    set_or_get_current
    echo "${confluent_current}"
}

destroy() {
    if [[ -f "${tmp_dir}confluent.current" ]]; then
        export confluent_current="$( cat "${tmp_dir}confluent.current" )"
    fi

    [[ ${confluent_current} == ${tmp_dir}confluent* ]] \
        && stop_command \
        && echo "Deleting: ${confluent_current}" \
        && rm -rf ${confluent_current} \
        && rm -f "${tmp_dir}confluent.current"
}

list_usage() {
    cat <<EOF
Usage: ${command_name} list

Description:
    A list of all the available services.

EOF
    exit 0
}

start_usage() {
    cat <<EOF
Usage: ${command_name} start [<service>]

Description:
    Start all services. If a specific <service> is given as an argument it starts this service
    along with all of its dependencies.

Output:
    Print a status messages after starting each service to indicate successful startup or an error.

Examples:
    confluent start
        Starts all available services.

    confluent start kafka
        Starts kafka and zookeeper as its dependency.

EOF
    exit 0
}

stop_usage() {
    cat <<EOF
Usage: ${command_name} stop [<service>]

Description:
    Stop all services. If a specific <service> is given as an argument it stops this service
    along with all of its dependencies.

Output:
    Print a status messages after stopping each service to indicate successful shutdown or an error.

Examples:
    confluent stop
        Stops all available services.

    confluent stop kafka
        Stops kafka and zookeeper as its dependency.

EOF
    exit 0
}

status_usage() {
    cat <<EOF
Usage: ${command_name} status [<service>]

Description:
    Return the status of all services. If a specific <service> is given as an argument the status of
    the requested service is returned along with the status of its dependencies.

Output:
    Print a status messages for each service.

Examples:
    confluent status
        Print the status of all the available services.

    confluent status kafka
        Prints the status of kafka and the status of zookeeper.

EOF
    exit 0
}

current_usage() {
    cat <<EOF
Usage: ${command_name} current

Description:
    Return the filesystem path of the data and logs of the services managed by the current
    confluent run. If such a path does not exist, it will be created.


Output:
    The filesystem directory path.


Examples:
    confluent current
        /tmp/confluent.SpBP4fQi
EOF
    exit 0
}

destroy_usage() {
    cat <<EOF
Usage: ${command_name} destroy

Description:
    Unpersist an existing confluent run. Any running services are stopped. The data and the log
    files of all services will be deleted.


Examples:
    confluent destroy
        Print the status of all the available services.
EOF
    exit 0
}

usage() {
    cat <<EOF
${command_name}: a command line interface to manage Confluent services

Usage: ${command_name} [<options>] <command> [<subcommand>] [<parameters>]

    list        List available services.

    start       Start all services or a service along with its dependencies

    stop        Stop all services or a service along with the services depending on it.

    status      Get the status of all services or a specific service along with its dependencies.

    current     Get the path of the data and logs of the services managed by the current confluent run.

    destroy     Delete the data and logs of the current confluent run.

'${command_name} help' lists available commands See 'git help <command>' to read about a
specific command.

EOF
    exit 0
}

invalid_command() {
    echo "Unknown command '${1}'."
    echo "Type '${command_name} help' for a list of available commands."
    exit 1
}

# Parse command-line arguments
[[ $# -lt 1 ]] && usage
command="${1}"
shift
case "${command}" in
    help)
        if [[ -n ${1} ]]; then
            command_exists ${1} && ( ${1}_usage || invalid_command ${1} )
        else
            usage
        fi;;

    list)
        list_command;;

    start)
        start_command $*;;

    stop)
        stop_command $*;;

    status)
        status_command $*;;

    current)
        print_current;;

    destroy)
        destroy;;

    *) invalid_command "${command}";;
esac

success=true

trap shutdown EXIT

#!/bin/bash
# $2 => script
# $3 => Env
# $4 => Options
cd "$(dirname "$0")"
. mashironrc
systemd-nspawn --volatile --private-network --register=no -qD $VM /usr/container.sh "$1" "$2" ${@:3} 

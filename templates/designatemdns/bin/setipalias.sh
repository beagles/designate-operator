#!/bin/bash
#
# Copyright 2024 Red Hat Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License. You may obtain
# a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
# License for the specific language governing permissions and limitations
# under the License.
set -ex

# expect that the common.sh is in the same dir as the calling script
#
SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
SVC_CFG_MERGED=/var/lib/config-data/merged/designate.conf

IPADDR=$(/usr/local/bin/container-scripts/setipalias.py)
if [ $? -eq 0 ] && [ -n "$IPADDR" ]; then
    crudini --set $SVC_CFG_MERGED 'service:mdns' 'listen' "${IPADDR}:5354"
fi

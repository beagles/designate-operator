#!/bin/bash

#cp -a /var/lib/config-data/named.conf /var/lib/config-data/merged/named.conf
# cp -a /var/lib/config-data/named /var/lib/config-data/merged

set -ex:

# expect that the common.sh is in the same dir as the calling script
SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
. ${SCRIPTPATH}/common.sh --source-only

# Merge all templates from config CM
for dir in /var/lib/config-data/default; do
  merge_config_dir ${dir}
done

mkdir /var/lib/config-data/merged/named
cp -f /var/lib/config-data/default/named/* /var/lib/config-data/merged/named/



#!/bin/bash

REPO=$(cat .git/config | grep url | awk -F 'https://' '{print $2}' \
    | rev | cut -c5- | rev)

if grep "replace ${REPO}" $@ 2>&1 >/dev/null ; then
    echo "ERROR: Don't commit a replacement in go.mod!"
    exit 1
fi

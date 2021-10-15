#!/usr/bin/env bash

IFS='
'
for line in $(grep -R "^type .* interface {$"); do
    file="$(echo "${line}" | cut -d: -f1)"
    interface="$(echo "${line}" | cut -d: -f2 | awk '{print $2}')"
    package="$(basename "$(pwd)")/$(dirname "${file}")"
    destinationDirectory="$(dirname "${file}")/mock"
    destination="${destinationDirectory}/mock_$(basename "${file}")"

    echo "Generating mock for ${interface}"
    
    if [[ ! -d "${destinationDirectory}" ]]; then
        mkdir -p "${destinationDirectory}"
    fi

    mockgen -package=mock -destination="${destination}" "${package}" "${interface}"

done

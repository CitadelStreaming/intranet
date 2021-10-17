#!/usr/bin/env bash

./bin/oci-exec build -f "./Dockerfile" --tag=citadel-intranet:latest --tag=citadel-intranet:$(git describe --tags --dirty)

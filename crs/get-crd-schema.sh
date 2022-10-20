#!/usr/bin/env bash

set -o errexit
set -o pipefail

command -v swagger >/dev/null 2>&1 || { echo >&2 "swagger not installed.  Aborting."; exit 1; }
command -v yq >/dev/null 2>&1 || { echo >&2 "yq not installed.  Aborting."; exit 1; }
command -v jq >/dev/null 2>&1 || { echo >&2 "jq not installed.  Aborting."; exit 1; }

CN_COMMIT=$(go list -m github.com/haproxytech/client-native/v3 | sed 's/^.*-//')

if [ -z "$1" ]; then echo >&2 "No model name supplied.  Aborting."; exit 1; fi
if [ -z "$CN_COMMIT" ]; then echo >&2 "Unable to get git commit for CN module.  Aborting."; exit 1; fi

swagger expand https://raw.githubusercontent.com/haproxytech/client-native/$CN_COMMIT/specification/build/haproxy_spec.yaml |
	yq |
  jq --arg MODEL $1 '.["definitions"][$MODEL]| 
	  walk(
        if type == "object" then with_entries(
          if .key == "x-nullable" then
            if .value == false then
              empty
            else
              .key = "nullable"
            end
          elif (.key | contains("x-")) then
            empty
          else
            .
          end
        ) else . end
    )' |
  yq -y

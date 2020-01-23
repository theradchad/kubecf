#!/usr/bin/env bash

set -o errexit -o nounset

workspace=$(bazel info workspace)

findcmd=(
    find "${workspace}"
    -type f
    '('
        -path "${workspace}/deploy/helm/kubecf/values.*"
        -or
            '(' -not -path "${workspace}/deploy/helm/kubecf/*" ')'
    ')'
    '('
        '(' -name '*.yaml' -not -name '*.tmpl.yaml' ')'
        -or
        '(' -name '*.yml' -not -name '*.tmpl.yaml' ')'
    ')'
)

# shellcheck disable=SC2046
# We want word splitting with find.
bazel run //dev/linters/yamllint -- \
      -d "{extends: relaxed, rules: {line-length: {max: 120}}}" \
      --strict $("${findcmd[@]}")

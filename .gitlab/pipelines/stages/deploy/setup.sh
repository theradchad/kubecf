#!/usr/bin/env bash

set -o errexit -o nounset

# shellcheck disable=SC1090
source "$(bazel info workspace)/.gitlab/pipelines/config/config.sh"

service docker start

wait_for_docker() {
  local timeout="120"
  until docker info 1> /dev/null 2> /dev/null || [[ "$timeout" == "0" ]]; do sleep 1; timeout=$((timeout - 1)); done
  if [[ "${timeout}" == 0 ]]; then return 1; fi
  return 0
}

echo "Waiting for docker to become available..."
wait_for_docker || {
  >&2 echo "Timed out waiting for docker"
  exit 1
}

bazel run //dev/kind:start

# shellcheck disable=SC1090
source "$(bazel info workspace)/.gitlab/pipelines/config/kind_kubeconfig.sh"

# Fix kind's coredns upstream name server to use a reachable name server, given that the one used by
# default is the outer k8s cluster coredns, which is not reachable from the kind network.
bazel run @kubectl//:kubectl -- patch deployment --namespace kube-system coredns --patch '{"spec": {"template": {"spec": {"dnsPolicy": "None", "dnsConfig": {"nameservers": ["8.8.8.8", "8.8.4.4"]}}}}}'

bazel run @kubectl//:kubectl -- create namespace "${KUBECF_NAMESPACE}"

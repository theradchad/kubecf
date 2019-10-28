#!/usr/bin/env bash

set -o errexit -o nounset

KUBECONFIG="$(bazel run @kind//:kind -- get kubeconfig-path --name="scf")"
export KUBECONFIG

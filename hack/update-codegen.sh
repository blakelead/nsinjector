#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

bash vendor/k8s.io/code-generator/generate-groups.sh \
  "deepcopy,client,informer,lister" \
  github.com/blakelead/nsinjector/pkg/generated \
  github.com/blakelead/nsinjector/pkg/apis\
  namespaceresourcesinjector:v1alpha1 \
  --go-header-file hack/boilerplate.go.txt
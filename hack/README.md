# Code generator

This folder contains files for generating custom resources dependencies.
This file is just a reminder of what I've done to make this work.

## Before generation

Change `k8s.io/client-go v11.0.0+incompatible` to `k8s.io/client-go v0.18.4` or higher.

Create the following files:

- `pkg/apis/namespaceresourcesinjector/v1alpha1/doc.go`
- `pkg/apis/namespaceresourcesinjector/v1alpha1/register.go`
- `pkg/apis/namespaceresourcesinjector/v1alpha1/types.go`

Add `// +genclient:nonNamespaced` comment after `// +genclient` to `pkg/apis/namespaceresourcesinjector/v1alpha1/types.go` in order to tell the generator that the resource is cluster scoped.

## Generation

Go mod vendor then execute generator. It will generate:

- `pkg/apis/namespaceresourcesinjector/v1alpha1/zz_generated.deepcopy.go`
- `pkg/generated/`

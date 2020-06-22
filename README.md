# nsinjector

**nsinjector** is a Kubernetes controller that automatically deploys resources into a namespace when it is created.

## How to use it

You can find a chart with an example `values.yaml` file in `deploy/helm`.
With helm3, you can deploy it with:

```bash
helm install nsinjector-controller ./deploy/helm
```

Alternatively, you can manually deploy manifest stored in `deploy/k8s`:

```bash
# Deploy CRD first. If your cluster is >= v1.16, you can use namespaceresourcesinjector-crd-1.16.yaml instead
kubectl apply -f deploy/k8s/namespaceresourcesinjector-crd.yaml
# Deploy the controller
kubectl apply -f deploy/k8s/nsinjector-controller.yaml
# Then deploy an injector custom resource
# This is the file that you'll want to customize to your needs
kubectl apply -f deploy/k8s/namespaceresourcesinjector.yaml
```

## More details

The controller watches:

- `namespaceresourcesinjector` creation, update or deletion
- `namespace` creation or deletion

When a custom resource `namespaceresourcesinjector` is created or updated, the controller will fetch all namespaces matching `.spec.namespaces` of the `namespaceresourcesinjector` and deploy (or update) `spec.resources` resources in those namespaces (or cluster-wide for cluster-scoped resources).

When a custom resource `namespaceresourcesinjector` is deleted, associated resources will be deleted from matching namespaces.

When a `namespace` is created, it will be checked against existing `namespaceresourcesinjector` and resources will be deployed in it if there's a match.

`namespaceresourcesinjector` keep a reference of matched namespaces in `status.injectedNamespaces`. When a corresponding namespace is deleted, it will be removed from the custom resource status.

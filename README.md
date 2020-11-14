# nsinjector

**nsinjector** is a Kubernetes controller that automatically deploys resources into a namespace when it is created.

Here's a blog post I wrote about it: [Deploy Kubernetes resources automatically with nsinjector
](https://blog.blakelead.com/posts/2020/11/12/nsinjector/).

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

## Example

When a namespace starting with `dev-` is created, the following resource will automatically inject a role and rolebinding in it:

```yaml
kind: NamespaceResourcesInjector
apiVersion: blakelead.com/v1alpha1
metadata:
  name: nri-test
spec:
  namespaces:
  - dev-.*
  resources:
  - |
    apiVersion: rbac.authorization.k8s.io/v1
    kind: Role
    metadata:
      name: dev-role
    rules:
      - apiGroups: [""]
        resources: ["pods","pods/portforward", "services", "deployments", "ingresses"]
        verbs: ["list", "get"]
      - apiGroups: [""]
        resources: ["pods/portforward"]
        verbs: ["create"]
      - apiGroups: [""]
        resources: ["namespaces"]
        verbs: ["list", "get"]
  - |
    apiVersion: rbac.authorization.k8s.io/v1
    kind: RoleBinding
    metadata:
      name: dev-rolebinding
    subjects:
    - kind: User
      name: dev
    roleRef:
      kind: Role
      name: dev-role
      apiGroup: rbac.authorization.k8s.io
```

- `namespaces`:  a list of namespace names or regex
- `resources`: a list of any Kubernetes resources

## Contributing

Although this project is currently used in production, it is relatively young and has not been extensively tested. Suggestions and contributions are therefore very welcome.

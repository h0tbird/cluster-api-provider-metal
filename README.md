# cluster-api-provider-metal
Bare metal CAPI provider

## Development

Start kind:
```
curl -s https://raw.githubusercontent.com/tilt-dev/kind-local/master/kind-with-registry.sh > ~/bin/kind-local
chmod +x ~/bin/kind-local
kind-local
```

Start tilt:
```
cd ${CAPI}
tilt up
```

## Scaffolding
Initial project scaffolding using *[kubebuilder](https://github.com/kubernetes-sigs/kubebuilder)* `2.3.1`:
```
go mod init github.com/h0tbird/cluster-api-provider-metal
kubebuilder init --domain cluster.x-k8s.io --license apache2 --owner "Open Source Community"
kubebuilder create api --group infrastructure --version v1alpha3 --kind BareMetalCluster
kubebuilder create api --group infrastructure --version v1alpha3 --kind BareMetalMachine
```

Patch the *Makefile* and the API types:
```
sed -i 's#controller:latest#localhost:5000/cluster-api-capm-controller:latest#' Makefile
sed -i 's/controller-gen@v0.2.5/controller-gen@v0.4.0/g' Makefile
sed -i '0,/+kubebuilder:object:root=true/s##&\
// +kubebuilder:subresource:status\
// +kubebuilder:storageversion#' \
api/v1alpha3/*_types.go
```

Generate the CRDs with *[controller-gen](https://github.com/kubernetes-sigs/controller-tools)*:
```
rm -f $(which controller-gen) && make manifests
```

Deploy the controller and the CRDs:
```
make deploy
kubectl apply -f config/samples
```

Setup CAPI tilt:
```
cat > tilt-settings.json << EOF
{
  "default_registry": "localhost:5000",
  "provider_repos": ["../cluster-api-provider-metal"],
  "enable_providers": ["metal", "kubeadm-bootstrap", "kubeadm-control-plane"]
}
EOF
```

## Use `kind` and `clusterctl` to deploy CAPD clusters

Allow CAPD to use the local docker:
```
KIND_EXPERIMENTAL_DOCKER_NETWORK=bridge kind create cluster --config <(
echo "kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraMounts:
    - hostPath: /var/run/docker.sock
      containerPath: /var/run/docker.sock") && \
while ! k wait --for=condition=Ready pod --all -A; do echo "Wait again"; done
```

Deploy the management cluster:
```
clusterctl init -v 1 --infrastructure docker && \
while ! k wait --for=condition=Ready pod --all -A; do echo "Wait again"; done
```

Deploy a workload `CLUSTER=foo`:
```
clusterctl config cluster ${CLUSTER} --flavor development \
--kubernetes-version v1.19.1 \
--control-plane-machine-count=1 \
--worker-machine-count=2 | \
k apply -f -
```

Get the kubeconfig:
```
clusterctl get kubeconfig ${CLUSTER} | sed -e "
  s/server:.*/server: https:\/\/$(docker port ${CLUSTER}-lb 6443/tcp | sed "s/0.0.0.0/127.0.0.1/")/g;
  s/certificate-authority-data:.*/insecure-skip-tls-verify: true/g;
" > /tmp/${CLUSTER}.kubeconfig
```

Upsert the kubeconfig:
```
KUBECONFIG=/tmp/${CLUSTER}.kubeconfig:${KUBECONFIG} \
k config view --flatten | sponge ${KUBECONFIG}
```

Deploy a CNI solution:
```
k apply -f https://docs.projectcalico.org/v3.15/manifests/calico.yaml && \
while ! k wait --for=condition=Ready pod --all -A; do echo "Wait again"; done
```

Verify:
```
k config use-context kind-kind
k get cluster,kubeadmcontrolplane -A
```

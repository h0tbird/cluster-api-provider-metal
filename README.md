# cluster-api-provider-metal
Bare metal CAPI provider

## Development
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

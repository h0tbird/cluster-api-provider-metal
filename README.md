# cluster-api-provider-metal
Bare metal CAPI provider

## Development
Initial project scaffolding with *kubebuilder* `2.3.1`:
```
go mod init github.com/h0tbird/cluster-api-provider-metal
kubebuilder init --domain cluster.x-k8s.io --license apache2 --owner "Open Source Community"
kubebuilder create api --group infrastructure --version v1alpha3 --kind BareMetalCluster
kubebuilder create api --group infrastructure --version v1alpha3 --kind BareMetalMachine
```

Patch the *Makefile* and the API types:
```
gsed -i 's#controller:latest#registry:5000/cluster-api-capm-controller:latest#' Makefile
gsed -i 's/controller-gen@v0.2.5/controller-gen@v0.2.9/g' Makefile
gsed -i 's/trivialVersions=true/crdVersions=v1/g' Makefile
gsed -i '0,/+kubebuilder:object:root=true/s##&\
// +kubebuilder:subresource:status\
// +kubebuilder:storageversion#' \
api/v1alpha3/*_types.go
```

Generate the CRDs with *controller-gen* `0.2.8`:
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

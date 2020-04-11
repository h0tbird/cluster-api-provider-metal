# cluster-api-provider-metal
Bare metal CAPI provider

## Development
Initial project scaffolding with kubebuilder `2.3.1`:
```
go mod init github.com/h0tbird/cluster-api-provider-metal
kubebuilder init --domain cluster.x-k8s.io --license apache2 --owner "Marc Villacorta"
kubebuilder create api --group infrastructure --version v1alpha3 --kind BareMetalCluster
kubebuilder create api --group infrastructure --version v1alpha3 --kind BareMetalMachine
```

Patch and generate the CRDs with controller-gen `0.2.8`:
```
gsed '0,/+kubebuilder:object:root=true/s##&\
// +kubebuilder:subresource:status\
// +kubebuilder:storageversion\#' api/v1alpha3/*_types.go
gsed -i 's/controller-gen@v0.2.5/controller-gen@v0.2.8/g' Makefile
gsed -i 's/trivialVersions=true/crdVersions=v1/g' Makefile
rm -f $(which controller-gen) && make manifests
```

Install the CRDs into the cluster:
```
make install
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

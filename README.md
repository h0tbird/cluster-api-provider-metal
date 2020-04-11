# cluster-api-provider-metal
Bare metal CAPI provider

## Development
Initial project scaffolding.
```
go mod init github.com/h0tbird/cluster-api-provider-metal
kubebuilder init --domain cluster.x-k8s.io --license apache2 --owner "Marc Villacorta"
kubebuilder create api --group infrastructure --version v1alpha1 --kind BareMetalCluster
kubebuilder create api --group infrastructure --version v1alpha1 --kind BareMetalMachine
```

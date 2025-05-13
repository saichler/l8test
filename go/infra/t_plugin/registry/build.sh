go build -buildmode=plugin -o Tests-0-registry.so TestRegistryPlugin.go
cp Tests-0-registry.so Tests-1-registry.so
cp Tests-0-registry.so Tests-2-registry.so

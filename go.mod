module operator/crda-operator

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/google/martian v2.1.0+incompatible
	github.com/jdob/visitors-operator v0.0.0-20191024200828-5b18c79fe98b
	github.com/operator-framework/operator-sdk v0.17.0
	github.com/spf13/pflag v1.0.5
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.5.2
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	k8s.io/client-go => k8s.io/client-go v0.17.4 // Required by prometheus-operator
)

# CRDA Operator

Operator for Code Ready Dependency Analytics Plateform Deployment

[![Go Report Card](https://goreportcard.com/badge/github.com/deepak1725/crda-operator)](https://goreportcard.com/report/github.com/deepak1725/crda-operator) 

This Operator will deploy necessary Services in [CRDA Plateform](https://github.com/fabric8-analytics). 

## Local Installation/ Testing:

* Have [Minikube Installed](https://kubernetes.io/docs/tasks/tools/install-minikube/) .
* Create Role `kubectl apply -f deploy/role.yaml`
* Create Service Account `kubectl apply -f deploy/service_account.yaml`
* Create Role Binding `kubectl apply -f deply/service_binding.yaml`
* Deploy CRD: `kubectl apply -f deploy/crds/openshift.com_codereadyanalytics_crd.yaml`
* Deploy CR: `kubectl apply -f deploy/crds/openshift.com_v1alpha1_codereadyanalytics_cr.yaml`

Now run Operator locally:
`operator-sdk run up --local`

This will deploy all the custom resources defined in Specified Kind.


## Production Installation via [OLM](https://sdk.operatorframework.io/docs/olm-integration/user-guide/) :

* Install OLM on Server `operator-sdk olm install`
* Install Operator `operator-sdk run packagemanifests --operator-version 0.1.0 --olm`
* Create Role `kubectl apply -f deploy/role.yaml`
* Create Service Account `kubectl apply -f deploy/service_account.yaml`
* Create Role Binding `kubectl apply -f deply/service_binding.yaml`  
* Deploy CR: `kubectl apply -f deploy/crds/openshift.com_v1alpha1_codereadyanalytics_cr.yaml`


Please note all resources are deployed on `default` namespace. You can overide this by `-n $mycustomnamespace` 



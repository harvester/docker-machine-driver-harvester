module github.com/harvester/docker-machine-driver-harvester

go 1.25.0

replace (
	github.com/docker/docker => github.com/moby/moby v1.4.2-0.20170731201646-1009e6a40b29 // oras dep requires a replace is set
	github.com/golang/glog => github.com/golang/glog v1.2.4 // CVE-2024-45339
	github.com/harvester/harvester => github.com/harvester/harvester v0.0.0-20250831134456-11f4f0204b16
	github.com/longhorn/longhorn-manager => github.com/longhorn/longhorn-manager v0.0.0-20240731060708-b10fa4a47fff
	github.com/openshift/api => github.com/openshift/api v0.0.0-20191219222812-2987a591a72c
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20200521150516-05eb9880269c
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring => github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.68.0
	github.com/rancher/machine => github.com/rancher/machine v0.15.0-rancher119
	github.com/rancher/rancher/pkg/apis => github.com/rancher/rancher/pkg/apis v0.0.0-20240919204204-3da2ae0cabd1
	github.com/rancher/rancher/pkg/client => github.com/rancher/rancher/pkg/client v0.0.0-20240919204204-3da2ae0cabd1
	golang.org/x/net => golang.org/x/net v0.38.0 // CVE-2025-22870, CVE-2025-22872
	golang.org/x/oauth2 => golang.org/x/oauth2 v0.27.0 // CVE-2025-22868

	k8s.io/api => k8s.io/api v0.31.13
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.31.13
	k8s.io/apimachinery => k8s.io/apimachinery v0.31.13
	k8s.io/apiserver => k8s.io/apiserver v0.31.13
	k8s.io/client-go => k8s.io/client-go v0.31.13
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.31.13
	k8s.io/component-base => k8s.io/component-base v0.31.13
	k8s.io/cri-api => k8s.io/cri-api v0.31.13
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.31.13
	k8s.io/dynamic-resource-allocation => k8s.io/dynamic-resource-allocation v0.31.13
	k8s.io/endpointslice => k8s.io/endpointslice v0.31.13
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.31.13
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.31.13
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20240228011516-70dd3763d340
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.31.13
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.31.13
	k8s.io/kubernetes => k8s.io/kubernetes v1.31.13
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.31.13
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.31.13

	kubevirt.io/api => github.com/kubevirt/api v1.5.0
	kubevirt.io/client-go => github.com/kubevirt/client-go v1.5.0
	kubevirt.io/kubevirt => kubevirt.io/kubevirt v1.5.0
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.19.7
	sigs.k8s.io/structured-merge-diff => sigs.k8s.io/structured-merge-diff v0.0.0-20190302045857-e85c7b244fd2
)

require (
	dario.cat/mergo v1.0.2
	github.com/ghodss/yaml v1.0.0
	github.com/harvester/harvester v1.7.0-dev-20250831
	github.com/k8snetworkplumbingwg/network-attachment-definition-client v1.7.5
	github.com/rancher/machine v0.15.0-rancher99
	github.com/rancher/wrangler v1.1.2
	github.com/stretchr/testify v1.10.0
	k8s.io/api v0.33.1
	k8s.io/apimachinery v0.33.1
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/utils v0.0.0-20250502105355-0f33e8f1c979
	kubevirt.io/api v1.5.0
)

require (
	emperror.dev/errors v0.8.1 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/c9s/goprocinfo v0.0.0-20210130143923-c95fcf8c64a8 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cisco-open/operator-tools v0.37.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/docker/docker v27.1.1+incompatible // indirect
	github.com/emicklei/go-restful/v3 v3.12.2 // indirect
	github.com/evanphx/json-patch v5.9.11+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.9.11 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/fxamacker/cbor/v2 v2.8.0 // indirect
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-openapi/jsonpointer v0.21.1 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.2.5 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gnostic-models v0.6.9 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/handlers v1.5.2 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/harvester/go-common v0.0.0-20250109132713-e748ce72a7ba // indirect
	github.com/harvester/harvester-network-controller v1.6.0-rc3 // indirect
	github.com/iancoleman/orderedmap v0.3.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/k8snetworkplumbingwg/whereabouts v0.8.0 // indirect
	github.com/kube-logging/logging-operator v0.0.0-20250424202944-7e1f9aad6e21 // indirect
	github.com/kube-logging/logging-operator/pkg/sdk v0.12.0 // indirect
	github.com/kubernetes-csi/external-snapshotter/client/v4 v4.2.0 // indirect
	github.com/longhorn/go-common-libs v0.0.0-20250215052214-151615b29f8e // indirect
	github.com/longhorn/longhorn-manager v1.8.1 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/onsi/gomega v1.37.0 // indirect
	github.com/openshift/api v0.0.0 // indirect
	github.com/openshift/client-go v3.9.0+incompatible // indirect
	github.com/openshift/custom-resource-status v1.1.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.82.0 // indirect
	github.com/prometheus/client_golang v1.22.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rancher/aks-operator v1.9.2 // indirect
	github.com/rancher/eks-operator v1.9.2 // indirect
	github.com/rancher/fleet/pkg/apis v0.10.0 // indirect
	github.com/rancher/gke-operator v1.9.2 // indirect
	github.com/rancher/lasso v0.2.3 // indirect
	github.com/rancher/norman v0.0.0-20241001183610-78a520c160ab // indirect
	github.com/rancher/rancher/pkg/apis v0.0.0 // indirect
	github.com/rancher/rke v1.6.2 // indirect
	github.com/rancher/system-upgrade-controller/pkg/apis v0.0.0-20240301001845-4eacc2dabbde // indirect
	github.com/rancher/wrangler/v3 v3.2.2 // indirect
	github.com/robfig/cron v1.2.0 // indirect
	github.com/shirou/gopsutil/v3 v3.24.5 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/cast v1.8.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/exp v0.0.0-20250506013437-ce4c2cf36ca6 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/oauth2 v0.30.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/term v0.32.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	gomodules.xyz/jsonpatch/v2 v2.5.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/evanphx/json-patch.v4 v4.12.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiextensions-apiserver v0.33.1 // indirect
	k8s.io/apiserver v0.33.1 // indirect
	k8s.io/component-base v0.33.1 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-aggregator v0.33.1 // indirect
	k8s.io/kube-openapi v0.31.9 // indirect
	k8s.io/kubernetes v1.32.2 // indirect
	kubevirt.io/client-go v1.5.0 // indirect
	kubevirt.io/containerized-data-importer-api v1.61.0 // indirect
	kubevirt.io/controller-lifecycle-operator-sdk/api v0.0.0-20220329064328-f3cc58c6ed90 // indirect
	kubevirt.io/kubevirt v1.5.0 // indirect
	sigs.k8s.io/cli-utils v0.37.2 // indirect
	sigs.k8s.io/cluster-api v1.7.3 // indirect
	sigs.k8s.io/controller-runtime v0.20.4 // indirect
	sigs.k8s.io/json v0.0.0-20241014173422-cfa47c3a1cc8 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.7.0 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

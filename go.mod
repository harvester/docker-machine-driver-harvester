module github.com/harvester/docker-machine-driver-harvester

go 1.18

replace (
	github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go v3.2.1-0.20200107013213-dc14462fd587+incompatible
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/docker/docker v1.4.2-0.20200203170920-46ec8731fbce
	github.com/go-kit/kit => github.com/go-kit/kit v0.3.0
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.1

	github.com/harvester/harvester => github.com/harvester/harvester v0.0.2-0.20220712022404-7158c858201a

	github.com/knative/pkg => github.com/rancher/pkg v0.0.0-20190514055449-b30ab9de040e
	github.com/openshift/api => github.com/openshift/api v0.0.0-20191219222812-2987a591a72c
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20200521150516-05eb9880269c
	github.com/operator-framework/operator-lifecycle-manager => github.com/operator-framework/operator-lifecycle-manager v0.0.0-20190128024246-5eb7ae5bdb7a
	github.com/rancher/rancher/pkg/apis => github.com/rancher/rancher/pkg/apis v0.0.0-20211208233239-77392a65423d
	github.com/rancher/rancher/pkg/client => github.com/rancher/rancher/pkg/client v0.0.0-20211208233239-77392a65423d

	helm.sh/helm/v3 => github.com/rancher/helm/v3 v3.7.1-rancher.1
	k8s.io/api => k8s.io/api v0.23.7
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.23.7
	k8s.io/apimachinery => k8s.io/apimachinery v0.23.7
	k8s.io/apiserver => k8s.io/apiserver v0.23.7
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.23.7
	k8s.io/client-go => k8s.io/client-go v0.23.7
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.23.7
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.23.7
	k8s.io/code-generator => k8s.io/code-generator v0.23.7
	k8s.io/component-base => k8s.io/component-base v0.23.7
	k8s.io/component-helpers => k8s.io/component-helpers v0.23.7
	k8s.io/controller-manager => k8s.io/controller-manager v0.23.7
	k8s.io/cri-api => k8s.io/cri-api v0.23.7
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.23.7
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.23.7
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.23.7
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.23.7
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.23.7
	k8s.io/kubectl => k8s.io/kubectl v0.23.7
	k8s.io/kubelet => k8s.io/kubelet v0.23.7
	k8s.io/kubernetes => k8s.io/kubernetes v1.23.7
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.23.7
	k8s.io/metrics => k8s.io/metrics v0.23.7
	k8s.io/mount-utils => k8s.io/mount-utils v0.23.7
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.23.7
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.23.7

	kubevirt.io/api => github.com/kubevirt/api v0.53.1
	kubevirt.io/client-go => github.com/kubevirt/client-go v0.53.1
	kubevirt.io/containerized-data-importer => github.com/kubevirt/containerized-data-importer v1.41.0
	kubevirt.io/containerized-data-importer-api => github.com/kubevirt/containerized-data-importer-api v1.47.0
	sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v0.4.4
	sigs.k8s.io/structured-merge-diff => sigs.k8s.io/structured-merge-diff v0.0.0-20190302045857-e85c7b244fd2
)

// Bump golang.org/x/crypto to fix https://github.com/advisories/GHSA-8c26-wmh5-6g9v
replace golang.org/x/crypto => golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4

require (
	github.com/ghodss/yaml v1.0.0
	github.com/harvester/harvester v1.0.2-rc2
	github.com/imdario/mergo v0.3.12
	github.com/k8snetworkplumbingwg/network-attachment-definition-client v0.0.0-20200331171230-d50e42f2b669
	github.com/rancher/machine v0.15.0-rancher85
	github.com/rancher/wrangler v1.0.1-0.20220520195731-8eeded9bae2a
	k8s.io/api v0.24.0
	k8s.io/apimachinery v0.24.0
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/utils v0.0.0-20211116205334-6203023598ed
	kubevirt.io/api v0.0.0-20220430221853-33880526e414
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/docker v20.10.12+incompatible // indirect
	github.com/emicklei/go-restful v2.15.0+incompatible // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/jinzhu/copier v0.3.5 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kubernetes-csi/external-snapshotter/v2 v2.1.1 // indirect
	github.com/longhorn/longhorn-manager v1.3.0-rc2 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/openshift/custom-resource-status v1.1.2 // indirect
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.12.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/rancher/lasso v0.0.0-20220519004610-700f167d8324 // indirect
	github.com/rancher/system-upgrade-controller/pkg/apis v0.0.0-20210727200656-10b094e30007 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.0.0 // indirect
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220319134239-a9b59b0215f8 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gotest.tools v2.2.0+incompatible // indirect
	k8s.io/apiextensions-apiserver v0.24.0 // indirect
	k8s.io/component-base v0.23.7 // indirect
	k8s.io/klog/v2 v2.40.1 // indirect
	k8s.io/kube-openapi v0.0.0-20220124234850-424119656bbf // indirect
	kubevirt.io/containerized-data-importer-api v1.47.0 // indirect
	kubevirt.io/controller-lifecycle-operator-sdk/api v0.0.0-20220329064328-f3cc58c6ed90 // indirect
	kubevirt.io/kubevirt v0.53.1 // indirect
	sigs.k8s.io/cluster-api v0.4.4 // indirect
	sigs.k8s.io/controller-runtime v0.11.2 // indirect
	sigs.k8s.io/json v0.0.0-20211020170558-c049b76a60c6 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

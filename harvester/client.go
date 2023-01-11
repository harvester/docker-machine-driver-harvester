package harvester

import (
	"fmt"
	"os"
	"strings"

	harvsterv1 "github.com/harvester/harvester/pkg/apis/harvesterhci.io/v1beta1"
	harvclient "github.com/harvester/harvester/pkg/generated/clientset/versioned"
	"github.com/harvester/harvester/pkg/generated/clientset/versioned/scheme"
	cniv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	"github.com/rancher/wrangler/pkg/kubeconfig"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

type Client struct {
	RestConfig                *rest.Config
	KubeVirtSubresourceClient *rest.RESTClient
	HarvesterClient           *harvclient.Clientset
	KubeClient                *kubernetes.Clientset
}

func NewClientFromRestConfig(restConfig *rest.Config) (*Client, error) {
	subresourceConfig := rest.CopyConfig(restConfig)
	subresourceConfig.GroupVersion = &schema.GroupVersion{Group: kubevirtv1.SubresourceGroupName, Version: kubevirtv1.ApiLatestVersion}
	subresourceConfig.APIPath = "/apis"
	subresourceConfig.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	kubeVirtSubresourceClient, err := rest.RESTClientFor(subresourceConfig)
	if err != nil {
		return nil, err
	}
	harvClient, err := harvclient.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	kubeClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return &Client{
		RestConfig:                restConfig,
		KubeVirtSubresourceClient: kubeVirtSubresourceClient,
		HarvesterClient:           harvClient,
		KubeClient:                kubeClient,
	}, nil
}

func NamespacedNameParts(namespacedName string) (string, string, error) {
	parts := strings.Split(namespacedName, "/")
	switch len(parts) {
	case 1:
		return "", parts[0], nil
	case 2:
		return parts[0], parts[1], nil
	default:
		err := fmt.Errorf("unexpected namespacedName format (%q), expected %q or %q. ", namespacedName, "namespace/name", "name")
		return "", "", err
	}
}

func NamespacedNamePartsByDefault(namespacedName string, defaultNamespace string) (string, string, error) {
	namespace, name, err := NamespacedNameParts(namespacedName)
	if err != nil {
		return "", "", err
	}
	if namespace == "" {
		namespace = defaultNamespace
	}
	return namespace, name, nil
}

func (d *Driver) getRestConfig() (*rest.Config, error) {
	if d.KubeConfigContent == "" {
		return kubeconfig.GetNonInteractiveClientConfig(os.Getenv("KUBECONFIG")).ClientConfig()
	}
	return clientcmd.RESTConfigFromKubeConfig([]byte(d.KubeConfigContent))
}

func (d *Driver) getClient() (*Client, error) {
	if d.client != nil {
		return d.client, nil
	}
	restConfig, err := d.getRestConfig()
	if err != nil {
		return nil, err
	}
	c, err := NewClientFromRestConfig(restConfig)
	if err != nil {
		return nil, err
	}
	d.client = c
	return d.client, nil
}

func (d *Driver) getSetting(name string) (*harvsterv1.Setting, error) {
	c, err := d.getClient()
	if err != nil {
		return nil, err
	}
	return c.HarvesterClient.HarvesterhciV1beta1().Settings().Get(d.ctx, name, metav1.GetOptions{})
}

func (d *Driver) getImage(imageName string) (*harvsterv1.VirtualMachineImage, error) {
	c, err := d.getClient()
	if err != nil {
		return nil, err
	}
	namespace, name, err := NamespacedNamePartsByDefault(imageName, d.VMNamespace)
	if err != nil {
		return nil, err
	}
	return c.HarvesterClient.HarvesterhciV1beta1().VirtualMachineImages(namespace).Get(d.ctx, name, metav1.GetOptions{})
}

func (d *Driver) getStorageClass(storageClassName string) (*storagev1.StorageClass, error) {
	c, err := d.getClient()
	if err != nil {
		return nil, err
	}
	return c.KubeClient.StorageV1().StorageClasses().Get(d.ctx, storageClassName, metav1.GetOptions{})
}

func (d *Driver) getKeyPair(keyPairName string) (*harvsterv1.KeyPair, error) {
	c, err := d.getClient()
	if err != nil {
		return nil, err
	}
	namespace, name, err := NamespacedNamePartsByDefault(keyPairName, d.VMNamespace)
	if err != nil {
		return nil, err
	}
	return c.HarvesterClient.HarvesterhciV1beta1().KeyPairs(namespace).Get(d.ctx, name, metav1.GetOptions{})
}

func (d *Driver) getNetwork(networkName string) (*cniv1.NetworkAttachmentDefinition, error) {
	c, err := d.getClient()
	if err != nil {
		return nil, err
	}
	namespace, name, err := NamespacedNamePartsByDefault(networkName, d.VMNamespace)
	if err != nil {
		return nil, err
	}
	return c.HarvesterClient.K8sCniCncfIoV1().NetworkAttachmentDefinitions(namespace).Get(d.ctx, name, metav1.GetOptions{})
}

func (d *Driver) getVMI() (*kubevirtv1.VirtualMachineInstance, error) {
	c, err := d.getClient()
	if err != nil {
		return nil, err
	}
	return c.HarvesterClient.KubevirtV1().VirtualMachineInstances(d.VMNamespace).Get(d.ctx, d.MachineName, metav1.GetOptions{})
}

func (d *Driver) getVM() (*kubevirtv1.VirtualMachine, error) {
	c, err := d.getClient()
	if err != nil {
		return nil, err
	}
	return c.HarvesterClient.KubevirtV1().VirtualMachines(d.VMNamespace).Get(d.ctx, d.MachineName, metav1.GetOptions{})
}

func (d *Driver) updateVM(newVM *kubevirtv1.VirtualMachine) (*kubevirtv1.VirtualMachine, error) {
	c, err := d.getClient()
	if err != nil {
		return nil, err
	}
	return c.HarvesterClient.KubevirtV1().VirtualMachines(d.VMNamespace).Update(d.ctx, newVM, metav1.UpdateOptions{})
}

func (d *Driver) deleteVM() error {
	c, err := d.getClient()
	if err != nil {
		return err
	}
	propagationPolicy := metav1.DeletePropagationForeground
	return c.HarvesterClient.KubevirtV1().VirtualMachines(d.VMNamespace).Delete(d.ctx, d.MachineName, metav1.DeleteOptions{
		PropagationPolicy: &propagationPolicy,
	})
}

func (d *Driver) putVMSubResource(subResource string) error {
	c, err := d.getClient()
	if err != nil {
		return err
	}
	return c.KubeVirtSubresourceClient.Put().Namespace(d.VMNamespace).Resource(vmResource).SubResource(subResource).Name(d.MachineName).Do(d.ctx).Error()
}

func (d *Driver) deleteVolume(name string) error {
	c, err := d.getClient()
	if err != nil {
		return err
	}
	return c.KubeClient.CoreV1().PersistentVolumeClaims(d.VMNamespace).Delete(d.ctx, name, metav1.DeleteOptions{})
}

func (d *Driver) createVM(vm *kubevirtv1.VirtualMachine) (*kubevirtv1.VirtualMachine, error) {
	c, err := d.getClient()
	if err != nil {
		return nil, err
	}
	return c.HarvesterClient.KubevirtV1().VirtualMachines(d.VMNamespace).Create(d.ctx, vm, metav1.CreateOptions{})
}

func (d *Driver) createSecret(secret *corev1.Secret) (*corev1.Secret, error) {
	c, err := d.getClient()
	if err != nil {
		return nil, err
	}
	return c.KubeClient.CoreV1().Secrets(d.VMNamespace).Create(d.ctx, secret, metav1.CreateOptions{})
}

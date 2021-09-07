package ks_devops_e2e

import (
	"context"
	"github.com/kubesphere-sigs/ks-devops-e2e/pkg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"testing"
	"time"
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{})
}

var _ = BeforeSuite(func() {
	client, _, err := GetClient()
	Expect(err).To(BeNil())

	Eventually(func() bool {
		if list, err := client.Resource(pkg.GetDeploySchema()).Namespace("kubesphere-devops-system").List(context.TODO(), metav1.ListOptions{}); err == nil {
			for i, _ := range list.Items {
				item := list.Items[i]
				ready := getNestedInt(item.Object, "status", "readyReplicas")
				if ready != 1 {
					return false
				}
			}
		}
		return true
	}, time.Second * 120, time.Second).Should(BeTrue())
})

func getNestedString(obj map[string]interface{}, fields ...string) string {
	val, found, err := unstructured.NestedString(obj, fields...)
	if !found || err != nil {
		return ""
	}
	return val
}

func getNestedInt(obj map[string]interface{}, fields ...string) int64 {
	val, found, err := unstructured.NestedInt64(obj, fields...)
	if !found || err != nil {
		return 0
	}
	return val
}

func GetClient() (client dynamic.Interface, clientSet *kubernetes.Clientset, err error) {
	KubernetesConfigFlags := genericclioptions.NewConfigFlags(false)
	if config, err := KubernetesConfigFlags.ToRESTConfig(); err == nil {
		client, err = dynamic.NewForConfig(config)
		clientSet, err = kubernetes.NewForConfig(config)
	}

	return
}

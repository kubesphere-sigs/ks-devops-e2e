package ks_devops_e2e

import (
	"bytes"
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/kubesphere-sigs/ks-devops-e2e/pkg"
	"github.com/linuxsuren/http-downloader/pkg/exec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"strings"
	"time"
)

var _ = Describe("Pipeline Test", func() {
	stdout := &bytes.Buffer{}
	errout := &bytes.Buffer{}

	var (
		ws, project, name, template string
		namespace string
	)

	ws = strings.ToLower(randomdata.SillyName())
	project = strings.ToLower(randomdata.SillyName())
	name = strings.ToLower(randomdata.SillyName())
	template = "simple"

	Context("Create and run pipeline", func() {
		It("Create Pipeline", func() {
			err := exec.RunCommandWithIO("ks", "", stdout, errout, "pip", "create", "--ws", ws,
				"--template", template, "--name", name, "--project", project)
			Expect(err).To(BeNil())
			result := errout.String()
			Expect(result).To(BeEmpty())

			client, _, err := GetClient()
			Expect(err).To(BeNil())
			Eventually(func() bool {
				var pipelineList *unstructured.UnstructuredList
				if pipelineList, err = client.Resource(pkg.GetPipelineSchema()).List(context.TODO(), metav1.ListOptions{}); err != nil {
					return false
				}

				for i := range pipelineList.Items {
					item := pipelineList.Items[i]
					if item.GetName() == name {
						namespace = item.GetNamespace()
						return true
					}
				}
				return false
			}, time.Minute * 2, time.Second * 3).Should(BeTrue())
		})

		It("Run a Pipeline", func() {
			err := exec.RunCommandWithIO("ks", "", stdout, errout, "pip", "run", "--project", project,
				"-p", name, "-b")
			Expect(err).To(BeNil())
			result := errout.String()
			Expect(result).To(BeEmpty())

			client, _, err := GetClient()
			Expect(err).To(BeNil())
			Eventually(func() bool {
				var pipelineRuns *unstructured.UnstructuredList
				if pipelineRuns, err = client.Resource(pkg.GetPipelineRunSchema()).Namespace(namespace).List(context.TODO(), metav1.ListOptions{}); err != nil {
					return false
				}

				for i := range pipelineRuns.Items {
					item := pipelineRuns.Items[i]
					if item.GetGenerateName() == name {
						phase := getNestedString(item.Object, "status", "phase")
						if phase == "Succeeded" {
							return true
						}
						break
					}
				}
				return false
			}, time.Minute * 10, time.Second * 3).Should(BeTrue())
		})
	})
})

package ks_devops_e2e

import (
	"bytes"
	"github.com/linuxsuren/http-downloader/pkg/exec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pipeline Test", func() {
	stdout := &bytes.Buffer{}
	errout := &bytes.Buffer{}

	Context("", func() {
		It("", func() {

			err := exec.RunCommandWithIO("ks", "", stdout, errout, "pip", "create", "--ws", "simple",
				"--template", "java", "--name", "java", "--project", "test")
			Expect(err).To(BeNil())
			result := errout.String()
			Expect(result).To(BeEmpty())
		})
	})
})

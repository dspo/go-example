package echo

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	sdkapp0 "gitee.com/huajinet/go-example/sdk/app0"
)

var _ = ginkgo.Describe("Echo", func() {
	ginkgo.Context("Test Echo", func() {
		ginkgo.It("normal case", func() {
			echo, err := sdkapp0.NewEcho().Echo()
			gomega.Ω(err).ShouldNot(gomega.HaveOccurred())
			gomega.Ω(echo).Should(gomega.HaveKey("headers"))
		})
	})
})

package e2e

import (
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	_ "gitee.com/huajinet/go-example/test/e2e/echo"
	"gitee.com/huajinet/go-example/test/framework"
)

func TestE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	f := framework.GetFramework()
	f.DeployComponents(t)
	ginkgo.RunSpecs(t, "E2E")
}

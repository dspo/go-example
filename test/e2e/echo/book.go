package echo

import (
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"gitee.com/huajinet/go-example/internal/model"
	sdkapp0 "gitee.com/huajinet/go-example/sdk/app0"
	"gitee.com/huajinet/go-example/test/framework"
)

var _ = ginkgo.Describe("Book Endpoints", func() {
	ginkgo.BeforeEach(func() {
		err := framework.GetFramework().DB().Exec("DELETE FROM `book`").Error
		gomega.Ω(err).ShouldNot(gomega.HaveOccurred())
	})

	ginkgo.Context("Test Books Endpoints", func() {
		ginkgo.It("Create & List", func() {
			cli := sdkapp0.NewBook()

			ginkgo.By("Create Book")
			var book = model.Book{
				ID:          "xxx-xxx",
				Title:       "Harry Potter",
				Author:      "J.K Rowling",
				ISBN:        "xxx",
				Pages:       1024,
				PublishedAt: time.Now(),
			}

			err := cli.Create(&book)
			gomega.Ω(err).ShouldNot(gomega.HaveOccurred())

			ginkgo.By("List Books")
			total, list, err := cli.List()
			gomega.Ω(err).ShouldNot(gomega.HaveOccurred())
			gomega.Ω(total).Should(gomega.BeEquivalentTo(1))
			gomega.Ω(list).Should(gomega.HaveLen(1))
			gomega.Ω(list[0].Title).Should(gomega.Equal("Harry Potter"))
		})
	})
})

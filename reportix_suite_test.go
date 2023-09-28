package reportix_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestReportix(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Reportix Suite")
}

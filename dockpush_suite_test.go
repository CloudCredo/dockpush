package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCdock(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cdock Suite")
}

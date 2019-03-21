package integration_tests

import (
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func TestIntegrationTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IntegrationTests Suite")
}

var _ = BeforeSuite(func() {
	testConfig := viper.New()
	testConfig.SetConfigName("integration")
	testConfig.AddConfigPath("$GOPATH/src/github.com/vulcanize/account_transformers/environments/")
	err := testConfig.ReadInConfig()
	ipc = testConfig.GetString("client.ipcPath")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(ioutil.Discard)
})

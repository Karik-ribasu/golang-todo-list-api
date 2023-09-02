package data_test

import (
	"testing"

	"github.com/Karik-ribasu/golang-todo-list-api/domain/data"
	"github.com/Karik-ribasu/golang-todo-list-api/infra/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var cfg config.Config
var dbManager data.DbManager
var errConn error

func TestData(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Suite")
}

var _ = BeforeSuite(func() {
	cfg, _ := config.ReadConfig()
	dbManager, errConn = data.InitDB(cfg)
})

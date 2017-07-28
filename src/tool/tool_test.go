package tool

import (
	"log"
	"testing"
)

func TestExistFile(t *testing.T) {
	log.Println(ExistFile("config.go"))
	log.Println(ExistFile("config"))
}

func TestConf(t *testing.T) {
	log.Println(Conf("addr"))
	log.Println(Conf("hello"))
}

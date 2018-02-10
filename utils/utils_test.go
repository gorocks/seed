package utils

import (
	"log"
	"testing"
)

func TestCheckGip(t *testing.T) {
	a := CheckGip()
	log.Println(a)
}

func TestDoGipInstall(t *testing.T) {
	//
	DoGipInstall("/Users/luan/go/src/github.com/Guazi-inc/seed/requirements.txt")
}

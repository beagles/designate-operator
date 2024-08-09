package main

import (
	"os"
	template "text/template"
)

// Run against the configuration file template to make sure results are as expected.
type RndcConfValue struct {
	RndcKeyAlgorithm string
	RndcKey          string
	PodAddr          string
	AllowCIDR        string
}

func main() {
	templateFile, err := os.ReadFile("../../../templates/designatebind9/config-named/rndc.conf")
	if err != nil {
		panic(err)
	}
	templateInstance, err := template.New("test").Parse(string(templateFile))
	if err != nil {
		panic(err)
	}

	testData := RndcConfValue{"hmac-sha384", "xxxyyy", "172.28.0.9", "172.28.0.0/24"}

	err = templateInstance.Execute(os.Stdout, testData)
	if err != nil {
		panic(err)
	}

	testData = RndcConfValue{"", "xxxyyy", "172.28.0.9", "172.28.0.0/24"}
	err = templateInstance.Execute(os.Stdout, testData)
	if err != nil {
		panic(err)
	}
}

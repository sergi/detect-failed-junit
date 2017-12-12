package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type testsuites struct {
	XMLName       xml.Name    `xml:"testsuites"`
	TestSuiteList []testsuite `xml:"testsuite"`
}

type testsuite struct {
	XMLName  xml.Name `xml:"testsuite"`
	Name     string   `xml:"name,attr"`
	Failures string   `xml:"failures,attr"`
	Errors   string   `xml:"errors,attr"`
}

// This program returns error if there are any test suites that failed or
// errored in a JUnit test report.
func main() {
	flag.Parse()
	args := flag.Args()
	fmt.Println(args)
	if len(args) != 1 {
		fmt.Println("One argument (file path) required")
		os.Exit(1)
	}

	d := testsuites{}
	xmlContent, _ := ioutil.ReadFile(args[0])
	err := xml.Unmarshal(xmlContent, &d)

	if err != nil {
		fmt.Println("Bad JUnit report format")
		panic(err)
	}

	for _, suite := range d.TestSuiteList {
		if suite.Errors == "0" && suite.Failures == "0" {
			fmt.Println("No errors in JUnit test reports, yay!")
			os.Exit(0)
		} else {
			fmt.Println("There were failures in JUnit test reports: " + suite.Name)
			os.Exit(1)
		}
	}
	fmt.Println("Unknown error")
	os.Exit(1)
}

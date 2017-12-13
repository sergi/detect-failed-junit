package main

import (
	"encoding/xml"
	"errors"
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
	XMLName      xml.Name   `xml:"testsuite"`
	Name         string     `xml:"name,attr"`
	Failures     string     `xml:"failures,attr"`
	Errors       string     `xml:"errors,attr"`
	TestCaseList []testcase `xml:"testcase"`
}

type testcase struct {
	XMLName     xml.Name        `xml:"testcase"`
	Classname   string          `xml:"classname,attr"`
	Name        string          `xml:"name,attr"`
	SkipMessage *skipMessage    `xml:"skipped,omitempty"`
	Failure     *failureMessage `xml:"failure,omitempty"`
}

type skipMessage struct {
	Message string `xml:"message,attr"`
}

type failureMessage struct {
	Message  string `xml:"message,attr"`
	Type     string `xml:"type,attr"`
	Contents string `xml:",chardata"`
}

func CheckForFailedTests(xmlContent []byte) error {
	suites := testsuites{}
	err := xml.Unmarshal(xmlContent, &suites)

	if err != nil {
		err := xml.Unmarshal(xmlContent, &suites.TestSuiteList)
		if err != nil {
			return errors.New("Wrong report format")
		}
	}

	for _, suite := range suites.TestSuiteList {
		if (suite.Errors != "" && suite.Errors != "0") || (suite.Failures != "" && suite.Failures != "0") {
			for _, testCase := range suite.TestCaseList {
				if testCase.SkipMessage == nil && testCase.Failure != nil {
					fmt.Println("Failed Test:", testCase.Name, testCase.Classname)
				}
			}
			return errors.New("There were failures in JUnit test reports: " + suite.Name)
		}
	}
	return nil
}

// This program returns error if there are any test suites that failed or
// errored in a JUnit test report.
func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		fmt.Println("One argument (file path) required")
		os.Exit(1)
	}

	xmlContent, err := ioutil.ReadFile(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	testsFailed := CheckForFailedTests(xmlContent)

	if testsFailed != nil {
		fmt.Println(testsFailed)
		os.Exit(1)
	} else {
		fmt.Println("No errors in JUnit test reports, yay!")
		os.Exit(0)
	}
}

package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/fatih/color"
)

type testsuites struct {
	XMLName       xml.Name    `xml:"testsuites"`
	TestSuiteList []testsuite `xml:"testsuite"`
}

type testsuite struct {
	XMLName      xml.Name   `xml:"testsuite"`
	Name         string     `xml:"name,attr"`
	Failures     int        `xml:"failures,attr"`
	Skipped      int        `xml:"skipped,attr"`
	Errors       int        `xml:"errors,attr"`
	TestCaseList []testcase `xml:"testcase"`
}

type testcase struct {
	XMLName        xml.Name        `xml:"testcase"`
	Classname      string          `xml:"classname,attr"`
	Name           string          `xml:"name,attr"`
	Time           string          `xml:"time,attr"`
	SkipMessage    *skipMessage    `xml:"skipped,omitempty"`
	FailureMessage *failureMessage `xml:"failure,omitempty"`
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
		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		fmt.Printf("%s: %d tests\n", red("ERRORS"), suite.Errors)
		var skipped,failed,passed int
		for _, testCase := range suite.TestCaseList {
			if testCase.SkipMessage != nil {
				fmt.Printf("%s: %s\n", yellow("SKIPPED"), testCase.Name)
				skipped++
			}
			if testCase.FailureMessage != nil {
				fmt.Printf("%s: %s (%s)\n", red("FAILED"), testCase.Name, testCase.FailureMessage)
				failed++
			}
			if testCase.FailureMessage == nil && testCase.SkipMessage == nil {
				fmt.Printf("%s:  %s\n", green("PASSED"), testCase.Name)
				passed++
			}
		}
		fmt.Printf("%s: %d tests\n", red("SKIPPED"), skipped)
		fmt.Printf("%s: %d tests\n", red("FAILED"), failed)
		fmt.Printf("%s: %d tests\n", green("PASSED"), passed)
		if (suite.Errors > 0) || (suite.Failures > 0) {
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
	fmt.Println("====== JUnit Test Validator ======")
	fmt.Println("Checking '" + args[0] + "' for failed/errored tests...")
	testsFailed := CheckForFailedTests(xmlContent)

	if testsFailed != nil {
		fmt.Println(testsFailed)
		os.Exit(1)
	} else {
		fmt.Println("No errors, yay!")
		os.Exit(0)
	}
}

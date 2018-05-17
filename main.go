package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/kstiehl/detect-failed-junit/parser"
)

func CheckForFailedTests(xmlContent []byte) error {
	suites, err := parser.Unmarshal(xmlContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse xml content %s", err.Error())
	}

	for _, suite := range suites.TestSuiteList {
		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		fmt.Printf("%s: %d tests\n", red("ERRORS"), suite.Errors)
		var skipped, failed, passed int
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

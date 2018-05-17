package parser

import (
	"encoding/xml"
)

// Testsuites xml element
type Testsuites struct {
	XMLName       xml.Name    `xml:"testsuites"`
	TestSuiteList []Testsuite `xml:"testsuite"`
}

// Testsuite xml element
type Testsuite struct {
	XMLName      xml.Name   `xml:"testsuite"`
	Name         string     `xml:"name,attr"`
	Failures     int        `xml:"failures,attr"`
	Skipped      int        `xml:"skipped,attr"`
	Errors       int        `xml:"errors,attr"`
	TestCaseList []Testcase `xml:"testcase"`
}

// Testcase xml element
type Testcase struct {
	XMLName        xml.Name        `xml:"testcase"`
	Classname      string          `xml:"classname,attr"`
	Name           string          `xml:"name,attr"`
	Time           string          `xml:"time,attr"`
	SkipMessage    *SkipMessage    `xml:"skipped,omitempty"`
	FailureMessage *FailureMessage `xml:"failure,omitempty"`
}

// SkipMessage xml element
type SkipMessage struct {
	Message string `xml:"message,attr"`
}

// FailureMessage xml element
type FailureMessage struct {
	Message  string `xml:"message,attr"`
	Type     string `xml:"type,attr"`
	Contents string `xml:",chardata"`
}

// Unmarshal testsuites
func Unmarshal(xmlContent []byte) (Testsuites, error) {
	testSuite := Testsuites{}
	var err error
	err = xml.Unmarshal(xmlContent, &testSuite)
	if err != nil {
		err = xml.Unmarshal(xmlContent, &testSuite.TestSuiteList)
		return testSuite, err
	}
	return testSuite, err
}

// Marshal testsuites
func Marshal(suites Testsuites) ([]byte, error) {
	return xml.Marshal(suites)
}

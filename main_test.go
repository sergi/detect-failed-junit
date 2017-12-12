package main

import (
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
)

func TestReportParsing(t *testing.T) {
	xmlContent, _ := ioutil.ReadFile("fixtures/report-bad2.xml")
	err1 := CheckForFailedTests(xmlContent)
	assert.EqualError(t, err1,  "There were failures in JUnit test reports: gda.device.detector.XHDetectorTest")

	xmlContent, _ = ioutil.ReadFile("fixtures/report-good.xml")
	err2 := CheckForFailedTests(xmlContent)
	assert.Equal(t, nil, err2)

	xmlContent, _ = ioutil.ReadFile("fixtures/report-bad.xml")
	err3 := CheckForFailedTests(xmlContent)
	assert.EqualError(t, err3,  "There were failures in JUnit test reports: org.sergi.test2")
}

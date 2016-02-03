package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrintSaveReport(t *testing.T) {
	expected := "10 items inserted.\n"

	saveReport := &SaveReport{}
	saveReport.NumItemsSaved = 10
	actual := getReport(saveReport)
	assert.Equal(t, expected, actual)
}

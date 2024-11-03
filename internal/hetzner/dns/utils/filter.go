// Package utils contains the functionality for filtering DNS records
package utils

import (
	"h3s/internal/hetzner/dns/api"
)

// isRecordExpected checks if a record is one of the expected records
func isRecordExpected(record api.Record, expectedRecords []api.CreateRecordOpts) bool {
	for _, expectedRecord := range expectedRecords {
		if record.Name == expectedRecord.Name && record.Type == expectedRecord.Type {
			return true
		}
	}
	return false
}

// FilterFoundRecords builds a list of expected records and filters the found records to the expected records
func FilterFoundRecords(foundRecords []api.Record) []api.Record {
	expectedRecords := GetExpectedRecords(nil, nil)

	var filteredRecords []api.Record
	for _, foundRecord := range foundRecords {
		if isRecordExpected(foundRecord, expectedRecords) {
			filteredRecords = append(filteredRecords, foundRecord)
		}
	}
	return filteredRecords
}

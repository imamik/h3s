package utils

import (
	"h3s/internal/hetzner/dns/api"
)

func isRecordExpected(record api.Record, expectedRecords []api.CreateRecordOpts) bool {
	for _, expectedRecord := range expectedRecords {
		if record.Name == expectedRecord.Name && record.Type == expectedRecord.Type {
			return true
		}
	}
	return false
}

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

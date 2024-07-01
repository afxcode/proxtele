package main

import (
	"regexp"
)

type BackupResult struct {
	VMID     string
	Name     string
	Status   string
	Time     string
	Size     string
	Filename string
}

func parseMessages(lines []string) (results []BackupResult, err error) {
	for _, line := range lines {
		re := regexp.MustCompile(`(\d+)\s+(\S+)\s+(\S+)\s+(\S+\s+\d+s)\s+(\d+\.\d+\s+\S+)\s+(\S+)`)
		matches := re.FindStringSubmatch(line)
		if matches != nil && len(matches) >= 7 {
			var result BackupResult
			result.VMID = matches[1]
			result.Name = matches[2]
			result.Status = matches[3]
			result.Time = matches[4]
			result.Size = matches[5]
			result.Filename = matches[6]
			results = append(results, result)
		}
	}
	return
}

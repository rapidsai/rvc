package rvc

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	ucxPyVersionPattern  = "^v?0*\\.[0-9]{1,2}(\\.[0-9]+)?$"
	rapidsVersionPattern = "^v?[0-9]{2}\\.[0-9]{2}(\\.[0-9]+)?$"
)

var versionMapping = map[string]string{
	"21.06": "0.20",
	"21.08": "0.21",
	"21.10": "0.22",
	"21.12": "0.23",
	"22.02": "0.24",
	"22.04": "0.25",
	"22.06": "0.26",
	"22.08": "0.27",
	"22.10": "0.28",
	"22.12": "0.29",
	"23.02": "0.30",
	"23.04": "0.31",
	"23.06": "0.32",
	"23.08": "0.33",
	"23.10": "0.34",
	"23.12": "0.35",
	"24.02": "0.36",
	"24.04": "0.37",
	"24.06": "0.38",
	"24.08": "0.39",
	"24.10": "0.40",
	"24.12": "0.41",
	"25.02": "0.42",
	"25.04": "0.43",
	"25.06": "0.44",
	"25.08": "0.45",
	"25.10": "0.46",
	"25.12": "0.47",
}

func GetUcxPyFromRapids(version string) (string, error) {
	parsedVersion, err := parseRapidsVersion(version)
	if err != nil {
		return "", err
	}

	if ucxPyVersion, ok := versionMapping[parsedVersion]; ok {
		return ucxPyVersion, nil
	}

	return "", fmt.Errorf("Ucx-py version not found for RAPIDS version %v", version)
}

func GetRapidsFromUcxPy(version string) (string, error) {
	parsedVersion, err := parseUcxPyVersion(version)
	if err != nil {
		return "", err
	}

	for rapidsVersion, ucxPyVersion := range versionMapping {
		if ucxPyVersion == parsedVersion {
			return rapidsVersion, nil
		}
	}

	return "", fmt.Errorf("RAPIDS version not found for ucx-py version %v", version)
}

func parseRapidsVersion(version string) (string, error) {
	rapidsRegex := regexp.MustCompile(rapidsVersionPattern)
	if !rapidsRegex.MatchString(version) {
		return "", fmt.Errorf("Version \"%v\" is not a RAPIDS version", version)
	}

	if version[0] == 'v' {
		version = version[1:]
	}

	split := strings.Split(version, ".")

	year, _ := strconv.Atoi(split[0])
	if year < 21 {
		return "", fmt.Errorf("Year value \"%v\" must be greater than 21 for RAPIDS version \"%v\"", year, version)
	}

	month, _ := strconv.Atoi(split[1])
	if month < 1 || month > 12 {
		return "", fmt.Errorf("Month value \"%v\" must be between 1 and 12 for RAPIDS version \"%v\"", month, version)
	}

	return fmt.Sprintf("%s.%s", split[0], split[1]), nil
}

func parseUcxPyVersion(version string) (string, error) {
	ucxPyRegex := regexp.MustCompile(ucxPyVersionPattern)
	if !ucxPyRegex.MatchString(version) {
		return "", fmt.Errorf("Version \"%v\" is not a ucx-py version", version)
	}

	if version[0] == 'v' {
		version = version[1:]
	}

	split := strings.Split(version, ".")

	minor, _ := strconv.Atoi(split[1])
	if minor == 0 {
		return "", fmt.Errorf("Minor value \"%v\" must not be 0 for ucx-py version \"%v\"", minor, version)
	}

	return fmt.Sprintf("%s.%s", split[0], split[1]), nil
}

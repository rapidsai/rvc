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
	if len(split) > 3 || len(split) < 2 {
		return "", fmt.Errorf("Unexpected number of \".\" in RAPIDS version \"%v\"", version)
	}

	year, err := strconv.Atoi(split[0])
	if err != nil {
		return "", fmt.Errorf("Non numerical year value \"%v\" for RAPIDS version \"%v\"", split[0], version)
	}

	if year < 21 {
		return "", fmt.Errorf("Year value \"%v\" must be greater than 21 for RAPIDS version \"%v\"", year, version)
	}

	month, err := strconv.Atoi(split[1])
	if err != nil {
		return "", fmt.Errorf("Non numerical month value \"%v\" for RAPIDS version \"%v\"", split[1], version)
	}

	if month < 1 || month > 12 {
		return "", fmt.Errorf("Month value \"%v\" must be between 1 and 12 for RAPIDS version \"%v\"", month, version)
	}

	if len(split) > 2 {
		_, err := strconv.Atoi(split[2])
		if err != nil {
			return "", fmt.Errorf("Non numerical patch value \"%v\" for RAPIDS version \"%v\"", split[2], version)
		}
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
	if len(split) > 3 || len(split) < 2 {
		return "", fmt.Errorf("Unexpected number of \".\" in ucx-py version \"%v\"", version)
	}

	major, err := strconv.Atoi(split[0])
	if err != nil {
		return "", fmt.Errorf("Non numerical major value \"%v\" for ucx-py version \"%v\"", split[0], version)
	}

	if major != 0 {
		return "", fmt.Errorf("Major value \"%v\" must be equal to 0 for ucx-py version \"%v\"", major, version)
	}

	_, err = strconv.Atoi(split[1])
	if err != nil {
		return "", fmt.Errorf("Non numerical minor value \"%v\" for ucx-py version \"%v\"", split[1], version)
	}

	if len(split) > 2 {
		_, err := strconv.Atoi(split[2])
		if err != nil {
			return "", fmt.Errorf("Non numerical patch value \"%v\" for ucx-py version \"%v\"", split[2], version)
		}
	}

	return fmt.Sprintf("%s.%s", split[0], split[1]), nil
}

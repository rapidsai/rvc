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
	rapidsStartYear      = 21
	rapidsStartMonth     = 6
	ucxStart             = 20
)

func GetUcxPyFromRapids(version string) (string, error) {
	year, month, err := parseRapidsVersion(version)
	if err != nil {
		return "", err
	}

	ucxVer := (((12 * (year - rapidsStartYear)) + month - rapidsStartMonth) / 2) + ucxStart

	return fmt.Sprintf("0.%d", ucxVer), nil
}

func GetRapidsFromUcxPy(version string) (string, error) {
	parsedVersion, err := parseUcxPyVersion(version)
	if err != nil {
		return "", err
	}
	distance := parsedVersion - ucxStart
	year := int64((distance*2.0+rapidsStartMonth)/12) + rapidsStartYear
	month := (distance*2 + rapidsStartMonth) % 12
	if month == 0 {
		year = year - 1
		month = 12
	}

	return fmt.Sprintf("%d.%02d", year, month), nil
}

func parseRapidsVersion(version string) (int, int, error) {
	rapidsRegex := regexp.MustCompile(rapidsVersionPattern)
	if !rapidsRegex.MatchString(version) {
		return -1, -1, fmt.Errorf("Version \"%v\" is not a RAPIDS version", version)
	}

	if version[0] == 'v' {
		version = version[1:]
	}

	split := strings.Split(version, ".")

	year, _ := strconv.Atoi(split[0])
	if year < 21 {
		return -1, -1, fmt.Errorf("Year value \"%v\" must be greater than 21 for RAPIDS version \"%v\"", year, version)
	}

	month, _ := strconv.Atoi(split[1])
	if month < 1 || month > 12 {
		return -1, -1, fmt.Errorf("Month value \"%v\" must be between 1 and 12 for RAPIDS version \"%v\"", month, version)
	}

	if month%2 != 0 {
		return -1, -1, fmt.Errorf("Month value \"%v\" must be even for RAPIDS version \"%v\"", month, version)
	}

	return year, month, nil
}

func parseUcxPyVersion(version string) (int, error) {
	ucxPyRegex := regexp.MustCompile(ucxPyVersionPattern)
	if !ucxPyRegex.MatchString(version) {
		return -1, fmt.Errorf("Version \"%v\" is not a ucx-py version", version)
	}

	if version[0] == 'v' {
		version = version[1:]
	}

	split := strings.Split(version, ".")

	minor, _ := strconv.Atoi(split[1])
	if minor < ucxStart {
		return -1, fmt.Errorf("Minor value \"%v\" must be >=%v for ucx-py version \"%v\"", minor, ucxStart, version)
	}

	return minor, nil
}

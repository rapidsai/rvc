package rvc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseRapidsVersion(t *testing.T) {

	validVersions := []string{"v22.02", "v22.02.00", "v22.02.01", "22.02", "22.02.00", "22.02.01"}
	expectedYear := 22
	expectedMonth := 2

	for _, version := range validVersions {
		year, month, err := parseRapidsVersion(version)
		require.NoError(t, err)
		require.Equal(t, expectedYear, year)
		require.Equal(t, expectedMonth, month)
	}

	invalidVersions := []string{"0.21", "v22.00", "v22.0", "22.0", "abc", "20.02"}

	for _, version := range invalidVersions {
		_, _, err := parseRapidsVersion(version)
		require.Error(t, err)
	}
}

func TestParseUcxPyVersion(t *testing.T) {

	validVersions := []string{"v0.21", "v0.21.0", "v0.21.1", "0.21", "0.21.0", "0.21.1"}
	expected := 21

	for _, version := range validVersions {
		v, err := parseUcxPyVersion(version)
		require.NoError(t, err)
		require.Equal(t, expected, v)
	}

	invalidVersions := []string{"22.02", "v22.02", "v1.01", "1.01", "12", "abc", "0.0"}

	for _, version := range invalidVersions {
		_, err := parseUcxPyVersion(version)
		require.Error(t, err)
	}
}

func TestGetRapidsFromUcxPy(t *testing.T) {
	validVersions := map[string]string{
		"v0.20":   "21.06",
		"0.20":    "21.06",
		"v0.20.0": "21.06",
		"0.20.0":  "21.06",
		"v0.20.1": "21.06",
		"0.20.1":  "21.06",
		"0.30":    "23.02",
		"0.35":    "23.12",
		"0.36":    "24.02",
		"0.37":    "24.04",
		"0.38":    "24.06",
		"0.39":    "24.08",
		"0.40":    "24.10",
		"0.41":    "24.12",
		"0.45":    "25.08",
	}

	for version, expected := range validVersions {
		v, err := GetRapidsFromUcxPy(version)
		require.NoError(t, err)
		require.Equal(t, expected, v)
	}

	invalidVersions := []string{"v321.02", "abc", "15", "1.15", "v1.24"}

	for _, version := range invalidVersions {
		_, err := GetRapidsFromUcxPy(version)
		require.Error(t, err)
	}
}

func TestGetUcxPyFromRapids(t *testing.T) {
	validVersions := map[string]string{
		"v21.06":   "0.20",
		"21.06":    "0.20",
		"v21.06.0": "0.20",
		"21.06.0":  "0.20",
		"v21.06.1": "0.20",
		"21.06.1":  "0.20",
		"23.02":    "0.30",
		"23.12":    "0.35",
		"24.02":    "0.36",
		"24.04":    "0.37",
		"24.06":    "0.38",
		"24.08":    "0.39",
		"24.10":    "0.40",
		"24.12":    "0.41",
		"25.08":    "0.45",
	}

	for version, expected := range validVersions {
		v, err := GetUcxPyFromRapids(version)
		require.NoError(t, err)
		require.Equal(t, expected, v)
	}

	invalidVersions := []string{"v21.00", "abc", "15", "1.15", "v1.24"}

	for _, version := range invalidVersions {
		_, err := GetUcxPyFromRapids(version)
		require.Error(t, err)
	}
}

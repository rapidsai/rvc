package version

const undefinedVersion string = "undefined"

// Will be set at build time using ldflags
var version = undefinedVersion

func Get() string {
	return version
}

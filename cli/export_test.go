package cli

import "fmt"

func SetVersion(v string) {
	version = v
}

func SetRevision(r string) {
	revision = r
}

func SetVersionString() {
	versionString = fmt.Sprintf("v%s-%s", version, revision)
}

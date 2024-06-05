package buildinfo

import "fmt"

var (
	Version = "0.0.0"
	Name    string
	GitSHA  string
)

type BuildInfo struct {
	Name    string
	Version string
	GitSHA  string
}

func (bi BuildInfo) String() string {
	return fmt.Sprintf("%s version: %s (%s)", bi.Name, bi.Version, bi.GitSHA)
}

func Get() BuildInfo {
	return BuildInfo{
		Version: Version,
		Name:    Name,
		GitSHA:  GitSHA,
	}
}

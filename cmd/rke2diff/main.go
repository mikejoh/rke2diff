package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"

	"github.com/google/go-github/v62/github"
	gversion "github.com/hashicorp/go-version"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mikejoh/rke2diff/internal/buildinfo"
)

type rke2diffOptions struct {
	version      bool
	rke2Versions rkeVersionSlice
	releases     bool
	skipRc       bool
	pick         bool
	perPage      int
}

type GitHubProject struct {
	Repo     string
	Owner    string
	Releases []*github.RepositoryRelease
}

type RKE2Release struct {
	Version    string
	Components []Component
}

type rkeVersionSlice []string

func (s *rkeVersionSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *rkeVersionSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// Icon constants.
const (
	versionParseError = "❌"
	bumped            = "🚀"
)

func main() {
	var rke2diffOpts rke2diffOptions
	flag.BoolVar(&rke2diffOpts.version, "version", false, "Print the version number.")
	flag.BoolVar(&rke2diffOpts.releases, "releases", false, "Show all releases.")
	flag.BoolVar(&rke2diffOpts.skipRc, "skip-rc", true, "Skip release candidate releases.")
	flag.BoolVar(&rke2diffOpts.pick, "pick", false, "Interactive release picker.")
	flag.IntVar(&rke2diffOpts.perPage, "per-page", 100, "Skip release candidate releases.")
	flag.Var(&rke2diffOpts.rke2Versions, "rke2", "RKE2 version to compare, can be set multiple times.")
	flag.Parse()

	if rke2diffOpts.version {
		fmt.Println(buildinfo.Get())
		os.Exit(0)
	}

	if len(rke2diffOpts.rke2Versions) > 2 {
		log.Fatal("only two RKE2 versions can be compared")
	}

	ghClient := github.NewClient(nil)

	project := GitHubProject{
		Owner:    "rancher",
		Repo:     "rke2",
		Releases: []*github.RepositoryRelease{},
	}

	ctx := context.Background()

	fetchedReleases, _, err := ghClient.Repositories.ListReleases(ctx, project.Owner, project.Repo, &github.ListOptions{
		PerPage: rke2diffOpts.perPage,
	})
	if err != nil {
		log.Fatal(err)
	}

	var releases []*github.RepositoryRelease

	for _, release := range fetchedReleases {
		if rke2diffOpts.skipRc && strings.Contains(release.GetTagName(), "rc") {
			continue
		}
		releases = append(releases, release)
	}

	if rke2diffOpts.releases {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleLight)
		header := table.Row{
			"Release",
			"Published At",
		}
		t.AppendHeader(header)
		t.Style().Title.Align = text.AlignCenter

		for _, release := range releases {
			t.AppendRow(table.Row{release.GetTagName(), release.GetPublishedAt()})
		}

		t.Render()

		os.Exit(0)
	}

	if rke2diffOpts.pick {
		var releaseURL string
		_, err := fuzzyfinder.FindMulti(
			releases,
			func(i int) string {
				releaseURL = *releases[i].HTMLURL
				return *releases[i].TagName
			},
		)
		if err != nil {
			log.Fatal(err)
		}

		err = openURL(releaseURL)
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}

	project.Releases = releases

	for _, rke2Version := range rke2diffOpts.rke2Versions {
		release := findRelease(project.Releases, rke2Version)
		if release == nil {
			log.Fatalf("release %s not found", rke2Version)
		}
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	header := table.Row{
		"Name",
	}
	for _, rke2Version := range rke2diffOpts.rke2Versions {
		r := getRelease(project.Releases, rke2Version)
		header = append(header, fmt.Sprintf("%s (%s)", rke2Version, r.GetPublishedAt().Format("2006-01-02")))
	}
	t.AppendHeader(header)
	t.AppendFooter(table.Row{
		fmt.Sprintf("%s = %s\n%s = %s", versionParseError, "Version string invalid", bumped, "Bumped"),
	})
	t.Style().Title.Align = text.AlignCenter

	rows := make(map[string]table.Row)

	componentVersionDiffs := make(map[string]map[string]string)

	// TODO: Compare version.Version instead of string to get proper version comparison
	for _, rke2Version := range rke2diffOpts.rke2Versions {
		componentVersionDiffs[rke2Version] = make(map[string]string)
		for _, release := range project.Releases {
			if !(release.GetTagName() == rke2Version) {
				continue
			}

			htmlBytes := mdToHTML([]byte(release.GetBody()))
			components := parseHTMLTable(string(htmlBytes), "h2", "charts-versions")
			packagedComponents := parseHTMLTable(string(htmlBytes), "h2", "packaged-component-versions")

			if len(components) == 0 {
				log.Fatalf("no components found in release %s", release.GetTagName())
			}

			if len(packagedComponents) == 0 {
				log.Fatalf("no packaged components found in release %s", release.GetTagName())
			}

			rke2Release := RKE2Release{
				Version:    release.GetTagName(),
				Components: append(packagedComponents, components...),
			}

			for _, component := range rke2Release.Components {
				var componentVersion string

				if _, ok := rows[component.Name]; !ok {
					rows[component.Name] = table.Row{
						strings.ToLower(component.Name),
					}
				}

				if _, ok := componentVersionDiffs[rke2Version][component.Name]; !ok {
					componentVersionDiffs[rke2Version][component.Name] = component.Version
				}

				if componentVersionDiffs[rke2diffOpts.rke2Versions[0]][component.Name] != component.Version {
					componentVersion = fmt.Sprintf("%s %s", component.Version, bumped)
				} else {
					componentVersion = component.Version
				}

				_, err := gversion.NewVersion(component.Version)
				if err != nil {
					componentVersion = fmt.Sprintf("%s %s", component.Version, versionParseError)
				}

				rows[component.Name] = append(rows[component.Name], componentVersion)
			}
		}
	}

	keys := make([]string, 0, len(rows))
	for k := range rows {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		t.AppendRow(rows[k])
	}

	t.Render()
}

func findRelease(releases []*github.RepositoryRelease, version string) *github.RepositoryRelease {
	for _, release := range releases {
		if release.GetTagName() == version {
			return release
		}
	}
	return nil
}

func getRelease(releases []*github.RepositoryRelease, version string) *github.RepositoryRelease {
	for _, release := range releases {
		if release.GetTagName() == version {
			return release
		}
	}
	return nil
}

func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

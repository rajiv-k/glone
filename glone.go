package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	usageHeader = "usage: glone [OPTIONS...] <url>"
	usageFooter = `
EXAMPLES:
glone git@github.com:rajiv/glone.git
glone https://github.com/rajiv/glone.git
`
)

var (
	gitSshUrl   = regexp.MustCompile(`^git@([a-zA-Z0-9._-]+):([a-zA-Z0-9./._-]+)(?:\?||$)(.*)$`)
	gitHttpsUrl = regexp.MustCompile(`^https://([a-zA-Z0-9._-]+)/([a-zA-Z0-9./._-]+)(?:\?||$)(.*)$`)
	Version     = "unknown"
	versionFlag = false
)

type RepoInfo struct {
	Scheme   string
	Host     string
	Owner    string
	RepoName string
}

func parse(repoUrl string) (*RepoInfo, error) {
	var scheme string

	var re *regexp.Regexp

	if strings.HasPrefix(repoUrl, "https://") {
		re = gitHttpsUrl
		scheme = "https"
	} else if strings.HasPrefix(repoUrl, "git@") {
		re = gitSshUrl
		scheme = "ssh"
	}

	matches := re.FindAllStringSubmatch(repoUrl, -1)
	if len(matches) < 1 {
		return nil, errors.New("invalid repo url")
	}
	first := matches[0]

	owner := filepath.Dir(first[2])
	repo := strings.TrimSuffix(filepath.Base(first[2]), ".git")

	info := &RepoInfo{
		Scheme:   scheme,
		Host:     first[1],
		Owner:    owner,
		RepoName: repo,
	}

	return info, nil

}

func (r *RepoInfo) String() string {
	if r.Scheme == "https" {
		return fmt.Sprintf("https://%v/%v/%v.git", r.Host, r.Owner, r.RepoName)
	} else if r.Scheme == "ssh" {
		return fmt.Sprintf("git@%v:%v/%v.git", r.Host, r.Owner, r.RepoName)
	}
	return ""
}

func clone(r *RepoInfo) error {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		return errors.New("GOPATH is not set")
	}
	projectDir := path.Join(goPath, "src", r.Host, r.Owner, r.RepoName)

	_, err := os.Stat(projectDir)
	if errors.Is(err, os.ErrExist) {
		return errors.New("directory %v already exists in $GOPATH/src")
	}

	repoUrl := r.String()
	cmd := exec.Command("git", "clone", repoUrl, projectDir)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func usage() {
	var builder strings.Builder
	fmt.Fprintln(&builder, "OPTIONS:")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(&builder, "\t--%v:\t\t%v", f.Name, f.Usage)
	})
	fmt.Printf("%v\n\n%v\n%v\n", usageHeader, builder.String(), usageFooter)
}

func version() string {
	return fmt.Sprintf("glone version %v", Version)
}

func main() {
	flag.BoolVar(&versionFlag, "version", false, "Print version and exit")
	flag.Usage = usage

	flag.Parse()

	if versionFlag {
		fmt.Println(version())
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		usage()
		os.Exit(0)
	}

	repoInfo, err := parse(os.Args[1])
	if err != nil {
		panic(err)
	}

	if err := clone(repoInfo); err != nil {
		fmt.Printf("ERROR: failed to clone: %v\n", err)
	}
}

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

var (
	gitSshUrl   = regexp.MustCompile(`^git@([a-zA-Z0-9._-]+):([a-zA-Z0-9./._-]+)(?:\?||$)(.*)$`)
	gitHttpsUrl = regexp.MustCompile(`^https://([a-zA-Z0-9._-]+)/([a-zA-Z0-9./._-]+)(?:\?||$)(.*)$`)
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
	repo := strings.TrimRight(filepath.Base(first[2]), ".git")

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

func main() {
	flags := flag.NewFlagSet("clone", flag.ContinueOnError)
	if err := flags.Parse(os.Args[1:]); err != nil {
		fmt.Printf("ERROR: could not parse flags: %v\n", err)
		os.Exit(1)
	}
	_ = flags
	repoInfo, err := parse(os.Args[1])
	if err != nil {
		panic(err)
	}

	if err := clone(repoInfo); err != nil {
		fmt.Printf("ERROR: failed to clone: %v\n", err)
	}
}

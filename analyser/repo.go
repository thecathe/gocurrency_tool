package main

import (
	"os"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"github.com/thecathe/gocurrency_tool/analyser/log"
)

// clone a repo and return its path
func CloneRepo(url string) (string, string) {

	var last_commit_hash string

	repo_dir := filepath.Join(clone_dir, ProjectName(url))

	// clone needs empty dir to clone into
	if _, empty_err := os.Stat(repo_dir); os.IsNotExist(empty_err) {
		log.GeneralLog("Starting clone of repo: %s\n", url)
		// Clones the repository into the given dir, just as a normal git clone does
		r, git_err := git.PlainClone(repo_dir, false, &git.CloneOptions{
			URL:      "https://github.com/" + url,
			Progress: os.Stdout,
		})

		// check cloned correctly
		if git_err != nil {
			log.WarningLog("CloneRepo,  %s: Error when cloning...\n\tpath: %s\n\terror: %v\n", url, repo_dir, git_err)
		} else {
			// sucess, pass
			log.GeneralLog("CloneRepo,  %s: Cloning successful\n\n", url)
		}

		head, _ := r.Head()
		if head != nil {
			cIter, err := r.Log(&git.LogOptions{From: head.Hash()})
			if err != nil {
				log.WarningLog("CloneRepo, %s: An error occured when reading the commit log\n\terror: %v\n.", url, err)
			}
			commits := []*object.Commit{}
			err = cIter.ForEach(func(c *object.Commit) error {
				commits = append(commits, c)
				return nil
			})

			last_commit_hash = commits[0].Hash.String()
			if err != nil {
				log.ExitLog(1, "CloneRepo,  %s: Error when iterating through commits.\n")
			}
		} else {
			log.WarningLog("CloneRepo,  %s: Repo was empty?\n")
		}

		// Write the commits of the projects to the commits.csv
		f, err := os.OpenFile("commits.csv", //filepath.Join(GenerateWDPath(), "commits.csv")
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.WarningLog("CloneRepo,  %s: Error when opening \"commits.csv\"\n\terror: %v\n", err)
		}
		defer f.Close()
		if _, err := f.WriteString(url + "," + last_commit_hash + "\n"); err != nil {
			log.WarningLog("CloneRepo,  %s: Error when writing to \"commits.csv\"\n\terror: %v\n", err)
		}
	} else {
		log.WarningLog("CloneRepo,  %s: Dir to clone repo into was not empty...\n\tpath: %s\n", url, repo_dir)
	}

	return repo_dir, last_commit_hash
}

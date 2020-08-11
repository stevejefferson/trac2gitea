package wiki

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// RepoClone clones our wiki repo to the provided directory.
func (accessor *Accessor) RepoClone() {
	isBare := false
	repository, err := git.PlainClone(accessor.repoDir, isBare, &git.CloneOptions{
		URL:               accessor.repoURL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		log.Fatal(err)
	}

	accessor.repo = repository
}

// RepoStageAndCommit stages any files added or updated since the last commit then commits them to our cloned wiki repo.
// We package the staging and commit together here because it is easier than embedding hooks to do the git staging
// deep into the wiki parsing process where files from the Trac worksapce can get copied over on-the-fly.
func (accessor *Accessor) RepoStageAndCommit(author string, authorEMail string, message string) {
	worktree, err := accessor.repo.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	status, err := worktree.Status()
	for file, filestatus := range status {
		worktreeStatus := filestatus.Worktree
		if worktreeStatus == git.Untracked || worktreeStatus == git.Modified {
			_, err = worktree.Add(file)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	fmt.Fprintf(os.Stderr, "\n")

	_, err = worktree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  author + "[trac user]",
			Email: authorEMail,
			When:  time.Now(),
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}

// RepoComplete indicates that changes to the local wiki repository are complete.
// In an ideal world we would push back to the remote here
// however, as I haven't worked out how to do the authentication for that yet, we just output a message telling the user to do it.
func (accessor *Accessor) RepoComplete() {
	fmt.Printf("Trac wiki has been imported into cloned wiki repository at %s. Please review changes and push back to remote when done.\n",
		accessor.repoDir)
}

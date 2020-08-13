package giteawiki

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"stevejefferson.co.uk/trac2gitea/log"
)

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}

	return cerr
}

// CopyFile copies an internal file into the Gitea Wiki.
func (accessor *DefaultAccessor) CopyFile(externalFilePath string, giteaWikiRelPath string) {
	_, err := os.Stat(externalFilePath)
	if os.IsNotExist(err) {
		log.Warnf("cannot copy non-existant file referenced from Wiki: \"%s\"\n", externalFilePath)
		return
	}

	giteaPath := filepath.Join(accessor.repoDir, giteaWikiRelPath)
	err = os.MkdirAll(path.Dir(giteaPath), 0775)
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(giteaPath)
	if os.IsExist(err) {
		return // if gitea path exists we'll just assume we've already created it as part of this run
	}

	err = copyFile(externalFilePath, giteaPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("Copied file %s to wiki path %s\n", externalFilePath, giteaWikiRelPath)
}

package gitea

import (
	"io"
	"os"

	"stevejefferson.co.uk/trac2gitea/log"
)

func (accessor *DefaultAccessor) copyFile(externalFilePath string, giteaPath string) error {
	_, err := os.Stat(externalFilePath)
	if os.IsNotExist(err) {
		log.Warnf("cannot copy non-existant attachment file: \"%s\"\n", externalFilePath)
		return nil
	}

	in, err := os.Open(externalFilePath)
	if err != nil {
		log.Error(err)
		return err
	}
	defer in.Close()

	out, err := os.Create(giteaPath)
	if err != nil {
		log.Error(err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		log.Error(err)
		return err
	}

	err = out.Close()
	if err != nil {
		log.Error(err)
		return err
	}

	log.Debugf("Copied file %s to %s\n", externalFilePath, giteaPath)
	return nil
}

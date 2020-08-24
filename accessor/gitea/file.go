// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"io"
	"os"

	"stevejefferson.co.uk/trac2gitea/log"
)

func (accessor *DefaultAccessor) copyFile(externalFilePath string, giteaPath string) error {
	_, err := os.Stat(externalFilePath)
	if os.IsNotExist(err) {
		log.Warn("Cannot copy non-existant attachment file: \"%s\"\n", externalFilePath)
		return nil
	}

	in, err := os.Open(externalFilePath)
	if err != nil {
		log.Error("Cannot open %s: %v\n", externalFilePath, err)
		return err
	}
	defer in.Close()

	out, err := os.Create(giteaPath)
	if err != nil {
		log.Error("Cannot create %s: %v\n", giteaPath, err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		log.Error("Failure trying to copy %s to %s: %v\n", externalFilePath, giteaPath, err)
		return err
	}

	err = out.Close()
	if err != nil {
		log.Error("Cannot close file %s: %v\n", giteaPath, err)
		return err
	}

	log.Debug("Copied file %s to %s\n", externalFilePath, giteaPath)
	return nil
}

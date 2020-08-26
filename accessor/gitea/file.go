// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"io"
	"os"

	"github.com/stevejefferson/trac2gitea/log"
)

func (accessor *DefaultAccessor) copyFile(externalFilePath string, giteaPath string) error {
	_, err := os.Stat(externalFilePath)
	if os.IsNotExist(err) {
		log.Warn("cannot copy non-existant attachment file: \"%s\"", externalFilePath)
		return nil
	}

	in, err := os.Open(externalFilePath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(giteaPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	err = out.Close()
	if err != nil {
		return err
	}

	log.Debug("copied file %s to %s", externalFilePath, giteaPath)
	return nil
}

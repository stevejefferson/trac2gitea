// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"github.com/stevejefferson/trac2gitea/log"
)

// CommitImport commits the import transaction
func (importer *Importer) CommitImport() error {
	log.Info("committing transaction")
	return importer.giteaAccessor.CommitTransaction()
}

// RollbackImport rolls back the import transaction.
func (importer *Importer) RollbackImport() error {
	log.Info("rolling back transaction")
	return importer.giteaAccessor.RollbackTransaction()
}

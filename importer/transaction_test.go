// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer_test

import (
	"testing"
)

func TestCommitImport(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	mockGiteaAccessor.
		EXPECT().
		CommitTransaction().
		Return(nil)

	dataImporter.CommitImport()
}

func TestRollbackImport(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	mockGiteaAccessor.
		EXPECT().
		RollbackTransaction().
		Return(nil)

	dataImporter.RollbackImport()
}

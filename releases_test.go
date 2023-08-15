// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package ghutils_test

import (
	"testing"

	"github.com/google/go-github/v53/github"
	"github.com/st3fan/ghutils"
)

func Test_EachRepoReleases(t *testing.T) {
	client, ctx := mustCreateTestClient()

	count := 0
	err := ghutils.IterateRepoReleases(client, ctx, "boltdb", "bolt", nil, func(commit *github.RepositoryRelease) (bool, error) {
		count += 1
		return true, nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if count != 5 {
		t.Errorf("got %d, expected %d", count, 5)
	}
}

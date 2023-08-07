// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package ghutils_test

import (
	"testing"

	"github.com/google/go-github/v53/github"
	"github.com/st3fan/ghutils"
)

func Test_EachRepoCommits(t *testing.T) {
	client, ctx := mustCreateTestClient()

	count := 0
	err := ghutils.IterateRepoCommits(client, ctx, "golang", "go", nil, func(commit *github.RepositoryCommit) (bool, error) {
		count += 1
		return count < 47, nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if count != 47 {
		t.Errorf("got %d, expected %d", count, 47)
	}
}

func Test_MapRepoCommits(t *testing.T) {
	client, ctx := mustCreateTestClient()

	count := 0
	hashes, _ := ghutils.MapRepoCommits(client, ctx, "golang", "go", nil, func(commit *github.RepositoryCommit) (string, bool, error) {
		count += 1
		return commit.GetSHA(), count < 47, nil
	})

	if len(hashes) != 47 {
		t.Errorf("got %d, wanted %d", len(hashes), 47)
	}
}

func Test_ReduceRepoCommits(t *testing.T) {
	client, ctx := mustCreateTestClient()

	count := 0
	total, _ := ghutils.ReduceRepoCommits(client, ctx, "golang", "go", nil, 0, func(acc int, commit *github.RepositoryCommit) (int, bool, error) {
		count += 1
		return acc + 1, count < 47, nil
	})

	if total != 47 {
		t.Errorf("got %d, wanted %d", total, 47)
	}
}

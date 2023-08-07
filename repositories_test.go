// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package ghutils_test

import (
	"testing"

	"github.com/google/go-github/v53/github"
	"github.com/st3fan/ghutils"
)

func Test_IterateOrgRepos(t *testing.T) {
	client, ctx := mustCreateTestClient()

	count := 0
	err := ghutils.IterateOrgRepos(ctx, client, "golang", nil, func(repo *github.Repository) (bool, error) {
		count += 1
		return true, nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if count != 58 { // TODO Depends on variable number of repos under golang
		t.Errorf("got %d, wanted %d", count, 12)
	}
}

func Test_MapOrgRepos(t *testing.T) {
	client, ctx := mustCreateTestClient()

	links, err := ghutils.MapOrgRepos(ctx, client, "golang", nil, func(repo *github.Repository) (string, bool, error) {
		return repo.GetHTMLURL(), true, nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(links) != 58 { // TODO Depends on variable number of repos under golang
		t.Errorf("got %d, wanted %d", len(links), 12)
	}
}

func Test_ReduceOrgRepos(t *testing.T) {
	client, ctx := mustCreateTestClient()

	total, err := ghutils.ReduceOrgRepos(ctx, client, "golang", nil, 0, func(total int, repo *github.Repository) (int, bool, error) {
		return total + 1, true, nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if total != 58 { // TODO Depends on variable number of repos under golang
		t.Errorf("got %d, wanted %d", total, 12)
	}
}

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package ghutils_test

import (
	"testing"

	"github.com/google/go-github/v53/github"
	"github.com/st3fan/ghutils"
)

func Test_EachPullRequest(t *testing.T) {
	client, ctx := mustCreateTestClient()

	count := 0
	err := ghutils.IterateRepoPullRequests(ctx, client, "golang", "go", nil, func(pull *github.PullRequest) (bool, error) {
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

func Test_MapRepoPullRequests(t *testing.T) {
	client, ctx := mustCreateTestClient()

	count := 0
	links, err := ghutils.MapRepoPullRequests(ctx, client, "golang", "go", nil, func(pull *github.PullRequest) (string, bool, error) {
		count += 1
		return pull.GetHTMLURL(), count < 47, nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(links) != 47 {
		t.Errorf("got %d, expected %d", len(links), 47)
	}
}

func Test_ReduceRepoPullRequests(t *testing.T) {
	client, ctx := mustCreateTestClient()

	count := 0
	total, err := ghutils.ReduceRepoPullRequests(ctx, client, "golang", "go", nil, 0, func(total int, repo *github.PullRequest) (int, bool, error) {
		count += 1
		return total + 1, count < 47, nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if total != 47 {
		t.Errorf("got %d, wanted %d", total, 47)
	}
}

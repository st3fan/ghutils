// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package ghutils

import (
	"context"

	"github.com/google/go-github/v53/github"
)

type IterateIssueFunc func(commit *github.Issue) (bool, error)

func IterateRepoIssues(ctx context.Context, client *github.Client, owner string, repo string, opts *github.IssueListByRepoOptions, f IterateIssueFunc) error {
	if opts == nil {
		opts = &github.IssueListByRepoOptions{}
	}

	for {
		issues, res, err := client.Issues.ListByRepo(ctx, owner, repo, opts)
		if err != nil {
			return err
		}

		for _, issue := range issues {
			cont, err := f(issue)
			if err != nil {
				return err
			}
			if !cont {
				return nil
			}
		}

		if res.NextPage == 0 {
			break
		}

		opts.Page = res.NextPage
	}
	return nil
}

//

type MapIssueFunc[T any] func(commit *github.Issue) (T, bool, error)

func MapRepoIssues[T any](ctx context.Context, client *github.Client, owner string, repo string, opts *github.IssueListByRepoOptions, mapper MapIssueFunc[T]) ([]T, error) {
	var results []T
	err := IterateRepoIssues(ctx, client, owner, repo, opts, func(issue *github.Issue) (bool, error) {
		res, cont, err := mapper(issue)
		if err != nil {
			return false, err
		}
		results = append(results, res)
		return cont, nil
	})
	return results, err
}

//

type ReduceIssueFunc[A any] func(acc A, repo *github.Issue) (A, bool, error)

func ReduceRepoIssues[A any](ctx context.Context, client *github.Client, org string, repo string, opts *github.IssueListByRepoOptions, initial A, reducer ReduceIssueFunc[A]) (A, error) {
	result := initial
	err := IterateRepoIssues(ctx, client, org, repo, opts, func(issue *github.Issue) (bool, error) {
		res, cont, err := reducer(result, issue)
		if err != nil {
			return false, err
		}
		result = res
		return cont, nil
	})
	return result, err
}

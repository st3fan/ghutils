// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package ghutils

import (
	"context"

	"github.com/google/go-github/v53/github"
)

type IterateCommitFunc func(commit *github.RepositoryCommit) (bool, error)

func IterateRepoCommits(client *github.Client, ctx context.Context, owner string, repo string, opts *github.CommitsListOptions, f IterateCommitFunc) error {
	if opts == nil {
		opts = &github.CommitsListOptions{}
	}

	for {
		commits, res, err := client.Repositories.ListCommits(ctx, owner, repo, opts)
		if err != nil {
			return err
		}

		for _, commit := range commits {
			cont, err := f(commit)
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

type MapCommitFunc[T any] func(commit *github.RepositoryCommit) (T, bool, error)

func MapRepoCommits[T any](client *github.Client, ctx context.Context, owner string, repo string, opts *github.CommitsListOptions, mapper MapCommitFunc[T]) ([]T, error) {
	var results []T
	err := IterateRepoCommits(client, ctx, owner, repo, opts, func(commit *github.RepositoryCommit) (bool, error) {
		res, cont, err := mapper(commit)
		if err != nil {
			return false, err
		}
		results = append(results, res)
		return cont, nil
	})
	return results, err
}

//

type ReduceCommitFunc[A any] func(acc A, commit *github.RepositoryCommit) (A, bool, error)

func ReduceRepoCommits[A any](client *github.Client, ctx context.Context, owner string, repo string, opts *github.CommitsListOptions, initial A, reducer ReduceCommitFunc[A]) (A, error) {
	result := initial
	err := IterateRepoCommits(client, ctx, owner, repo, opts, func(commit *github.RepositoryCommit) (bool, error) {
		res, cont, err := reducer(result, commit)
		if err != nil {
			return false, err
		}
		result = res
		return cont, nil
	})
	return result, err
}

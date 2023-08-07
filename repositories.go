// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package ghutils

import (
	"context"
	"fmt"

	"github.com/google/go-github/v53/github"
)

type IterRepoFunc func(commit *github.Repository) (bool, error)

func IterateOrgRepos(ctx context.Context, client *github.Client, org string, opts *github.RepositoryListByOrgOptions, f IterRepoFunc) error {
	if opts == nil {
		opts = &github.RepositoryListByOrgOptions{}
	}

	for {
		commits, res, err := client.Repositories.ListByOrg(ctx, org, opts)
		if err != nil {
			return err
		}

		for _, commit := range commits {
			if cont, err := f(commit); !cont || err != nil {
				return err
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

type MapRepoFunc[T any] func(repo *github.Repository) (T, bool, error)

func MapRepos[T any](ctx context.Context, client *github.Client, repos []string, mapper MapRepoFunc[T]) ([]T, error) {
	return nil, nil
}

func MapOrgRepos[T any](ctx context.Context, client *github.Client, org string, opts *github.RepositoryListByOrgOptions, mapper MapRepoFunc[T]) ([]T, error) {
	var results []T
	err := IterateOrgRepos(ctx, client, org, opts, func(commit *github.Repository) (bool, error) {
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

type ReduceRepoFunc[A any] func(acc A, repo *github.Repository) (A, bool, error)

func ReduceRepos[A any](ctx context.Context, client *github.Client, repos []string, initial A, reducer ReduceRepoFunc[A]) (A, error) {
	result := initial

	for i := range repos {
		owner, name, err := SplitFullRepoName(repos[i])
		if err != nil {
			return initial, fmt.Errorf("could not split full repo name <%s>: %w", repos[i], err)
		}

		repo, _, err := client.Repositories.Get(ctx, owner, name)
		if err != nil {
			return initial, fmt.Errorf("could not get repo <%s>: %w", repos[i], err)
		}

		res, cont, err := reducer(result, repo)
		if !cont || err != nil {
			return initial, err
		}

		result = res // TODO Confusing naming
	}

	return result, nil
}

func ReduceOrgRepos[A any](ctx context.Context, client *github.Client, org string, opts *github.RepositoryListByOrgOptions, initial A, reducer ReduceRepoFunc[A]) (A, error) {
	result := initial
	err := IterateOrgRepos(ctx, client, org, opts, func(commit *github.Repository) (bool, error) {
		res, cont, err := reducer(result, commit)
		if err != nil {
			return false, err
		}
		result = res
		return cont, nil
	})
	return result, err
}

//

func ListAllReposByOrg(ctx context.Context, client *github.Client, org string, opts *github.RepositoryListByOrgOptions) ([]*github.Repository, error) {
	return MapOrgRepos[*github.Repository](ctx, client, org, opts, func(repo *github.Repository) (*github.Repository, bool, error) {
		return repo, true, nil
	})
}

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package ghutils

import (
	"context"

	"github.com/google/go-github/v53/github"
)

type IteratePullRequestFunc func(commit *github.PullRequest) (bool, error)

func IterateRepoPullRequests(ctx context.Context, client *github.Client, owner string, repo string, opts *github.PullRequestListOptions, f IteratePullRequestFunc) error {
	if opts == nil {
		opts = &github.PullRequestListOptions{}
	}

	for {
		pulls, res, err := client.PullRequests.List(ctx, owner, repo, opts)
		if err != nil {
			return err
		}

		for _, pull := range pulls {
			cont, err := f(pull)
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

type MapPullRequestFunc[T any] func(commit *github.PullRequest) (T, bool, error)

func MapRepoPullRequests[T any](ctx context.Context, client *github.Client, owner string, repo string, opts *github.PullRequestListOptions, mapper MapPullRequestFunc[T]) ([]T, error) {
	var results []T
	err := IterateRepoPullRequests(ctx, client, owner, repo, opts, func(pull *github.PullRequest) (bool, error) {
		res, cont, err := mapper(pull)
		if err != nil {
			return false, err
		}
		results = append(results, res)
		return cont, nil
	})
	return results, err
}

//

type ReducePullRequestFunc[A any] func(acc A, repo *github.PullRequest) (A, bool, error)

func ReduceRepoPullRequests[A any](ctx context.Context, client *github.Client, org string, repo string, opts *github.PullRequestListOptions, initial A, reducer ReducePullRequestFunc[A]) (A, error) {
	result := initial
	err := IterateRepoPullRequests(ctx, client, org, repo, opts, func(pull *github.PullRequest) (bool, error) {
		res, cont, err := reducer(result, pull)
		if err != nil {
			return false, err
		}
		result = res
		return cont, nil
	})
	return result, err
}

//

// func ListAllPullsByOrg(ctx context.Context, client *github.Client, org string, opts *github.RepositoryListByOrgOptions) ([]*github.Repository, error) {
// 	return MapRepoPullRequests[*github.PullRequest](ctx, client, org, opts, func(pull *github.PullRequest) (*github.PullRequest, bool, error) {
// 		return pull, true, nil
// 	})
// }

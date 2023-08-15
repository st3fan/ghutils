// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package ghutils

import (
	"context"

	"github.com/google/go-github/v53/github"
)

type IterateReleaseFunc func(release *github.RepositoryRelease) (bool, error)

func IterateRepoReleases(client *github.Client, ctx context.Context, owner string, repo string, opts *github.ListOptions, f IterateReleaseFunc) error {
	if opts == nil {
		opts = &github.ListOptions{}
	}

	for {
		releases, res, err := client.Repositories.ListReleases(ctx, owner, repo, opts)
		if err != nil {
			return err
		}

		for _, release := range releases {
			cont, err := f(release)
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


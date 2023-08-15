// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package ghutils

import (
	"context"
	"net/url"
	"strings"

	"github.com/google/go-github/v47/github"
)

type IterateWorkflowRunsFunc func(run *github.WorkflowRun) (bool, error)

func IterateRepositoryWorkflowRuns(ctx context.Context, client *github.Client, owner string, repo string, opts *github.ListWorkflowRunsOptions, f IterateWorkflowRunsFunc) error {
	if opts == nil {
		opts = &github.ListWorkflowRunsOptions{}
	}

	for {
		runs, res, err := client.Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, opts)
		if err != nil {
			return err
		}

		for _, run := range runs.WorkflowRuns {
			cont, err := f(run)
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

func repoAndOwnerFromHTMLURL(rawURL string) (string, string) {
	u, _ := url.Parse(rawURL)
	fields := strings.Split(u.Path, "/")
	return fields[1], fields[2]
}

func IterateWorkflowRunsForRepositoryRelease(ctx context.Context, client *github.Client, release *github.RepositoryRelease, opts *github.ListWorkflowRunsOptions, f IterateWorkflowRunsFunc) error {
	if opts == nil {
		opts = &github.ListWorkflowRunsOptions{}
	}

	// Bit of a hack but saves us an API call or an extra function parameter
	owner, repo := repoAndOwnerFromHTMLURL(release.GetHTMLURL())

	ref, _, err := client.Git.GetRef(ctx, owner, repo, "tags/"+release.GetTagName())
	if err != nil {
		return err
	}

	opts.HeadSHA = ref.GetObject().GetSHA()

	return IterateRepositoryWorkflowRuns(ctx, client, owner, repo, opts, f)
}

//

type IterateWorkflowJobFunc func(job *github.WorkflowJob) (bool, error)

func IterateWorkflowJobsForWorkflowRun(ctx context.Context, client *github.Client, run *github.WorkflowRun, opts *github.ListWorkflowJobsOptions, f IterateWorkflowJobFunc) error {
	if opts == nil {
		opts = &github.ListWorkflowJobsOptions{}
	}

	for {
		jobs, res, err := client.Actions.ListWorkflowJobs(ctx, run.GetRepository().GetOwner().GetLogin(), run.GetRepository().GetName(), run.GetID(), opts)
		if err != nil {
			return err
		}

		for _, job := range jobs.Jobs {
			cont, err := f(job)
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

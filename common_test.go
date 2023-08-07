// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package ghutils_test

import (
	"context"
	"log"

	"github.com/google/go-github/v53/github"
	"github.com/st3fan/ghutils"
)

// TODO Move into common test functions?
func mustCreateTestClient() (*github.Client, context.Context) {
	token, err := ghutils.NewTokenFromEnvironment()
	if err != nil {
		log.Fatalln("could not find GITHUB_TOKEN: ", err)
	}

	ctx := context.Background()

	client, err := ghutils.NewClient(ctx, token)

	if err != nil {
		log.Fatalln("could not create a GitHub client: ", err)
	}

	return client, ctx
}

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package ghutils

import (
	"context"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

func NewClient(ctx context.Context, token string) (*github.Client, error) {
	return github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))), nil
}

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package ghutils

import (
	"errors"
	"os"
	"strings"

	"github.com/dvcrn/go-1password-cli/op"
)

func NewTokenFromEnvironment() (string, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return "", errors.New("no GITHUB_TOKEN set")
	}

	if strings.HasPrefix(token, "op://") {
		client := op.NewOpClient()
		return client.Read(token)
	}

	return token, nil
}

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package ghutils

import "strings"

func SplitFullRepoName(fullName string) (string, string, error) {
	fields := strings.Split(fullName, "/")
	return fields[0], fields[1], nil // TODO Check
}

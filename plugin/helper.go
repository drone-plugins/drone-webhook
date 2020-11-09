// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import "fmt"

// intInSlice checks if int is in slice of ints
func intInSlice(s []string, e int) bool {
	eStr := fmt.Sprintf("%d", e)
	for _, a := range s {
		if a == eStr {
			return true
		}
	}

	return false
}

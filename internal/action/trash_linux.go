//go:build linux

package action

import trash "github.com/hymkor/trash-go"

func trashOS(path string) error {
	return trash.Throw(path)
}

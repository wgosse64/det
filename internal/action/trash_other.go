//go:build !darwin && !linux

package action

import "fmt"

func trashOS(path string) error {
	return fmt.Errorf("trash not supported on this platform")
}

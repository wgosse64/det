package action

import "github.com/atotto/clipboard"

func CopyToClipboard(s string) error {
	return clipboard.WriteAll(s)
}

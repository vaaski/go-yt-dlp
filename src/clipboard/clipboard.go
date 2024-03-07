//go:build !cgo

package clipboard

func GetClipboardUrl() string {
	return ""
}

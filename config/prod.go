//go:build prod
// +build prod

package config

func init() {
	isDevMode = false
}

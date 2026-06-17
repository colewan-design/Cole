//go:build !windows

package main

func encryptSecret(value string) (string, error) {
	return value, nil
}

func decryptSecret(value string) (string, error) {
	return value, nil
}

//go:build windows

package main

import (
	"encoding/base64"
	"unsafe"

	"golang.org/x/sys/windows"
)

func encryptSecret(value string) (string, error) {
	input := makeBlob([]byte(value))
	var output windows.DataBlob

	if err := windows.CryptProtectData(&input, nil, nil, 0, nil, 0, &output); err != nil {
		return "", err
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(output.Data)))

	encrypted := unsafe.Slice(output.Data, output.Size)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func decryptSecret(value string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	input := makeBlob(raw)
	var output windows.DataBlob
	if err := windows.CryptUnprotectData(&input, nil, nil, 0, nil, 0, &output); err != nil {
		return "", err
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(output.Data)))

	decrypted := unsafe.Slice(output.Data, output.Size)
	return string(decrypted), nil
}

func makeBlob(data []byte) windows.DataBlob {
	if len(data) == 0 {
		return windows.DataBlob{}
	}
	return windows.DataBlob{
		Size: uint32(len(data)),
		Data: &data[0],
	}
}

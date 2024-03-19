package gocommand

import (
	"bytes"
	"errors"
	"strings"
)

type errorBuffer struct{ bytes.Buffer }

func (buffer *errorBuffer) error() error {
	text := strings.TrimRight(buffer.String(), "\n")

	if len(text) > 0 {
		return errors.New(text)
	}

	return nil
}

type versionsErrorBuffer struct{ errorBuffer }

var errNotFound = errors.New("404 not found")

func (buffer *versionsErrorBuffer) error() error {
	err := buffer.errorBuffer.error()
	if err == nil {
		return nil
	}

	if strings.HasSuffix(err.Error(), ": 404 Not Found") {
		return errNotFound
	}

	return err
}

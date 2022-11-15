package repos

import (
	"io"
)

type ImageRepo interface {
	UploadImage(name string, r io.Reader) (string, error)
}

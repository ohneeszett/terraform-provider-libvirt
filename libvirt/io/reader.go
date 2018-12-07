package io

type Sized interface {
	Size() (int64, error)
}

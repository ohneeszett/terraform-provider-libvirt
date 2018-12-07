package ioutil

type Sized interface {
	Size() (int64, error)
}

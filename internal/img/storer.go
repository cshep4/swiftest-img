package img

type Storer interface {
	Upload(img []byte) (string, error)
	Get(uri string) ([]byte, error)
}

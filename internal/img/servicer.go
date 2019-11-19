package img

type Servicer interface {
	Upload(req UploadRequest) (*UploadResult, error)
	GetDocument(uri string) (*Document, error)
}

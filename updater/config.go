package updater

type (
	Options struct {
		Bucket string
		Path   string
	}
	MetaData struct {
		ContentType string
		Hash        string
	}
)

package core

type BrowserProvider interface {
	GetRenderedHTML(url string) (string, error)
}

type ResultPublisher interface {
	Publish(result interface{})
}

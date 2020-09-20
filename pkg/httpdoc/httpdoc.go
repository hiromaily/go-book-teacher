package httpdoc

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// GetHTMLDocs is to get http resource as document data from set URL
func GetHTMLDocs(url string) (*goquery.Document, error) {
	// goquery.NewDocument is Deprecated

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(res.Body)
}

package debug

import (
	"github.com/bookerzzz/grok"
)

// PrintDetail is wrapper of grok.Value
func PrintDetail(value interface{}, options ...grok.Option) {
	grok.Value(value, options...)
}

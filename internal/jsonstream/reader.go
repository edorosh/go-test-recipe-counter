package jsonstream

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/edorosh/go-test-recipe-counter/internal/app"
)

type readFunc func(r app.Recipe)
type errReadFunc func(err error)

// Reader decodes recipes in JSON format in streams.
// It requires small memory footprint for operation.
type Reader struct {
	reader     io.Reader
	handler    readFunc
	errHandler errReadFunc
}

// ErrJSONParse represents error in parsing a JSON file.
var ErrJSONParse = errors.New("JSON parse error")

// NewReader creates new instance.
func NewReader(r io.Reader) *Reader {
	return &Reader{reader: r}
}

// Handler set recipe handler.
func (d *Reader) Handler(f readFunc) *Reader {
	d.handler = f
	return d
}

// ErrHandler set error handler on read.
func (d *Reader) ErrHandler(f errReadFunc) *Reader {
	d.errHandler = f
	return d
}

// Read publishes Recipes in the channel.
func (d *Reader) Read(ctx context.Context) {
	dec := json.NewDecoder(d.reader)

	// read open bracket
	if _, err := dec.Token(); err != nil && err != io.EOF {
		d.errHandler(wrapDecodeJSONErr(err))
		return
	}

	// read recipes
	for dec.More() {
		var r app.Recipe
		if err := dec.Decode(&r); err != nil {
			d.errHandler(wrapDecodeJSONErr(err))
			return
		}

		select {
		case <-ctx.Done():
			d.errHandler(ctx.Err())
			return
		default:
		}

		d.handler(r)
	}

	// read closing bracket
	if _, err := dec.Token(); err != nil && err != io.EOF {
		d.errHandler(wrapDecodeJSONErr(err))
		return
	}
}

func wrapDecodeJSONErr(err error) error {
	return fmt.Errorf("[jsonstream] %w: %v", ErrJSONParse, err)
}

package notice

import (
	"io"
)

func ToReadCloser(r io.Reader) io.ReadCloser {
	if rc, ok := r.(io.ReadCloser); ok {
		return rc
	}
	return io.NopCloser(r)
}

func ToWriteCloser(w io.Writer) io.WriteCloser {
	if wc, ok := w.(io.WriteCloser); ok {
		return wc
	}
	return nopWriteCloser{Writer: w}
}

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }

func newProxyReader(r io.Reader, prog Progressor) io.ReadCloser {
	if r == nil {
		prog.AutoComplete(true)
		return io.NopCloser(&proxyNonReader{prog: prog})
	}
	rc := ToReadCloser(r)
	pr := &proxyReader{ReadCloser: rc, prog: prog}
	if _, ok := r.(io.WriterTo); ok {
		return proxyWriterTo{proxyReader: pr}
	}
	return pr
}

type proxyReader struct {
	io.ReadCloser
	prog Progressor
}

func (x proxyReader) Read(p []byte) (int, error) {
	n, err := x.ReadCloser.Read(p)
	x.prog.Done(int64(n))
	if err == io.EOF {
		x.prog.Complete()
	}
	return n, err
}

type proxyNonReader struct {
	prog Progressor
}

func (x *proxyNonReader) Read(p []byte) (int, error) {
	n := len(p)
	x.prog.Done(int64(n))
	return n, nil
}

func newProxyWriter(w io.Writer, prog Progressor) io.WriteCloser {
	if w == nil {
		prog.AutoComplete(true)
		return nopWriteCloser{Writer: &proxyNonWriter{prog: prog}}
	}
	wc := ToWriteCloser(w)
	pw := &proxyWriter{WriteCloser: wc, prog: prog}
	if _, ok := w.(io.ReaderFrom); ok {
		return proxyReaderFrom{proxyWriter: pw}
	}
	return pw
}

type proxyWriter struct {
	io.WriteCloser
	prog Progressor
}

func (x *proxyWriter) Write(p []byte) (int, error) {
	n, err := x.WriteCloser.Write(p)
	x.prog.Done(int64(n))
	if err == io.EOF {
		x.prog.Complete()
	}
	return n, err
}

type proxyNonWriter struct {
	prog Progressor
}

func (x *proxyNonWriter) Write(p []byte) (int, error) {
	n := len(p)
	x.prog.Done(int64(n))
	return n, nil
}

type proxyWriterTo struct {
	*proxyReader
}

func (x *proxyWriterTo) WriteTo(w io.Writer) (int64, error) {
	n, err := x.ReadCloser.(io.WriterTo).WriteTo(w)
	x.prog.Done(int64(n))
	if err == io.EOF {
		x.prog.Complete()
	}
	return n, err
}

type proxyReaderFrom struct {
	*proxyWriter
}

func (x *proxyReaderFrom) ReadFrom(r io.Reader) (int64, error) {
	n, err := x.WriteCloser.(io.ReaderFrom).ReadFrom(r)
	x.prog.Done(int64(n))
	if err == io.EOF {
		x.prog.Complete()
	}
	return n, err
}

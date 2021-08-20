package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/url"

	"github.com/tkuchiki/alp/errors"
)

type Parser interface {
	Parse() (*ParsedLog, error)
	ReadBytes() int
	SetReadBytes(n int)
	Seek(n int) error
}

type ParsedLog struct {
	Uri    string
	Method string
	User   string
}

func (pl *ParsedLog) Request() string {
	return fmt.Sprintf("%s %s", pl.Method, pl.Uri)
}

type LogKeys struct {
	uri    string
	method string
	user   string
}

type logKey func(*LogKeys)

func uriKey(s string) logKey {
	return func(lk *LogKeys) {
		if s != "" {
			lk.uri = s
		}
	}
}

func methodKey(s string) logKey {
	return func(lk *LogKeys) {
		if s != "" {
			lk.method = s
		}
	}
}

func userKey(s string) logKey {
	return func(lk *LogKeys) {
		if s != "" {
			lk.user = s
		}
	}
}

func newLogKeys(lk ...logKey) *LogKeys {
	lks := &LogKeys{
		uri:    "uri",
		method: "method",
		user:   "user",
	}

	for _, l := range lk {
		l(lks)
	}

	return lks
}

func readline(reader *bufio.Reader) ([]byte, int, error) {
	var b []byte
	var i int
	var err error
	for {
		line, _err := reader.ReadBytes('\n')
		if _err == io.EOF && len(line) == 0 {
			err = io.EOF
			break
		}

		if _err != io.EOF && _err != nil {
			return []byte{}, 0, err
		}
		trimedLine := bytes.TrimRight(line, "\r\n")
		if len(trimedLine) > 0 {
			b = append(b, trimedLine...)
		} else {
			err = errors.SkipReadLineErr
		}

		size := len(line)
		i += size

		if line[size-1] == byte('\n') {
			break
		}
	}

	return b, i, err
}

func NewParsedLog(uri, method, user string) *ParsedLog {
	return &ParsedLog{
		Uri:    uri,
		Method: method,
		User:   user,
	}
}

func toStats(parsedValue map[string]string, keys *LogKeys, strictMode, queryString, qsIgnoreValues bool) (*ParsedLog, error) {
	u, err := url.Parse(parsedValue[keys.uri])
	if err != nil {
		return nil, errSkipReadLine(strictMode, err)
	}

	uri := normalizeURL(u, queryString, qsIgnoreValues)
	if uri == "" {
		return nil, errSkipReadLine(strictMode, err)
	}

	method := parsedValue[keys.method]
	user := parsedValue[keys.user]

	return NewParsedLog(uri, method, user), nil
}

func normalizeURL(src *url.URL, queryString, qsIgnoreValues bool) string {
	if src.RawQuery == "" {
		return src.String()
	}

	u := *src // basic clone
	if !queryString {
		u.RawQuery = ""
		return u.String()
	}

	if qsIgnoreValues {
		values := u.Query()
		for q := range values {
			values.Set(q, "xxx")
		}
		u.RawQuery = values.Encode()
	} else {
		u.RawQuery = u.Query().Encode() // re-encode to sort queries
	}
	return u.String()
}

func errSkipReadLine(strictMode bool, err error) error {
	if strictMode {
		return err
	}

	return errors.SkipReadLineErr
}

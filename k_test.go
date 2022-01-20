package k

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
)

func CountLines(r io.Reader) (int, error) {
	sc := bufio.NewScanner(r)
	lines := 0

	for sc.Scan() {
		lines++
	}

	return lines, sc.Err()
}

////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////

type Header struct {
	Key, Value string
}

type Status struct {
	Code   int
	Reason string
}

func WriteResponse(w io.Writer, st Status, headers []Header, body io.Reader) error {
	_, err := fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", st.Code, st.Reason)
	if err != nil {
		return err
	}

	for _, h := range headers {
		_, err := fmt.Fprintf(w, "%s: %s\r\n", h.Key, h.Value)
		if err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(w, "\r\n"); err != nil {
		return err
	}

	_, err = io.Copy(w, body)
	return err
}

type errWriter struct {
	io.Writer
	err error
}

func (e *errWriter) Write(buf []byte) (int, error) {
	if e.err != nil {
		return 0, e.err
	}

	var n int
	n, e.err = e.Writer.Write(buf)
	return n, nil
}

func WriteResponse2(w io.Writer, st Status, headers []Header, body io.Reader) error {
	ew := &errWriter{Writer: w}
	fmt.Fprintf(ew, "HTTP/1.1 %d %s\r\n", st.Code, st.Reason)

	for _, h := range headers {
		fmt.Fprintf(ew, "%s: %s\r\n", h.Key, h.Value)
	}

	fmt.Fprint(ew, "\r\n")
	io.Copy(ew, body)

	return ew.err
}

////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////

// 處理一個錯誤時，帶了兩個任務
func WriterAll(w io.Writer, buf []byte) error {
	_, err := w.Write(buf)
	if err != nil {
		log.Println("unable to write:", err)
		return err
	}

	return nil
}

type MyErrors struct {
	err error
	msg string
}

func (e *MyErrors) Error() string { return e.err.Error() }

////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////

func TestMain(t *testing.T) {
	err := t5()
	if err != nil {
		switch errors.Cause(err).(type) {
		case *MyErrors:
			fmt.Printf("MyErrors::: %+v\n", err.err)
		default:
			fmt.Printf("Fatal: %T\n", err)
		}
	}
}

////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////

func ReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "open failed")
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Wrap(err, "read failed")
	}
	return buf, nil
}

func ReadConfig() ([]byte, error) {
	home := os.Getenv("HOME")
	config, err := ReadFile(filepath.Join(home, ".settings.xml"))
	return config, errors.WithMessage(err, "could not read config")
}

func TestRead(t *testing.T) {
	_, err := ReadConfig()
	if err != nil {
		fmt.Printf("Error's string: %v\n\n\n", err)
		fmt.Printf("Type: %T\nString: %v\n\n", errors.Cause(err), errors.Cause(err))
		fmt.Printf("Just error: %v\n\n", err)
		fmt.Printf("stack: %+v\n", err)
		os.Exit(1)
	}
}

////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////

func t1() *MyErrors { return &MyErrors{err: errors.New("gg"), msg: "msg"} }
func t2() error     { return errors.WithMessage(t1(), "with message") }
func t3() error     { return errors.Wrap(t2(), "Wrap") }
func t4() error     { return errors.Errorf("Errorf :%v", t3()) }
func t5() *MyErrors { return &MyErrors{err: t4(), msg: "haha"} }

func TestThirdParty(t *testing.T) {
	t.Run("NormalS5():", func(t *testing.T) {
		fmt.Printf("Type: %T\nString: %v\n", errors.Cause(s5()), s5().Error())
		fmt.Printf("stack: %+v\n", s5())
	})

	t.Run("s5() use errorf again:", func(t *testing.T) {
		fmt.Printf("Type: %T\nString: %v\n", errors.Cause(s5UseErrorfAgain()), s5UseErrorfAgain().Error())
		fmt.Printf("stack: %+v\n", s5UseErrorfAgain())
	})
}

func s1() error               { return errors.Errorf("Without Wrap") }
func s2() error               { return errors.WithMessage(s1(), "Where is MSG?") }
func s3() error               { return errors.WithMessagef(s2(), "err: %v", "gg") }
func s4() error               { return errors.WithStack(s3()) }
func s5UseErrorfAgain() error { return errors.Errorf("%v", s4()) }
func s5() error               { return s4() }

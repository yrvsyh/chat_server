package e

import (
	"errors"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

type withMessage struct {
	cause error
	msg   string
}

func New(message string) error {
	return errors.New(message)
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   message,
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	return Wrap(err, message)
}

func (w *withMessage) Error() string { return w.msg + ": " + w.cause.Error() }

func (w *withMessage) Cause() error { return w.cause }

func (w *withMessage) Unwrap() error { return w.cause }

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			io.WriteString(s, w.msg)
			return
		}
		io.WriteString(s, w.Error())
	case 's', 'q':
		io.WriteString(s, w.msg)
	}
}

func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

func Throw(msg string) {
	panic(New(msg))
}

func Check(err error, msg string) {
	if err != nil {
		panic(Wrap(err, msg))
	}
}

func Must0(err error) {
	if err != nil {
		panic(err)
	}
}

func Must1[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}

	return val
}

func Must2[T1 any, T2 any](val1 T1, val2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}

	return val1, val2
}

func Must3[T1 any, T2 any, T3 any](val1 T1, val2 T2, val3 T3, err error) (T1, T2, T3) {
	if err != nil {
		panic(err)
	}

	return val1, val2, val3
}

func Must4[T1 any, T2 any, T3 any, T4 any](val1 T1, val2 T2, val3 T3, val4 T4, err error) (T1, T2, T3, T4) {
	if err != nil {
		panic(err)
	}

	return val1, val2, val3, val4
}

func Must5[T1 any, T2 any, T3 any, T4 any, T5 any](val1 T1, val2 T2, val3 T3, val4 T4, val5 T5, err error) (T1, T2, T3, T4, T5) {
	if err != nil {
		panic(err)
	}

	return val1, val2, val3, val4, val5
}

func Try(try func()) (e error) {
	e = nil
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if ok {
				e = err
			} else {
				panic(r)
			}
		}
	}()
	try()
	return
}

func TryCatch(try func(), catch func(error)) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if ok {
				catch(err)
			} else {
				logrus.Errorf("TryCatch error [%v]", r)
			}
		}
	}()

	try()
}

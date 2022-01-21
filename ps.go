// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pdf

import (
	"errors"
	"fmt"
	"io"
)

// A Stack represents a stack of values.
type Stack struct {
	stack []Value
}

// Len returns the amount of values in this stack.
func (stk *Stack) Len() int {
	return len(stk.stack)
}

// Push puts the given value onto the top of the stack.
func (stk *Stack) Push(v Value) {
	stk.stack = append(stk.stack, v)
}

// Pop removes and returns the value at the top of the stack.
func (stk *Stack) Pop() Value {
	n := len(stk.stack)
	if n == 0 {
		return Value{}
	}
	v := stk.stack[n-1]
	stk.stack[n-1] = Value{}
	stk.stack = stk.stack[:n-1]
	return v
}

func newDict() Value {
	return Value{nil, objptr{}, make(dict)}
}

// Interpret interprets the content in a stream as a basic PostScript program,
// pushing values onto a stack and then calling the do function to execute
// operators. The do function may push or pop values from the stack as needed
// to implement op.
//
// Interpret handles the operators "dict", "currentdict", "begin", "end", "def", and "pop" itself.
//
// Interpret is not a full-blown PostScript interpreter. Its job is to handle the
// very limited PostScript found in certain supporting file formats embedded
// in PDF files, such as cmap files that describe the mapping from font code
// points to Unicode code points.
//
// There is no support for executable blocks, among other limitations.
//
func Interpret(strm Value, do func(stk *Stack, op string) (bool, error)) error {
	var rd io.Reader
	var err error
	switch strm.data.(type) {
	case stream:
		rd, err = strm.Reader()
		if err != nil {
			return err
		}
	case array:
		values, err := strm.Values()
		if err != nil {
			return err
		}
		rd, err = MultiReader(values)
		if err != nil {
			return err
		}
	default:
		return errors.New("interpret: value is neither stream nor array")
	}

	b := newBuffer(rd, 0)
	b.allowEOF = true
	b.allowObjptr = false
	b.allowStream = false
	var stk Stack
	var dicts []dict
Reading:
	for {
		tok, err := b.readToken()
		if err != nil {
			return err
		}
		if tok == io.EOF {
			break
		}
		if kw, ok := tok.(keyword); ok {
			switch kw {
			case "null", "[", "]", "<<", ">>":
				break
			default:
				for i := len(dicts) - 1; i >= 0; i-- {
					if v, ok := dicts[i][name(kw)]; ok {
						stk.Push(Value{nil, objptr{}, v})
						continue Reading
					}
				}
				cont, err := do(&stk, string(kw))
				if err != nil {
					return err
				}
				if !cont {
					return nil
				}
				continue
			case "dict":
				stk.Pop()
				stk.Push(Value{nil, objptr{}, make(dict)})
				continue
			case "currentdict":
				if len(dicts) == 0 {
					return errors.New("no current dictionary")
				}
				stk.Push(Value{nil, objptr{}, dicts[len(dicts)-1]})
				continue
			case "begin":
				d := stk.Pop()
				if d.Kind() != Dict {
					return errors.New("cannot begin non-dict")
				}
				dicts = append(dicts, d.data.(dict))
				continue
			case "end":
				if len(dicts) <= 0 {
					return errors.New("mismatched begin/end")
				}
				dicts = dicts[:len(dicts)-1]
				continue
			case "def":
				if len(dicts) <= 0 {
					return errors.New("def without open dict")
				}
				val := stk.Pop()
				key, ok := stk.Pop().data.(name)
				if !ok {
					// return errors.New("def of non-name")
					// TODO: Examine what’s wrong with the `testdata/buggy-pdf-v1.6.pdf` and implement a proper fix.
					fmt.Println("warn: def of non-name")
					continue
				}
				dicts[len(dicts)-1][key] = val.data
				continue
			case "pop":
				stk.Pop()
				continue
			}
		}
		b.unreadToken(tok)
		obj, err := b.readObject()
		if err != nil {
			return err
		}
		stk.Push(Value{nil, objptr{}, obj})
	}
	return nil
}

type seqReader struct {
	rd     io.Reader
	offset int64
}

func (r *seqReader) ReadAt(buf []byte, offset int64) (int, error) {
	if offset != r.offset {
		return 0, fmt.Errorf("non-sequential read of stream")
	}
	n, err := io.ReadFull(r.rd, buf)
	r.offset += int64(n)
	return n, err
}

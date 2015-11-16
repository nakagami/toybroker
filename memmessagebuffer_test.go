/*******************************************************************************
The MIT License (MIT)

Copyright (c) 2015 Hajime Nakagami

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*******************************************************************************/

package toybroker

import (
	"errors"
	"reflect"
	"testing"
)

func TestMemMessageBuffer(t *testing.T) {
	var err error
	mb1 := NewMemoryMessageBuffer("foo")
	mb2 := NewMemoryMessageBuffer("bar")
	mb1.Set(1, []byte{1, 2, 3})
	mb1.Set(2, []byte{4, 5, 6})
	mb2.Set(1, []byte{7, 8, 9})

	if len(mb1.List()) != 2 {
		err = errors.New("mb1.List()")
	}
	if !reflect.DeepEqual(mb1.Get(1), []byte{1, 2, 3}) {
		err = errors.New("mb1.Get(1)")
	}
	if !reflect.DeepEqual(mb1.Get(2), []byte{4, 5, 6}) {
		err = errors.New("mb1.Get(2)")
	}

	if len(mb2.List()) != 1 {
		err = errors.New("mb2.List()")
	}

	mb1.Delete(1)
	if len(mb1.List()) != 1 {
		err = errors.New("mb1.Delete()")
	}
	if mb1.Get(1) != nil {
		err = errors.New("mb1.Get(1) is not nil")
	}

	if err != nil {
		t.Error(err.Error())
	}
}

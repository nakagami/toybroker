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

package main

import (
	"errors"
	"testing"
)

func TestEncode(t *testing.T) {
	var err error
	if len(encodeRemainLength(0)) != 1 {
		err = errors.New("encodeRemainLenght(0)")
	}
	if len(encodeRemainLength(127)) != 1 {
		err = errors.New("encodeRemainLenght(127)")
	}
	if len(encodeRemainLength(128)) != 2 {
		err = errors.New("encodeRemainLenght(128)")
	}
	if len(encodeRemainLength(16383)) != 2 {
		err = errors.New("encodeRemainLenght(16383)")
	}
	if len(encodeRemainLength(16384)) != 3 {
		err = errors.New("encodeRemainLenght(16384)")
	}
	if len(encodeRemainLength(2097151)) != 3 {
		err = errors.New("encodeRemainLenght(2097151)")
	}
	if len(encodeRemainLength(2097152)) != 4 {
		err = errors.New("encodeRemainLenght(2097152)")
	}
	if len(encodeRemainLength(268435455)) != 4 {
		err = errors.New("encodeRemainLenght(268435455)")
	}

	if err != nil {
		t.Error(err.Error())
	}

}

func TestDecode(t *testing.T) {
	var err error
	if decodeRemainLength(encodeRemainLength(0)) != 0 {
		err = errors.New("decodeRemainLength(0)")
	}
	if decodeRemainLength(encodeRemainLength(127)) != 127 {
		err = errors.New("decodeRemainLength(127)")
	}
	if decodeRemainLength(encodeRemainLength(128)) != 128 {
		err = errors.New("decodeRemainLength(128)")
	}
	if decodeRemainLength(encodeRemainLength(16383)) != 16383 {
		err = errors.New("decodeRemainLength(16383)")
	}
	if decodeRemainLength(encodeRemainLength(16384)) != 16384 {
		err = errors.New("decodeRemainLength(16384)")
	}
	if decodeRemainLength(encodeRemainLength(2097151)) != 2097151 {
		err = errors.New("decodeRemainLength(2097151)")
	}
	if decodeRemainLength(encodeRemainLength(2097152)) != 2097152 {
		err = errors.New("decodeRemainLength(2097152)")
	}
	if decodeRemainLength(encodeRemainLength(268435455)) != 268435455 {
		err = errors.New("decodeRemainLength(268435455)")
	}

	if err != nil {
		t.Error(err.Error())
	}

}

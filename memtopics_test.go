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
	"sort"
	"testing"
)

func TestTopics(t *testing.T) {
	var err error
	topics := NewMemoryTopics()
	topics.Add("foo/bar", "client1", 0)
	topics.Add("foo/bar", "client2", 0)
	topics.Add("foo/baz", "client1", 0)

	if !reflect.DeepEqual(topics.TopicList(), []string{"foo/bar", "foo/baz"}) {
		err = errors.New("topic.TopicList()")
	}

	r := topics.List("foo/bar")
	sort.Strings(r)
	if !reflect.DeepEqual(r, []string{"client1", "client2"}) {
		err = errors.New("topics.List(\"foo/bar\")")
	}

	if !reflect.DeepEqual(topics.List("foo/baz"), []string{"client1"}) {
		err = errors.New("topics.List(\"foo/baz\")")
	}

	topics.Remove("foo/bar", "client1")
	if !reflect.DeepEqual(topics.List("foo/bar"), []string{"client2"}) {
		err = errors.New("Remove(\"boo/bar\", \"client1\")")
	}

	if err != nil {
		t.Error(err.Error())
	}
}

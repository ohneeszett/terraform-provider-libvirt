// Derived from https://github.com/paulcager/aio
// Copyright 2017 Paul Cager
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package io

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"testing"

	"fmt"
	"io/ioutil"

	"encoding/base64"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadEmpty(t *testing.T) {
	r, err := NewAnyReader(strings.NewReader(""))
	assert.NoError(t, err)
	b := make([]byte, 12)
	n, err := r.Read(b)
	assert.EqualValues(t, 0, n)
	assert.True(t, err == nil || err == io.EOF, "Err was %s", err)
}

func TestReadPlain(t *testing.T) {
	const str = "HelloWorld"
	r, err := NewAnyReader(strings.NewReader(str))
	assert.NoError(t, err)

	b := make([]byte, len(str)+12)
	n, err := r.Read(b)
	assert.NoError(t, err)
	assert.EqualValues(t, len(str), n)
	assert.EqualValues(t, str, string(b[:n]))

	n, err = r.Read(b)
	assert.EqualValues(t, 0, n)
	assert.EqualValues(t, io.EOF, err)
}

func TestReadPlainShortReads(t *testing.T) {
	const str = "HelloWorld"
	r, err := NewAnyReader(strings.NewReader(str))
	assert.NoError(t, err)

	b := make([]byte, 1)
	for i := range str {
		n, err := r.Read(b)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, n)
		assert.EqualValues(t, str[i], b[0])
	}
	n, err := r.Read(b)
	assert.EqualValues(t, 0, n)
	assert.EqualValues(t, io.EOF, err)
}

func TestReadEmptyGZIP(t *testing.T) {
	buff := new(bytes.Buffer)
	gz := gzip.NewWriter(buff)
	gz.Close()

	r, err := NewAnyReader(buff)
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(r)
	require.NoError(t, err)
	assert.Empty(t, b)
}

func TestReadGZIP(t *testing.T) {
	buff := new(bytes.Buffer)
	gz := gzip.NewWriter(buff)
	fmt.Fprint(gz, "Hello World")
	gz.Close()

	r, err := NewAnyReader(buff)
	assert.NoError(t, err)
	b, err := ioutil.ReadAll(r)
	require.NoError(t, err)
	assert.EqualValues(t, "Hello World", string(b))
}

func TestReadXZ(t *testing.T) {
	compressed := `/Td6WFoAAATm1rRGAgAhARYAAAB0L+WjAQAKSGVsbG8gV29ybGQAAMbNtcdndHQ+AAEjC8Ib/QkftvN9AQAAAAAEWVo=`
	r := base64.NewDecoder(base64.StdEncoding, strings.NewReader(compressed))

	r, err := NewAnyReader(r)
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(r)
	require.NoError(t, err)
	assert.EqualValues(t, "Hello World", string(b))
}

func TestReadBZ2(t *testing.T) {
	compressed := `QlpoOTFBWSZTWQZcidoAAACXgEAAAEAAgAYEkAAgADEMCCAxqRbEHUHi7kinChIAy5E7QA==`
	d := base64.NewDecoder(base64.StdEncoding, strings.NewReader(compressed))

	r, err := NewAnyReader(d)
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(r)
	require.NoError(t, err)
	assert.EqualValues(t, "Hello World", string(b))
}

// Copyright 2018 the Service Broker Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by assets/pack.go. DO NOT EDIT.

package tflinuxamd64

import (
	"archive/zip"
	"bytes"
	"io"
)

func NewZipReader() (*zip.Reader, error) {
	fd := File{
		chunk0,
		chunk1,
		chunk2,
		chunk3,
		chunk4,
		chunk5,
		chunk6,
		chunk7,
		chunk8,
		chunk9,
		chunk10,
		chunk11,
		chunk12,
		chunk13,
		chunk14,
		chunk15,
		chunk16,
		chunk17,
		chunk18,
		chunk19,
		chunk20,
		chunk21,
		chunk22,
		chunk23,
		chunk24,
		chunk25,
		chunk26,
		chunk27,
		chunk28,
		chunk29,
		chunk30,
		chunk31,
		chunk32,
		chunk33,
		chunk34,
		chunk35,
		chunk36,
		chunk37,
		chunk38,
		chunk39,
		chunk40,
		chunk41,
		chunk42,
		chunk43,
		chunk44,
		chunk45,
		chunk46,
		chunk47,
		chunk48,
		chunk49,
		chunk50,
		chunk51,
		chunk52,
		chunk53,
		chunk54,
		chunk55,
		chunk56,
		chunk57,
		chunk58,
		chunk59,
		chunk60,
		chunk61,
		chunk62,
		chunk63,
		chunk64,
		chunk65,
		chunk66,
		chunk67,
		chunk68,
		chunk69,
		chunk70,
		chunk71,
		chunk72,
		chunk73,
		chunk74,
		chunk75,
		chunk76,
		chunk77,
		chunk78,
		chunk79,
		chunk80,
		chunk81,
		chunk82,
		chunk83,
		chunk84,
		chunk85,
		chunk86,
		chunk87,
		chunk88,
		chunk89,
		chunk90,
		chunk91,
		chunk92,
		chunk93,
		chunk94,
		chunk95,
		chunk96,
		chunk97,
		chunk98,
		chunk99,
		chunk100,
		chunk101,
		chunk102,
		chunk103,
		chunk104,
		chunk105,
		chunk106,
		chunk107,
		chunk108,
		chunk109,
		chunk110,
		chunk111,
		chunk112,
		chunk113,
		chunk114,
		chunk115,
		chunk116,
		chunk117,
		chunk118,
		chunk119,
		chunk120,
		chunk121,
		chunk122,
		chunk123,
		chunk124,
		chunk125,
		chunk126,
		chunk127,
		chunk128,
		chunk129,
		chunk130,
		chunk131,
		chunk132,
		chunk133,
		chunk134,
		chunk135,
		chunk136,
		chunk137,
		chunk138,
		chunk139,
		chunk140,
		chunk141,
		chunk142,
		chunk143,
		chunk144,
		chunk145,
		chunk146,
		chunk147,
		chunk148,
		chunk149,
		chunk150,
		chunk151,
		chunk152,
		chunk153,
		chunk154,
		chunk155,
		chunk156,
		chunk157,
		chunk158,
		chunk159,
		chunk160,
		chunk161,
		chunk162,
		chunk163,
		chunk164,
		chunk165,
		chunk166,
		chunk167,
		chunk168,
		chunk169,
		chunk170,
		chunk171,
		chunk172,
		chunk173,
		chunk174,
		chunk175,
		chunk176,
		chunk177,
		chunk178,
		chunk179,
		chunk180,
		chunk181,
		chunk182,
		chunk183,
		chunk184,
		chunk185,
		chunk186,
		chunk187,
		chunk188,
		chunk189,
		chunk190,
		chunk191,
		chunk192,
		chunk193,
		chunk194,
		chunk195,
		chunk196,
		chunk197,
		chunk198,
		chunk199,
		chunk200,
		chunk201,
		chunk202,
		chunk203,
		chunk204,
		chunk205,
		chunk206,
		chunk207,
		chunk208,
		chunk209,
		chunk210,
		chunk211,
		chunk212,
		chunk213,
		chunk214,
		chunk215,
		chunk216,
		chunk217,
		chunk218,
		chunk219,
		chunk220,
		chunk221,
		chunk222,
		chunk223,
		chunk224,
		chunk225,
		chunk226,
		chunk227,
		chunk228,
		chunk229,
		chunk230,
		chunk231,
		chunk232,
		chunk233,
		chunk234,
		chunk235,
		chunk236,
		chunk237,
		chunk238,
		chunk239,
		chunk240,
		chunk241,
		chunk242,
	}

	return zip.NewReader(fd, fd.Length())
}

type chunk []byte
type File []chunk

func (f File) ReadAt(p []byte, off int64) (n int, err error) {
	iOff := int(off)
	slice := f.slice(iOff, iOff+len(p))

	for i, b := range slice {
		p[i] = b
	}

	n = len(slice)
	if n < len(p) {
		err = io.EOF
	}

	return
}

func (f File) Length() int64 {
	l := 0
	for _, chunk := range f {
		l += len(chunk)
	}
	return int64(l)
}

func (f File) slice(start, end int) []byte {
	buf := new(bytes.Buffer)

	chunkStart := 0
	for _, chunk := range f {
		chunkEnd := chunkStart + len(chunk)

		overlapStart := imin(chunkEnd, imax(chunkStart, start))
		overlapEnd := imax(chunkStart, imin(chunkEnd, end))

		if overlapStart != overlapEnd {
			buf.Write(chunk[overlapStart-chunkStart : overlapEnd-chunkStart])
		}
		chunkStart = chunkEnd
	}

	return buf.Bytes()
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

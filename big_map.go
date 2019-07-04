package main

import (
	"errors"
	"math/big"
	"strconv"
)

type bigMapElement struct {
	v interface{}
	m map[string]*bigMapElement
}

type BigMap struct {
	partLen int
	m       map[string]*bigMapElement
	length  *big.Int
}

func NewBigMap(partLen int) *BigMap {
	if partLen <= 0 {
		partLen = strconv.IntSize >> 1
	}
	return &BigMap{partLen: partLen, length: big.NewInt(0)}
}

func (bm *BigMap) IsEmpty() bool {
	return bm == nil || bm.m == nil
}

func (bm *BigMap) Len() *big.Int {
	if bm == nil {
		return big.NewInt(0)
	}
	bm.ckLength()
	return new(big.Int).Set(bm.length)
}

func (bm *BigMap) PartLen() int {
	if bm == nil {
		return 0
	}
	bm.ckPartLen()
	return bm.partLen
}

func (bm *BigMap) Load(key string) (value interface{}, ok bool) {
	if bm == nil {
		return
	}
	bm.ckPartLen()
	m := bm.m
	var e *bigMapElement
	n := len(key)
	if n == 0 {
		n = 1
	}
	var end int
	for i := 0; i < n; i += bm.partLen {
		if m == nil {
			return
		}
		end = i + bm.partLen
		if end > len(key) {
			end = len(key)
		}
		e = m[key[i:end]]
		if e == nil {
			return
		}
		m = e.m
	}
	return e.v, true
}

func (bm *BigMap) Store(key string, value interface{}) {
	if bm == nil {
		panic(errors.New("BigMap is nil"))
	}
	bm.ckPartLen()
	if bm.m == nil {
		bm.m = make(map[string]*bigMapElement)
	}
	m := bm.m
	var e *bigMapElement
	isNew := false
	n := len(key)
	if n == 0 {
		n = 1
	}
	var end int
	for i := 0; i < n; i += bm.partLen {
		if m == nil {
			m = make(map[string]*bigMapElement)
			e.m = m
		}
		end = i + bm.partLen
		if end > len(key) {
			end = len(key)
		}
		e = m[key[i:end]]
		if e == nil {
			e = new(bigMapElement)
			m[key[i:end]] = e
			isNew = true
		}
		m = e.m
	}
	e.v = value
	if isNew {
		bm.ckLength()
		bm.length.Add(bm.length, One)
	}
}

func (bm *BigMap) ckPartLen() {
	if bm == nil {
		return
	}
	if bm.partLen == 0 {
		bm.partLen = strconv.IntSize >> 1
	}
}

func (bm *BigMap) ckLength() {
	if bm == nil {
		return
	}
	if bm.length == nil {
		bm.length = big.NewInt(0)
	}
}

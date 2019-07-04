package main

import (
	"errors"
	"math/big"
)

type bigStackElement struct {
	v    interface{}
	prev *bigStackElement
}

type BigStack struct {
	top    *bigStackElement
	length *big.Int
}

func NewBigStack() *BigStack {
	return &BigStack{length: big.NewInt(0)}
}

func (bs *BigStack) IsEmpty() bool {
	return bs == nil || bs.top == nil
}

func (bs *BigStack) Len() *big.Int {
	if bs == nil {
		return big.NewInt(0)
	}
	bs.ckLength()
	return new(big.Int).Set(bs.length)
}

func (bs *BigStack) Top() (value interface{}, ok bool) {
	if bs == nil || bs.top == nil {
		return
	}
	return bs.top.v, true
}

func (bs *BigStack) Push(value interface{}) {
	if bs == nil {
		panic(errors.New("BigStack is nil"))
	}
	bs.ckLength()
	e := &bigStackElement{v: value, prev: bs.top}
	bs.top = e
	bs.length.Add(bs.length, One)
}

func (bs *BigStack) Pop() (value interface{}, ok bool) {
	if bs == nil || bs.top == nil {
		return
	}
	bs.ckLength()
	e := bs.top
	bs.top = e.prev
	bs.length.Sub(bs.length, One)
	return e.v, true
}

func (bs *BigStack) Clear() {
	if bs.IsEmpty() {
		return
	}
	bs.top = nil
	if bs.length != nil {
		bs.length.SetInt64(0)
	}
}

func (bs *BigStack) ckLength() {
	if bs == nil {
		return
	}
	if bs.length == nil {
		bs.length = big.NewInt(0)
	}
}

/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2018 Yu Jing <yu@argcv.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package helpers

import (
	"github.com/argcv/webeh/log"
	"runtime"
	"sync/atomic"
)

type WaitGroupWithState struct {
	st int64
}

func NewWaitGroupWithState() *WaitGroupWithState {
	return &WaitGroupWithState{
		st: 0,
	}
}

func (wg *WaitGroupWithState) Add(delta int64) int64 {
	newSt := atomic.AddInt64(&(wg.st), delta)
	if newSt < 0 {
		log.Fatalf("ERROR: status is lower than 0!!! (%v)", newSt)
	}
	return newSt
}

func (wg *WaitGroupWithState) Done() int64 {
	return wg.Add(-1)
}

func (wg *WaitGroupWithState) State() int64 {
	return atomic.LoadInt64(&(wg.st))
}

func (wg *WaitGroupWithState) Wait() {
	for wg.State() > 0 {
		runtime.Gosched()
	}
}

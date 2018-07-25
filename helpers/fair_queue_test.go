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
	"testing"
	"time"
)

func TestNewFairQueue(t *testing.T) {
	fq := NewFairQueue()

	prev := -1
	for i := 0; i < 100; i++ {
		ci := i
		fq.Enqueue(func() {
			if ci-prev != 1 {
				t.Errorf("#1 ci(%v) - prev(%v) != -1!!!", ci, prev)
			}
			time.Sleep(20 * time.Millisecond)
			if ci-prev != 1 {
				t.Errorf("#2 ci(%v) - prev(%v) != -1!!!", ci, prev)
			}
			prev = ci
		})
	}

	t.Logf("Close..")
	fq.Close()
	if prev != 99 {
		t.Errorf("#3 the latest on is NOT 99 (%v)", prev)
	}
}

func TestFairQueue_Close(t *testing.T) {
	fq := NewFairQueue()
	fq.Close()
}

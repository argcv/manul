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

/**
 * Fair Queue:
 * a fifo queue, which is used to help us schedule the
 * sequence of execution
 *
 * It could also treated as some kind of Lock
 */
type FairQueue struct {
	c  chan func()
	wg *WaitGroupWithState
	sd *SingletonDesc
}

func NewFairQueue() *FairQueue {
	return &FairQueue{
		c:  make(chan func()),
		wg: NewWaitGroupWithState(),
		sd: NewSingleton(),
	}
}

func (q *FairQueue) Perform() {
	go q.sd.Acquire(func() {
		for f := range q.c {
			f()
			if q.wg.Done() == 0 {
				return
			}
		}
	})
}

func (q *FairQueue) Enqueue(f func()) {
	q.wg.Add(1)
	q.Perform()
	q.c <- f
}

func (q *FairQueue) Wait() {
	q.Perform()
	q.wg.Wait()
}

// It's OK to leave a Go channel open forever and never close it.
// When the channel is no longer used, it will be garbage collected.
// -- <https://stackoverflow.com/questions/8593645>
// However we could provide a close + wait interface, which is used
// to indicate its finishing
func (q *FairQueue) Close() {
	close(q.c)
	q.Wait()
}

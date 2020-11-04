package promise

import (
	"fmt"
	"sync"
)

//Promise type
type Promise struct {
	mutex  sync.Mutex
	status Status
	result interface{}
	err    error
	wg     sync.WaitGroup
}

// Executor to create a Promise
type Executor func(Reject, Resolve)

// Resolve ...
type Resolve func(interface{})

// Reject ...
type Reject func(error)

//Status ...
type Status int

const (
	pending Status = iota
	rejected
	fulfilled
)

//Create new promise
func Create(executor Executor) *Promise {
	var promise = new(Promise)
	promise.mutex = sync.Mutex{}
	promise.wg = sync.WaitGroup{}
	promise.status = pending

	promise.wg.Add(1)

	go func() {
		defer promise.resolvePanic()
		executor(promise.reject, promise.resolve)
	}()

	return promise
}

func (promise *Promise) resolve(val interface{}) {
	promise.mutex.Lock()
	if promise.status == pending {
		switch result := val.(type) {
		case *Promise:
			result.wg.Wait()
			if result.err != nil {
				promise.status = rejected
				promise.err = result.err

			} else {
				promise.result = result.result
				promise.status = fulfilled
			}
		default:
			promise.result = result
			promise.status = fulfilled
		}
	}
	promise.wg.Done()
	promise.mutex.Unlock()
}

func (promise *Promise) reject(err error) {
	promise.mutex.Lock()
	if promise.status == pending {
		promise.err = err
		promise.status = rejected
	}
	promise.wg.Done()
	promise.mutex.Unlock()
}

func (promise *Promise) resolvePanic() {
	err := recover()
	if err != nil {
		switch e := err.(type) {
		case error:
			promise.reject(fmt.Errorf("panic recovery with error: %s", e.Error()))
		default:
			promise.reject(fmt.Errorf("panic recovery with unknown error: %s", fmt.Sprint(e)))
		}
	}
}

//Then ...
func (promise *Promise) Then(onFulfilled func(interface{}) interface{}, onRejected func(error) interface{}) *Promise {
	return Create(func(reject Reject, resolve Resolve) {
		promise.wg.Wait()
		if promise.err != nil {
			resolve(onRejected(promise.err))
		} else {
			if onFulfilled != nil {
				resolve(onFulfilled(promise.result))
			}
		}
	})
}

//Catch ...
func (promise *Promise) Catch(onRejected func(error) interface{}) *Promise {
	return promise.Then(nil, onRejected)
}

//Finally ...
func (promise *Promise) Finally(onSettled func() interface{}) *Promise {
	go func() {
		promise.wg.Wait()
		onSettled()
	}()
	return promise
}

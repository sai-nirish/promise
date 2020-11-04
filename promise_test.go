package promise

import (
	"errors"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	var promise = Create(func(reject Reject, resolve Resolve) {
		resolve("test")
	})
	if promise == nil {
		t.Error("Promise is nil")
	}
}

func TestPromise_Then(t *testing.T) {

	var promise = Create(
		func(reject Reject, resolve Resolve) {
			time.Sleep(3 * time.Second)
			resolve("test")
		})
	promise.Then(func(value interface{}) interface{} {
		var val = value.(string)
		if val != "test" {
			t.Error("Incorrect result propogated")
		}
		return val
	}, func(err error) interface{} {
		return err
	})
	promise.wg.Wait()
}

func TestPromise_ChaninedThen(t *testing.T) {

	var promise = Create(
		func(reject Reject, resolve Resolve) {
			time.Sleep(3 * time.Second)
			resolve("test")
		})
	promise.Then(func(value interface{}) interface{} {
		var val = value.(string)
		if val != "test" {
			t.Error("Incorrect result propogated")
		}
		return "chain"
	}, func(err error) interface{} {
		return err
	}).Then(func(value interface{}) interface{} {
		var val = value.(string)
		if val != "chain" {
			t.Error("Incorrect chain propogated")
		}
		return val

	}, func(err error) interface{} {
		return err
	})
	promise.wg.Wait()
}

func TestPromise_NestedThenPromise(t *testing.T) {
	var promise = Create(
		func(reject Reject, resolve Resolve) {
			time.Sleep(3 * time.Second)
			resolve("test")
		})

	promise.Then(func(value interface{}) interface{} {
		var prom = Create(
			func(reject Reject, resolve Resolve) {
				time.Sleep(3 * time.Second)
				var nested = "nested"
				resolve(nested)
			})
		return prom
	}, func(err error) interface{} {
		return err
	}).Then(func(value interface{}) interface{} {
		var val = value.(string)
		if val != "nested" {
			t.Error("Incorrect nesting propogation")
		}
		return val
	}, func(err error) interface{} {
		return err
	})
	promise.wg.Wait()
}
func TestPromise_Catch(t *testing.T) {
	var promise = Create(
		func(reject Reject, resolve Resolve) {
			time.Sleep(3 * time.Second)
			reject(errors.New("test"))
		})
	promise.Catch(func(err error) interface{} {
		if err.Error() != "test" {
			t.Error("Incorrect catch propogated")
		}
		return err

	})
	promise.wg.Wait()
}

func TestPromise_Then_After_Catch(t *testing.T) {
	var promise = Create(
		func(reject Reject, resolve Resolve) {
			time.Sleep(3 * time.Second)
			reject(errors.New("test"))
		})
	promise.Catch(func(err error) interface{} {
		if err.Error() != "test" {
			t.Error("Incorrect catch propogated")
		}
		return err
	}).Then(func(value interface{}) interface{} {
		var val = value.(error)
		if val.Error() != "test" {
			t.Error("Incorrect then after catch propogated")
		}
		return val
	}, func(err error) interface{} {
		return err
	})
	promise.wg.Wait()
}

func TestPromise_Catch_After_Then(t *testing.T) {

	var promise = Create(
		func(reject Reject, resolve Resolve) {
			time.Sleep(3 * time.Second)
			reject(errors.New("test"))
		})
	promise.Then(func(value interface{}) interface{} {
		var val = value.(string)
		if val != "test" {
			t.Error("Incorrect result propogated")
		}
		return val
	}, func(err error) interface{} {
		return err
	}).Catch(func(err error) interface{} {
		t.Error("Catch should not execute")
		return err
	})

	promise.wg.Wait()
}

func TestPromise_Finally(t *testing.T) {
	var promise = Create(
		func(reject Reject, resolve Resolve) {
			time.Sleep(3 * time.Second)
			resolve("test")
		})
	var finalPromise = promise.Then(func(value interface{}) interface{} {
		var val = value.(string)
		if val != "test" {
			t.Error("Incorrect result propogated")
		}
		return val
	}, func(err error) interface{} {
		return err
	}).Finally(
		func() interface{} {
			return nil
		})

	finalPromise.Then(func(value interface{}) interface{} {
		var val = value.(string)
		if val != "test" {
			t.Error("Incorrect finally resolution ")
		}
		return val
	}, func(err error) interface{} {
		return err
	})

	finalPromise.wg.Wait()

}

func TestPromise_onPanic(t *testing.T) {
	var promise = Create(
		func(reject Reject, resolve Resolve) {
			time.Sleep(2 * time.Second)
			panic(errors.New("panic"))
		})
	promise.Then(func(value interface{}) interface{} {
		return value
	}, func(err error) interface{} {
		if err.Error() != "panic" {
			t.Error("panic not handled")
		}
		return err
	})
	promise.wg.Wait()
}

## 0.0.3

* Fixed a number of bugs with the handling of `nil` values. If you called
  `d.Is(nil, d.Map(...))` you'd get a panic instead of a failure. The same
  applied to `d.Slice`, `d.Struct`, and the `FuncComparer` returned by
  `d.Func` and `d.NamedFunc`. If you called `d.Is(nil, ...)` it would
  sometimes consider the two argument not nil, if the second argument was a
  nil of a specific type (like you get if you write `var s *string` or `var s
  []int`).


## 0.0.2  2021-03-26

* Added a `d.Passes` method. This works just like `d.Is` except that the
  second argument must be a `detest.Comparer`. This reads better when using a
  `detest.FuncComparer`, since we're not testing exact equality in this case.
* Added a `d.Require` method. This is used something like this:

  ```go
  d.Require(d.Is(got, 42, "answer is 42"))
  ```
  
  If the `d.Is` test fails, then `t.Fatal` will be called and the current
  goroutine is aborted. This lets you exit a set of tests early if a key
  assertion fails.
* Change `d.Is` and `d.ValueIs` to take multiple final arguments. If only one
  argument is provided, this is used as the test name. If multiple are
  provided, then these are passed to `fmt.Sprintf` to create the test name.


## 0.0.1  2020-12-27

* First release upon an unsuspecting world.

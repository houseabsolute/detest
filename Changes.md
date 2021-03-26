## 0.0.2  2021-03-26

* Added a `d.Passes` method. This works just like `d.Is` except that the
  second argument must be a `detest.Comparer`. This reads better when using a
  `detest.FuncComparer`, since we're not testing exact equality in this case.


## 0.0.1  2020-12-27

* First release upon an unsuspecting world.

# Property-Based Testing and Functional Programming tools in Golang that use Go 1.18 Generics.

General Purpose tools
- Generic set operations on Golang arrays
- Generic Linked List
- Generic function composition
- Generic sort

A Properties-based testing library that is based upon ScalaCheck and Haskell Quickcheck
- Everything has strong typing thanks to Golang Generics
- Many different kinds of generators
- Every generator is composable. Lots of generators are included already and new ones are compositions of existing generators.
- A property is composable with other properties
- Assertions are composable with And and Or logic
- Programmers can better cover the scope of all possible inputs to a test(i.e. zero values, empty things, etc).
- Tests outcomes are reproducible
- Programmers can eliminate a lot of code duplication and get better tests at the same time because properties-based testing uses random test data.
- Properties-based testing is useful for all sorts of tests: unit, integration, course-grained functional/system/black-box.

A few properties-based testing library exist in Golang. [Gopter](https://github.com/leanovate/gopter/) is an example. This library has much less code than [Gopter](https://github.com/leanovate/gopter/) and provides a very important feature that Gopter does not, namely that all abstractions are fully composable.

### Two Key Abstractions:
- Generators - Generators are functions that produce random test data. 
  - They are composable. You can combine them to make other generators. 
generator.
  - They obey algebraic laws. You can guarantee the safety of their compositions.
  - They are pure functions, freely shareable between Go Routines.
  - Generators allow you to reproduce the exact same test data by passing in the same integer seed value into a SimpleRNG.  You should never need to save files of test data again. 
- Properties - Properties are functions that execute a predicate-like function over a set of test data generated using a given Generator. 
  - They are composable - You can combine them in arbitrary ways to make new properties. 
  - They obey algebraic laws so that you can guarantee the safety of their compositions.
  - They are pure functions, freely shareable between Go Routines.  
  - A property is a function that takes two function arguments:
    - one transforming a generated value into a potentially different type 
    - and one evaluating the preceding value for correctness. 
    - These two characteristics are a very important feature for test reuse.

### Initializing Project
- `go1.18rc1 mod init pkg`
- `go1.18rc1 mod tidy`
- `go1.18rc1 test -v -count=1 ./...` 
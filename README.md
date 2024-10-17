# Go interpreter for the "monkey" language

This is an implementation of an interpreter for the "monkey" language detailed in [Writing an Interpreter in Go](https://interpreterbook.com/)

## Language example

```
let age = 1;
let name = "Monkey";
let result = 10 * (20 / 2);
let myArray = [1, 2, 3, 4, 5];
let myHash = {"key1": "lorem ipsum", "key2": 42};

let fibonacci = fn(x) {
  if (x <= 1 ) {
    x
  } else {
    fibonacci(x - 1) + fibonacci(x - 2);
  }
};
```

Complicated Hello world example:
```
puts(
  ["Hello", 42][0],
  {"second word": "world"}["second word"],
  fn(b) { if (b) { "!" } else { "something else" } }(true)
)
```

## Usage

The interpreter is provided as a REPL, to run it use:

``` sh
go run github.com/ManuelGarciaF/go-interpreter@latest
```

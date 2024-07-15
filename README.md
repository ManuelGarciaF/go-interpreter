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
  if (x == 0) {
    0
  } else {
    if (x == 1) {
      1
    } else {
      fibonacci(x - 1) + fibonacci(x - 2);
    }
  }
};
```

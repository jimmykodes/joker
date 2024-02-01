# Joker
The Joker programming language.

Don't take programming so seriously...

---

## Installing

```shell
go install github.com/jimmykodes/joker/cmd/joker@latest
```

## Running

```shell
joker  # starts the repl

joker run fib.jk # compiles and runs the fib.jk file

joker build fib.jk # builds the fib.jk file into fib.jkb
joker run fib.jkb # runs the compiled fib.jkb file
```

---

# Language Spec

## Data Types

### Int

A integer is any numeric literal that does not contain a decimal point:

```
12 // => 12
```

#### Conversions

Floats and strings can be converted to integers using the `int` builtin

```
int(10) // => 10 - redundant cast of int -> int

// floats are truncated, not rounded
int(10.0) // => 10
int(10.1) // => 10
int(10.8) // => 10

// strings
int("10")   // => 10
int("10.1") // => 10
int("10.8") // => 10
```

### Float

A float is any numeric literal that contains a decimal point:

```
12.0 // => 12.0
0.05 // => 0.05
```

Note:
Float literals < 1 must include a preceding zero. `.5` is not valid, `0.5` is.

#### Conversions

Integers and strings can be converted to integers using the `float` builtin

```
float(10.0) // => 10.0 - redundant cast of float to float

// integers
float(10) // => 10.0

// strings
float("10")   // => 10.0
float("10.1") // => 10.1
```

### String

String literals are values contained in double quotes (`"`).

```
"Hello, world!"
"foo"
"bar"
"baz"
```

Note:
There is currently no handling of escaped quotes or alternate wrappers (like `'` or `\``).

#### Conversions

Integers and floats can be converted to strings using the `string` builtin

```
string("10") // => "10" - redundant cast of string to string

// integers
string(10)   // => "10"
string(9999) // => "9999"

// floats
string(10.0) //   => "10"
string(10.1) //   => "10.1"
string(10.959) // => "10.959"
```

### Boolean

Booleans are the values `true` and `false`

#### Conversions

Conversion to booleans can be done with the `!` unary operation. This will invert the bool, so doing it twice will
return the original "truthiness" of the value:

```
!10  // => false

!!10    // => true
!!10.0  // => true
!!"foo" // => true

!!0   // => false
!!0.0 // => false
!!""  // => false
```

### Array

Arrays are lists of elements between square brackets (`[` and `]`) where each element is separated by a comma(`,`).

```
[] // empty array
[1]
[1, 2]
["1", 2, 3.0] // arrays can contain multiple types
```

#### Element access

Array elements are access by index, starting at 0
```
[][0]                    // => error - index out of range
[1][0]                   // => 1
["foo", "bar", "baz"][1] // => "bar"
```

#### Element assignment

Element assignment is currently unsupported.
```
[1][0] = 12 // throws a parser error
```

### Map

Maps are lists of key value pairs between curly braces (`{` and `}`) where the keys are hashable values and the values are anything.
Keys and values should be separated by a colon(`:`) and pairs should be separated by a comma(`,`)
```
{} // empty map
{"1": 12}
{"1": 2, 3: 5} // keys and values can be mixed types
{"foo": [1, 2, 3], "bar": {"baz": "taco"}}
```

#### Element access

Map values are access by key:
```
{"foo": 12}["foo"] // => 12
{"foo": 12}["bar"] // => error key not present
```

#### Element assignment

Element assignment is currently unsupported.

```
{"foo": 1}["bar"] = 12 // throws a parser error
```

## Functions

Declarations:

```
fn add(a, b) {
    return a + b;
}
```

Anonymous:
```
fn(a + b) {
    return a + b;
}
```

Immediately invoked:
```
fn(a, b){ return a + b; }(12, 5);
```

Functions can also be values in Arrays and maps:
```
[
    fn(a, b) { return a + b; },
    fn(a, b) { return a - b; },
    fn(a, b) { return a * b; },
    fn(a, b) { return a / b; },
]

{
    "add": fn(a, b) { return a + b; }, 
    "sub": fn(a, b) { return a - b; }, 
    "mul": fn(a, b) { return a * b; },
    "div": fn(a, b) { return a / b; },
}
```

They can then be called directly from their access statements:
```
let x = [
    fn(a, b) { return a + b; },
    fn(a, b) { return a - b; },
    fn(a, b) { return a * b; },
    fn(a, b) { return a / b; },
]

let y = {
    "add": fn(a, b) { return a + b; }, 
    "sub": fn(a, b) { return a - b; }, 
    "mul": fn(a, b) { return a * b; },
    "div": fn(a, b) { return a / b; },
}

x[0](1, 2) // => 3
y["sub"](5, 3) // => 2
```

### Recursion

```
fn fib(i) {
    if i == 0 {
        return 0
    }
    if i == 1 {
        return 1
    }
    return fib(i-1) + fib(i-2);
}

fib(30); // => 832040
```

### Closures

Simple closures:
```
fn adder(a) {
    return fn(b) {
        return a + b;
    }
}

let plusTwo = adder(2);
plusTwo(4);  // => 6
plusTwo(10); // => 12
```

Accumulator closures:
```
fn accumulator() {
    acc := 0;
    return fn(a) {
        acc = acc + a;
        return acc
    }
}
let a = accumulator();
a(10); // => 10
a(10); // => 20
a(5); // => 25
```
or
```
fn add(a) {
    acc := 0;
    return fn() {
        acc = acc + a;
        return acc;
    }
}

let addTwo = add(2);
addTwo(); // => 2
addTwo(); // => 4
addTwo(); // => 6
```

Recursive closures _don't_ work (at least not yet):
```
fn map(arr, f) {
  fn iter(arr, acc) {
    if len(arr) == 0 {
      return acc;
    } else {
      next := arr[0];
      rest := slice(arr, 1, len(arr));
      return iter(rest, append(acc, f(next)));
    }
  }
  return iter(arr, []);
}

items := [1, 2, 3];
fn double(a) {
  return a * 2;
}
print(map(items, double));
```

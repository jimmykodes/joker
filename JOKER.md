# Joker Language Definition

## Style

These are the guiding principles for how Joker code should look and feel. The language
enforces many of these at the syntax level — where it doesn't, tooling (formatter/linter)
should.

### Naming

- **Classes** use `PascalCase`.
- **Everything else** — variables, functions, methods, constants — uses `camelCase`.
- **Package names** are single lowercase words. Multi-word package names (camelCase) are
  technically allowed but strongly discouraged.
- Visibility semantics are TBD. A Go-like capitalization approach is under consideration,
  in which case exported function names would be `PascalCase`.

### Indentation

Indentation is **tabs**, not spaces.

### Semicolons

Semicolons are not used as statement terminators. Statements are newline-terminated. The
only place semicolons appear is as separators in C-style `for` loop headers.

### Braces

Opening braces go on the **same line** as their keyword or closing parenthesis — never on
the next line. This applies to all constructs: `if`, `elseif`, `else`, `while`, `for`,
`switch`, `fn`, `class`, etc.

For multi-line declarations (e.g. long parameter lists), the opening brace follows the
closing parenthesis on the same line:

```
fn foo(
	bar,
	bazz,
	bing,
) {
	print(bar, bazz, bing)
}
```

The following is **invalid** — the brace must not appear on its own line:

```
// INVALID
if foo.bar
{
	print("foo")
}
```

### Parentheses Around Conditions

Conditions in `if`, `elseif`, `while`, and `switch` are **not** wrapped in parentheses.
Parentheses around the entire condition are invalid:

```
// INVALID
if (foo.bar) {
	print(foo)
}
```

Parentheses used for grouping _within_ a condition are fine:

```
// valid — parens add clarity to sub-expressions, not the whole condition
if (foo.bar < 12) && (foo.bing > 10) {
	print(foo)
}
```

### Trailing Commas

Trailing comma behavior follows Go's rules:

- **Single-line** — no trailing comma.
- **Multi-line** (closing delimiter on its own line) — trailing comma is **required**.

```
let x = [1, 2, 3]

let y = [
	1,
	2,
	3,
]
```

---

## Primitives

### Integers

Integers support decimal, binary (prefixed with `0b`), and hexadecimal (prefixed with `0x`) literal notations.

```
let x = 10
let y = 0b100
let z = 0xFFF
```

### Floats

Floats are decimal numbers with a required digit on both sides of the decimal point (i.e. `0.1` not `.1`).

```
let x = 10.0
let y = 0.1
```

### Strings

Strings are delimited by double quotes.

```
let x = "test"
```

### Byte

Essentially a uint8.

### Chars

Chars are an alias for a byte. They are delimited by single quotes and support escape sequences (e.g. `'\n'`).

```
let x = 't'
let y = '\n'
```

### Boolean

The two boolean literals are `true` and `false`.

```
let x = true
let y = false
```

### Nil

`nil` represents the absence of a value. A variable declared with `let` but no assignment is implicitly `nil`. It can also be explicitly assigned.

```
let x // inherently nil
let y = nil
```

---

## Variables

### Constants

Constants are declared with the `const` keyword. They are immutable bindings.

```
const x = 12
const y = "test"
```

### Top level variable declarations

Top-level variables are declared with `let`. Multiple declarations can be grouped with parentheses in a single `let` block.

```
let x = 12
let y = "test"
let (
    z = "foo"
    alpha = "bazz"
)
```

### Scoped variable declarations

Inside function (or block) bodies, the standard `let` declaration is still valid. Additionally, the short-hand `:=` declare-and-assign operator is available, but it **cannot** be used at the top level.

```
fn hello(name) {
    let x = 12 // still valid
    y := "test" // := equivalent, but cannot be used as a top level declaration
}
```

---

## Operators

**Arithmetic**
- `+` addition
  Valid for ints, floats, bytes, chars, and strings (concatenation for strings). When used on chars, the result is a numeric value (the sum of the byte values), and mixing a char with an int produces a numeric result as well.
  ```
  't' + 'a' // -> 213
  't' + 2 // -> 118
  ```
- `-` subtraction — valid for ints, floats, bytes, and chars.
- `*` multiplication — valid for ints, floats, bytes, and chars.
- `/` division — valid for ints, floats, bytes, and chars.
- `%` modulo — valid for ints, floats, bytes, and chars.
- `**` exponentiation — valid for ints, floats, bytes, and chars.

**Comparison**
- `==` equal
- `!=` not equal
- `<` less than
- `>` greater than
- `<=` less than or equal
- `>=` greater than or equal

**Logical**
- `!` NOT (unary prefix)
- `&&` logical AND
- `||` logical OR

**Bitwise**
- `&` AND
- `|` OR
- `^` XOR
- `<<` left shift
- `>>` right shift

**Assignment**
- `=` assign
- `:=` declare and assign (scoped only, not top-level)
- `+=`, `-=`, `*=`, `/=`, etc. — compound assignment operators

**Other**
- `--` / `++` — decrement / increment (postfix)
- `?.` — optional chaining operator. Evaluates a chain of member accesses left to right; if any link in the chain is `nil`, evaluation short-circuits and the entire expression returns `nil`.
- `...` — spread operator

---

## Strings

Strings support several forms:

- **Double-quoted strings** (`"..."`) — standard string literals.
- **Backtick strings** (`` ` `` ... `` ` ``) — template literals with interpolation. Expressions inside `{...}` are evaluated and inserted into the string.
- **Concatenation** via `+`.
- **Formatted strings** via a stdlib call like `strings.fmt`, which takes a format string with `%d`-style verbs and arguments.

```
empty = "" // empty string
x = "test" // normal
concat = "fizz" + "buzz" // -> "fizzbuzz"
interpolation = `your number was {guess}?` // guess = 12 -> "your number was 12"
fmt = strings.fmt("your number was %d", guess) // guess = 12 -> "your number was 12"
```

---

## Control Flow

### If

`if` / `elseif` / `else` blocks. Condition expressions are not wrapped in parentheses; bodies are always brace-delimited. The `elseif` keyword is a single word (not `else if`).

```
if condition {
    statements
} elseif condition {
    statements
} else {
    statements
}
```

### Switch

`switch` evaluates a condition and matches against `case` clauses. A `default` clause handles the fallback. Cases are terminated by `:`. Each case breaks implicitly — fallthrough is **not** the default behavior. To explicitly fall into the next case, use the `fallthrough` keyword.

```
switch condition {
    case "":
        fallthrough
    case "empty":
        print("empty or blank")
    default:
        print("non-empty")
}
```

---

## Loops

### C style for loop

A traditional three-part for loop with init (using `:=`), condition, and post statement. Body is brace-delimited.

```
for i := 0; i < 10; i++ {
    print(i)
}
```

### Go style range loop

`range` can iterate over arrays (yielding index and value), maps (yielding key and value), or a bare integer (yielding indices from 0 up to but not including that number, i.e. exclusive upper bound). Variables are declared with `:=`.

```
for i, val := range myArray {
    print("idx", i, "val", val)
}
for k, v := range myMap {
    print("key", k, "value", v)
}
for i := range 5 {
    print("idx", i)
}
```

### While loops

`while` loops repeat as long as their condition is true. If the condition is omitted, the loop runs indefinitely (equivalent to `while true`), and must be exited with `break`.

```
while {
    // equivalent to `while true`
    next, err := getThing()
    if err != nil { 
        break
    }
}

let numActive = 0
while numActive < 10 {
    thing := getThing()
    if thing.active {
        numActive++
    }
}
```

### Break & Continue

`break` exits the innermost enclosing loop. `continue` skips the rest of the current iteration and advances to the next iteration of the innermost enclosing loop. Both keywords apply only to the innermost loop scope — there is no labeled break/continue.

```
for i := range 10 {
    for j := range 100 {
        if j % 2 == 0 {
            continue // this will always continue the j loop, never the i loop
        }
        print("odd")
    }
}
```

### Iterators

The `iter()` builtin wraps a collection into a lazy iterator. Iterators support chainable functional methods like `.filter()`, `.map()`, and `.collect()`. Callbacks use pipe-delimited parameter syntax (`|param| expr`), which is a closure / lambda shorthand. Lambdas are single-expression only — `|x| x ** 2` is shorthand for `fn(x) { return x ** 2 }`. For multi-line logic, use a full anonymous function instead.

```
myArr = [1, 2, 3, 4, 5, 6]
evenSquared = iter(myArr).filter(|f| f % 2 == 0).map(|f| f ** 2).collect() // -> [4, 16, 36]
```

---

## Functions

Functions are declared with the `fn` keyword. Parameters are listed in parentheses without type annotations. Functions are first-class values — they can be assigned to variables, returned from other functions, and passed as arguments.

Named function declarations:

```
fn Hello(name) {
    print(`Hello, {name}`)
}
```

Anonymous function expressions assigned to a variable:

```
let hello = fn(name) {
    print(`Hello, {name}`)
}
```

Functions can return other functions, forming closures that capture variables from their enclosing scope:

```
fn closure(greeting) {
    return fn(name) {
        return `{greeting}, {name}!`
    }
}

let hello = closure("Hello")
hello("world") // -> "Hello, world!"
```

Recursion is supported:

```
fn recurFib(n) {
	if n <= 1 {
		return n
	}
	return recurFib(n-1) + recurFib(n-2)
}
```

### Lambdas

Lambdas are a shorthand for single-expression anonymous functions, using pipe-delimited parameters. `|x| expr` is equivalent to `fn(x) { return expr }`. They are single-expression only — for multi-line logic, use a full anonymous function.

```
let square = |x| x ** 2
square(5) // -> 25

let add = |a, b| a + b
add(1, 2) // -> 3

// lambdas are commonly used with iterators
iter(myArr).filter(|f| f % 2 == 0).map(|f| f ** 2).collect()

// for multi-line logic, use anonymous functions instead
iter(myArr).map(fn(f) {
    let result = f ** 2
    return result + 1
}).collect()
```

---

## Collections

- **Lists** are ordered, heterogeneous sequences delimited by `[...]`.
- **Maps** are key-value collections delimited by `{key: value, ...}` with string keys using colons.
Sets are not a builtin type. They may be provided via a `collections` stdlib package.

```
list = [1, 2, "thing"]
map = {"foo": "bar"}
```

---

## Classes

Classes are declared with the `class` keyword. Constructor parameters are listed in parentheses after the class name. Inside the constructor body, fields are assigned to `self`.

```
class Foo(bar) {
    self.bar = bar
}
```

**Static methods** are defined with `fn ClassName.methodName(...)` using dot notation. They can be called on the class itself or on an instance — static methods are always accessible from instances as well. They can also serve as alternative constructors (factory methods).

```
fn Foo.fromBing(bing) {
    return Foo(bing.bar)
}

fn Foo.string() {
    return "foo"
}
```

**Instance methods** are defined with `fn ClassName:methodName(...)` using colon notation. They have access to `self` and must be called on an instance. Calling an instance method on the class directly (e.g. `Foo:bazz(2)`) is invalid. A static method can also be called with an explicit instance as the first argument (e.g. `Foo.bazz(myFoo, 2)`), acting like an unbound method call. Defining an instance method with the same name as an existing static method (or vice versa) is invalid.

```
fn Foo:bazz(exp) {
    return self.bar ** exp
}

let myFoo = Foo(5)
print(myFoo.bar) // -> 5

print(myFoo.string()) // -> foo
print(Foo.string()) // -> foo

Foo:bazz(2) // invalid
myFoo:bazz(3) // -> 125
Foo.bazz(myFoo, 2) // -> 25
```

Note: `myFoo:bazz(3)` yields `125` (`5 ** 3`), while `Foo.bazz(myFoo, 2)` yields `25` (`5 ** 2`).

---

## Error Handling

Errors are values. They should be returned from functions as values rather than thrown as exceptions. *(Further details TBD.)*

---

## Scoping & Closures

Joker uses lexical scoping. Inner functions (closures) capture variables from their enclosing scope. Variables declared inside a block (e.g. an `if` body) are not visible outside that block.

```
fn closure(greeting) {
    return fn(name) {
        return `{greeting}, {name}!`
    }
}

let hello = closure("Hello")
hello("world") // -> "Hello, world!"
```

```
if x == 12 {
    z := x * 2
    print(z) // -> 24
}
print(z) // invalid. no z in scope
```

---

## Comments

Two comment styles are supported:

- **Line comments** — `//` through end of line.
- **Block comments** — `/* ... */`.

```
// comments go from here to the end of the line

/* comments go from here to the closing symbol */
```

---

## Modules & Imports

Packages are defined at the folder level (like Go). Every file must begin with a `package xxx` declaration.

Dependencies are imported via a `require(...)` block at the top of the file (below the package declaration), one import per line. Stdlib packages are listed first, separated from third-party packages by a blank line.

```
require(
    strings

    github.com/jimmykodes/jokesort
    github.com/jimmykodes/foo
)
```

Project configuration uses:
- **`Joker.toml`** — declares project dependencies outside the stdlib in a `[dependencies]` section.
- **`Joker.lock`** — stores SHAs of installed packages for reproducible builds.

There is no package registry; packages are fetched directly from GitHub (similar to Go modules and Deno).

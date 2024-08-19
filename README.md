<p>
  <img src="./dodo_banner.svg" alt="Dodo Programming Language">
</p>

Dodo Lang is an extended/modified version of the Monkey Programming Language and its interpreter implemented throughout the book [Writing An Interpreter In Go](https://interpreterbook.com/) by Thorsten Ball. Highly recommend it.

The following are additional features I've added or things I've changed in the language.

1. Immutable and mutable variables using the "mut" keyword, inspired by Rust.
1. Conditional loops using the "for" keyword, like in Go.
1. Ability to pick a character in a string by index.
1. Ability to quickly get the last element of an array or string by using index -1.
1. Ability to run .dodo files using the -f <filename> flag.
1. Ability to call built in functions on objects using dot syntax.
1. Ability to index arrays and hashmaps using dot syntax.
1. _(Work in progress)_ Pipe operator to pass result of one function to another directly after.
1. A typeof() function.
1. A preliminary debug() print function.
1. Improved the interactive mode/terminal REPL with some autocomplete, double parenthesis/bracket completion and colored output.

## Examples

### Variables

```rust
let foo = "Hello World!";
let bar = 10;
let baz = [1, 2, 3, 4];
let foobarbaz = {"country": "France", "capital": "Paris"};

let mut foo = 5;
let mut bar = "Snickers";
```

### If Statements

```rust
let foo = 8;

let result = if (10 > foo) {
    return true;
} else {
    return false;
}
```

### Functions

```rust
let add = fn(x, y) { x + y };

add(10, 5);

let newAdder = fn(x) { fn(y) { x + y }; };
let addTwo = newAdder(2);

addTwo(2);
```

### Built-in Functions and Dot Syntax

```rust
println("Hello World");
printf("Hello %s", "World");

len("How long am I?");
"How long am I?".len();

len([1, 2, 3, 4]);
[1, 2, 3].push(4);
[1, 2, 3].first();
[1, 2, 3].last();
[1, 2, 3].rest();

typeof([4, 5, 6])
```

### _(WIP)_ Pipe Operator

```rust
let add = fn(x, y) { x + y };
let sub = fn(x, y) { x - y };

sub(10, 3) |> add(5, $);
```

_[...] and more._

---

This project was just for learning and not intended to be production-ready in any way.

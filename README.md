### Simple Javascript Compiler using using Golang

The purpose of this project was to create a simple Javascript comiler from scratch using Golang. It's not effecient or optimized, and has very small lexicon, but it works!

It uses the code inside the file `src/example.lg` to output working javascript function.

#### Available Functionality

- Handles multiple function declarations
- Ability to pass in zero, one, or multiple arguments to a function
- Function body can evaluate Integers, function calls, and addition operations
- Additon Operations can evalulate Varibles or Integer values

#### Sample Source Code Syntax

```
def f(x,y)
    add(x,y)
end

def add(x, y)
    x + y
end
```

#### Sample Output

```
function f(x,y) { return add(x,y) };
function add(x,y) { return x+y };
function tacos() { return 1 };
console.log(f(1,2));
```

**NOTE:** Currently, there's no input to change values or call created functions. There is a placeholder `console.log(f(1,2))` to run the generated code.

#### Running The Code

To see the compiled JS code, from the root directory: `go run ./src/compiler.go`
To see the comipled JS code execute, from the root directory: `go run ./src/compiler.go | node`

#### #Next

- Add other operators
- Add ability to use function calls in operations

Thanks to https://www.destroyallsoftware.com/screencasts/catalog/a-compiler-from-scratch that was used as a guild, albight in Ruby.

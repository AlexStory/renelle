# Renelle

Renelle is a wip programming language for me to experimeent with writing a parser, interpreter, compiler, etc. It will probably never be ready for production use, but I can lay out my goals with the project and what I want in the language.

I'm going to try and design the language how I'd want my ideal language to be, since it is my language.

### Assignments
simple `let` binding, no need for semicolon.

`let x = 5`

### Functions

`fn` keyword declares a function. Followed by arguments, and function body in a `{}` block.

functions have an implicit return, but will allow an early return with a `return` keyword.

Call by using function name with `()`.

```
fn add(x y) {
    x + y
}

add(2 3)
```

`\` and `=>` are used for lambda functions, like so.

```
map(list, \x => x + 1)
```

### Conditionals

Simple `if` and `else` keywords, no parens, with a `{}` block.

```
if x > y {
    true
} else {
    false
}
```

I also will use the `and` and `or` keywords and am hoping to make them short circuit.

In the following example `run_code()` is never called.
```
if true or run_code() {
    return
}
```

### Pipelining

Renelle will have function pipelines with `|>` piping to the first argument.

The following would evaluate to 100.

```
2
|> add(8)
|> times(10)
```

### Comments

Comments will be started with a `#` and continue until the end of the line

```
# this is a comment
```

### Data Types

All of the simple data types are there.

```
"string" # string
123 # int
3.14 # float
true # boolean
```

a more unique one that renelle will have is atoms. which are simple pieces of data that are their value.

```
:ok # the value is :ok
```

For compound data types, we have objects, which are anonymous simple structures.

```
{
    name: "renelle"
    type: "dynamic"
}
```

As well as lists which is the collection data structure.

```
[1 2 3]
```

We will also have tuples, which when combined with atoms can represent values very well.

```
(:ok 100)
(:err "bad request")
```

### Modules

**WIP: This is not finalized and i'm still testing different things**

Modules can be declared at top level of a file, or they can be created with a block for scope.

```
module MyApp.Dog 
# the rest of the file will be included in the MyApp.Dog module

fn bark() {
    "woof"
}
```

```
module Script {
    fn do_work() {
        ...
    }
}

module Helper {
    fn util_function() {
        ...
    }
}
```

### Structs

**WIP: This is not finalized and i'm still testing differend things**

Each module can have a struct, which is like an object with predifined fields.

```
module MyApp.Dog

struct {
    name
    age
}

fn bark(dog) {
    $"Woof, my name is {dog.name}"
    # WIP string interpolation
}
```
in another file (or this one)

```
let dog = MyApp.Dog{
    name = "Fido"
    age = 3
}
```

I'm still undecided between only allowing functions to be called with full module resolution.

```
Dog.bark(dog)
```

or allowing both that and common dot notation.

```
Dog.bark(dog)
dog.bark() # both valid
```

#### Small Bits.

Renelle allows `?` in variable and function names, so you could have the following.

```
let my_list = [1 2 3]
List.empty?(my_list) # false
```

I have very early ideas for a in-code documentation system. early prototypes look like this.

```
##: Returns the sum of two numbers
##: example:
##:  let x = 1
##:  let y = 2
##:  add(x y)
##: returns: 3 
fn add(x y) {
    x + y
}
```


# Runny

With runny you can keep all your project's commands in one place.

## Getting Started

Runny's vocabulary is deliberately very simple. There are just 3 core keywords: `var`, `target` and `run` (and some other peripheral ones).

`var` is for defining variables. These can be static:
```
var {
    name "Tim"
}
```
or runny can capture a script's stdout:
```
var {
    name {
        run { echo "Tim" }
    }
}
```
Variables that you define are created as environment variables in the shell.

A `target` is for commands you want to run later:
```
target say_hello {
    run {
        echo "hello $name"
    }
}
```

The `run` keyword (in addition to being how you execute shell commands) is how your targets are executed.
```
run say_hello
```
You can also define scoped variables within your run statement.
```
run say_hello {
    var {
        name "Tim"
    }
}
```
A runny config can be extended from another using the `extends` keyword.
```
extends {
    "./parent.rny"
}
```

Lastly, can specify the shell you'd prefer to use with the `config` keyword.
```
config {
    shell "/bin/bash"
}
```

Use the `runny` executable in your terminal to specify targets to run:
```
runny {my_target}
```

By default runny will look for a `runny.rny` file or you can specify the path to a file with:
```
runny -f {path to file.rny}
```

## Ongoing Development
runny is the first (working) language I've written. Much of its inner workings are based on lox, from the wonderful book [Crafting Interpreters](https://craftinginterpreters.com).

If you'd like to contribute to runny please do raise an issue or open a pull request. I'm eager to improve the language.
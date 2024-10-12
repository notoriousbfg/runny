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
or runny can capture a script's stdout with a `run` statement:
```
var {
    name {
        run { echo "Tim" }
    }
}
```
Variables are transferred to your environment.


A `target` is for commands you want to run later:
```
target say_hello {
    run {
        echo "hello $name"
    }
}
```
Use the `desc` keyboard to describe your targets:
```
target say_hello {
    desc {
        "it says hello"
    }
    ...
}
```

The `run` keyword can be used to execute targets
```
run say_hello
```
or to run arbitrary shell commands, which can also be the output of run statements themselves.
```
run { echo "hello world" }
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

Lastly, one can specify the shell you'd prefer to use with the `config` keyword.
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
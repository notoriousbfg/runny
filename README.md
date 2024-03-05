# Runny

With runny you can keep all your project's commands in one place.

## Getting Started

Runny's vocabulary is deliberately very simple. There are just 3 terms.

`var` for defining variables. These can be static:
```
var {
    name "Tim"
}
```
or runny can capture a script's output:
```
var {
    name {
        run { echo "Tim" }
    }
}
```
Variables that you define are created as environment variables in the shell. Variables nested in targets and runs are scoped.

A `target` is for commands you want to run later:
```
target say_hello {
    run {
        echo "hello $name"
    }
}
```

`run` (which we've already seen) is how you run commands. You may also run targets:
```
run say_hello
```
or even define scoped variables:
```
run say_hello {
    var {
        name "Tim"
    }
}
```

Use the `runny` executable in your terminal to specify targets to run:
```
runny {my_target}
```

By default it will look for a `runny.rny` file, or you can specify the path to a file with:
```
runny -f {path to file.rny}
```
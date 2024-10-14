# Runny

Runny is a `Make` or `just` alternative. With runny you can keep all of your project's commands in one place.

You can then run them with `runny {?target}`
```
$ runny monitor_disk_usage -f examples/kitchensink.rny

> Monitor disk usage and warn if it exceeds 80%
if [ "$current_usage" -ge "$threshold" ]; then
    echo "Warning: Disk usage is at ${current_usage}%"
else
    echo "Disk usage is under control: ${current_usage}%"
fi
Disk usage is under control: 5%
```

Runny's vocabulary is deliberately very simple. There are just 3 core keywords: `var`, `target` and `run` (and some other peripheral ones).

`var` is for defining variables:
```
var {
    name "Tim"
}
```

A `target` contains things you want to run later:
```
target say_hello {
    run {
        echo "hello $name"
    }
}
```

`run` is for executing shell commands:
```
run {
    echo "hello world"
}
```

See the <a href="./examples/kitchensink.rny">kitchen sink</a> for some practical examples of all of the language's features.

## Config files
By default Runny looks for a `runny.rny` file in the current directory. If you want to use a different config file you can pass the `-f` flag.

## Editor Support
Syntax highlighting for Runny is currently supported in VSCode by installing the <a href="./editor/runny-0.0.1.vsix">editor/runny-0.0.1.vsix</a> file. Support will be added for other editors in the near future.

## Ongoing Development
Runny is the first (working) language I've written. Much of its inner workings are based on lox, from the wonderful book [Crafting Interpreters](https://craftinginterpreters.com).

If you'd like to contribute to my project please raise an issue or open a pull request. I'm eager to improve the language with the advice & experience of others.
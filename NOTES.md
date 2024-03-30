## 30/12/23

Looking at other task runner src, it looks like the lexer reads recipe vocabulary and recipe body differently using a long set of conditions.

Perhaps I should look at the sh library, since this handles multiple other shells.

It doesn't look like what I need. I should still probably split out runny and commands but I don't think the libraries do what I want.

For placeholders, I should choose a special syntax that doesn't exist in bash or zsh.

Or perhaps we use a shell syntax for placeholders, and then vars are exported as variables that bash/zsh/etc understands. Use sed? `${VAR}` or `$VAR`

I think I'll just treat anything inside a target or run as a string for now (excluding var overrides).

Why does Go even need to do the interpolation? We could export runny vars as environment variables, then the command would simple evaluate these in the (bash) shell. Simple.

For _documentation_ purposes later, we could parse the command string, but let's just lex it as a COMMAND_STRING and do nothing with it.

I would very much like to _bake in_ documentation later, like Elixir does.

## 30/03/24

I've gone down a particular avenue where the parser will look for variables defined an action's script and create them as variables. However sometimes a variable might have been declared in a script but not necessarily as a runny variable e.g. in a loop.

I think I should condense action and run statements into one. The interpreter will catch any variables declared in a run statement or target and the script will be a new field on the run statement. I _hope_ this will be easy to do.


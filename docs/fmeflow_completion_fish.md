## fmeflow completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	fmeflow completion fish | source

To load completions for every new session, execute once:

	fmeflow completion fish > ~/.config/fish/completions/fmeflow.fish

You will need to start a new shell for this setup to take effect.


```
fmeflow completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow completion](fmeflow_completion.md)	 - Generate the autocompletion script for the specified shell


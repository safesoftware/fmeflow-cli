## fmeflow completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(fmeflow completion zsh); compdef _fmeflow fmeflow

To load completions for every new session, execute once:

#### Linux:

	fmeflow completion zsh > "${fpath[1]}/_fmeflow"

#### macOS:

	fmeflow completion zsh > $(brew --prefix)/share/zsh/site-functions/_fmeflow

You will need to start a new shell for this setup to take effect.


```
fmeflow completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow completion](fmeflow_completion.md)	 - Generate the autocompletion script for the specified shell


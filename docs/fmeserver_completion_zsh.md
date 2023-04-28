## fmeserver completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(fmeserver completion zsh); compdef _fmeserver fmeserver

To load completions for every new session, execute once:

#### Linux:

	fmeserver completion zsh > "${fpath[1]}/_fmeserver"

#### macOS:

	fmeserver completion zsh > $(brew --prefix)/share/zsh/site-functions/_fmeserver

You will need to start a new shell for this setup to take effect.


```
fmeserver completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeserver-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeserver completion](fmeserver_completion.md)	 - Generate the autocompletion script for the specified shell


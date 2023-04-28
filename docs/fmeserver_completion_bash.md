## fmeserver completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(fmeserver completion bash)

To load completions for every new session, execute once:

#### Linux:

	fmeserver completion bash > /etc/bash_completion.d/fmeserver

#### macOS:

	fmeserver completion bash > $(brew --prefix)/etc/bash_completion.d/fmeserver

You will need to start a new shell for this setup to take effect.


```
fmeserver completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeserver-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeserver completion](fmeserver_completion.md)	 - Generate the autocompletion script for the specified shell


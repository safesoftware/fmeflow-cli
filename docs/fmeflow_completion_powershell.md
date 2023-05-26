## fmeflow completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	fmeflow completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
fmeflow completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow completion](fmeflow_completion.md)	 - Generate the autocompletion script for the specified shell


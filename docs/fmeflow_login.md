## fmeflow login

Save credentials for an FME Server

### Synopsis

Update the config file with the credentials to connect to FME Server. If just a URL is passed in, you will be prompted for a user and password for the FME Server. This will be used to generate an API token that will be saved to the config file for use connecting to FME Server.
	Use the --token flag to pass in an existing API token. To log in with a password on the command line without being prompted, place the password in a text file and pass that in using the --password-file flag.
	This will overwrite any existing credentials saved.

```
fmeflow login [URL] [flags]
```

### Examples

```

  # Prompt for user and password for the given FME Server URL  
  fmeflow login https://my-fmeflow.internal
	
  # Login to an FME Server using a pre-generated token
  fmeflow login https://my-fmeflow.internal --token 5937391ad3a87f19ba14dc6082867373087d031b
	
  # Login to an FME Server using a passed in user and password file (The password is contained in a file at the path /path/to/password-file). You will be prompted to enter the user and password if you did not pass it using the flags. However any flags in the group [user password-file] are set they must all be set. 
  fmeflow login https://my-fmeflow.internal --user admin --password-file /path/to/password-file
```

### Options

```
      --expiration int         The length of time to generate the token for in seconds. (default 2592000)
  -h, --help                   help for login
  -p, --password-file string   A file containing the FME Server password for the user to generate an API token for.
  -t, --token string           The existing API token to use to connect to FME Server
  -u, --user string            The FME Server user to generate an API token for.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow](fmeflow.md)	 - A command line interface for interacting with FME Flow.


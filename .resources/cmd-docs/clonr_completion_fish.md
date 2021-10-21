## clonr completion fish

generate the autocompletion script for fish

### Synopsis


Generate the autocompletion script for the fish shell.

To load completions in your current shell session:
$ clonr completion fish | source

To load completions for every new session, execute once:
$ clonr completion fish > ~/.config/fish/completions/clonr.fish

You will need to start a new shell for this setup to take effect.


```
clonr completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

* [clonr completion](clonr_completion.md)	 - generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 21-Oct-2021
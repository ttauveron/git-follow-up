package cmd

import (
	"git-follow-up/internal"
	"strings"
)

var bash_completion_func = `
__display_values()
{
    COMPREPLY=( $( compgen -W "`+strings.Join(internal.DisplayArgs, " ")+`" -- "$cur" ) )
}

__from_values()
{
	COMPREPLY=( $( compgen -W "`+strings.Join(internal.FromArgs, " ")+`" -- "$cur" ) )
}
`




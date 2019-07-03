package cmd

import (
	"github.com/ttauveron/git-follow-up/git"
	"strings"
)

var bash_completion_func = `
__display_values()
{
    COMPREPLY=( $( compgen -W "`+strings.Join(git.DisplayArgs, " ")+`" -- "$cur" ) )
}

__from_values()
{
	COMPREPLY=( $( compgen -W "`+strings.Join(git.FromArgs, " ")+`" -- "$cur" ) )
}
`




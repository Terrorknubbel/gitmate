package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func (c *Config) newMenuViewPreviewCmd() *cobra.Command {
	gitCmd := &cobra.Command{
		Use:   "preview",
		Short: "Aktuelle Menü Vorschau",
		RunE:  c.preview,
	}

	return gitCmd
}

func (c *Config) preview(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		//lint:ignore ST1005 Capitalization is fine…
		return errors.New("Keine Preview Option angegeben.")
	}

	option := args[0]
	var updateLine string

	if option == "Staging" {
		updateLine = "Aktualisiert den Feature, Staging & Master Branch"
	} else {
		updateLine = "Aktualisiert den Feature und Master Branch"
	}

	output :=
		`Bringt den aktuellen Feature Branch in den %s Branch.

Folgende Schritte werden durchgeführt:
↳ Überprüfung auf sauberen Branch Status (alle Commits gepusht, remote Status, …)
↳ %s
↳ Merged Master in den Feature Branch
↳ Merged den Feature Branch in den %s Branch

Im Falle von Merge Konflikten oder anderen Fehlern wird sofortig abgebrochen.
Es wird NICHT automatisch gepusht. Dies passiert erst nach einer manuellen Bestätigung.
`
	s := fmt.Sprintf(output, option, updateLine, option)
	fmt.Print(s)

	return nil
}

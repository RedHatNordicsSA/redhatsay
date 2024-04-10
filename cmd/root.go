/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/BourgeoisBear/rasterm"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/RedHatNordicsSA/redhatsay/assets"
)

var (
	vintage bool
	think   bool
	text    string
	err     error
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "redhatsay",
	Short: "Make the Red Hat say something",
	Long:  `redhatsay generates a Red Hat saying something provided by the user.`,
	Run: func(cmd *cobra.Command, args []string) {
		style := lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			Bold(true).
			Align(lipgloss.Center).
			BorderForeground(lipgloss.Color("#FF0000"))
		if len(args) > 0 {
			text = strings.Join(args, "\n")
		} else {
			if terminal.IsTerminal(0) {
				cobra.CheckErr(fmt.Errorf("no text provided"))
			}
			reader := cmd.InOrStdin()
			buf := new(strings.Builder)
			_, err = io.Copy(buf, reader)
			cobra.CheckErr(err)
			text = buf.String()
			// remove trailing newline
			text = strings.TrimSuffix(text, "\n")

		}
		fmt.Println(style.Render(text))
		styleSay := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)
		if think {
			fmt.Println(styleSay.Render("     o"))
			fmt.Print(styleSay.Render("      o"))
		} else {
			fmt.Println(styleSay.Render("     \\"))
			fmt.Print(styleSay.Render("      \\"))
		}
		// if vintage flag is set, use the vintage Red Hat logo
		// otherwise, use the modern Red Hat logo
		var file string
		if vintage {
			file = "RedHat_vintage.png"
		} else {
			file = "RedHat.png"
		}
		f, err := assets.FS.Open(file)
		cobra.CheckErr(err)
		defer f.Close()
		if rasterm.IsKittyCapable() {
			opts := rasterm.KittyImgOpts{
				PlacementId: 1,
			}
			err := rasterm.KittyCopyPNGInline(cmd.OutOrStdout(), f, opts)
			cobra.CheckErr(err)
		}
		fmt.Println("")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.redhatsay.yaml)")

	// --vintage flag to use the vintage Red Hat logo
	rootCmd.Flags().BoolVarP(&vintage, "vintage", "v", false, "Use the vintage Red Hat logo")
	// --think flag to make the Red Hat think instead of say
	rootCmd.Flags().BoolVarP(&think, "think", "t", false, "Make the Red Hat think instead of say")
}

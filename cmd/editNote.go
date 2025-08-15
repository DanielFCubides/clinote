/*
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 *
 * Copyright (C) Joakim Kennedy, 2016
 */

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/TcM1911/clinote"
	"github.com/spf13/cobra"
)

var editNoteCmd = &cobra.Command{
	Use:   "edit \"note title\"",
	Short: "Edit note.",
	Long: `
Edit allows you to edit the note. If no flags are set, the note is opened
with the editor defined by the environment variable $EDITOR.

The first line will be used as the note title and the rest is encoded as
the note content.

To change to title, the title flag can be used.

The note can be moved to another notebook by defining the new notebook
with the notebook flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		raw, err := cmd.Flags().GetBool("raw")
		if err != nil {
			fmt.Printf("❌ Invalid raw flag value: %v\n", err)
			fmt.Println("💡 Tip: Use --raw (no value needed) to edit XML content directly")
			return
		}
		title, err := cmd.Flags().GetString("title")
		if err != nil {
			fmt.Printf("❌ Failed to parse new title: %v\n", err)
			fmt.Println("💡 Tip: Use --title \"New Title\" or -t \"New Title\"")
			return
		}
		notebook, err := cmd.Flags().GetString("notebook")
		if err != nil {
			fmt.Printf("❌ Failed to parse notebook name: %v\n", err)
			fmt.Println("💡 Tip: Use --notebook \"Notebook Name\" or -b \"Notebook Name\"")
			return
		}
		recover, err := cmd.Flags().GetBool("recover")
		if err != nil {
			return
		}
		client := defaultClient()
		defer client.Close()
		ns, err := client.GetNoteStore()
		if err != nil {
			fmt.Printf("❌ Cannot connect to Evernote: %v\n", err)
			fmt.Println("💡 Troubleshooting:")
			fmt.Println("   • Check internet connection")
			fmt.Println("   • Verify authentication: clinote user login")
			fmt.Println("   • Check credentials: clinote user list")
			return
		}
		opts := clinote.DefaultNoteOption
		if raw {
			opts = opts | clinote.RawNote
		}
		if recover {
			c := clinote.NewClient(client.Config, client.Config.Store(), ns, clinote.DefaultClientOptions)
			err := clinote.EditNote(c, "", opts|clinote.UseRecoveryPointNote)
			if err != nil {
				fmt.Printf("❌ Failed to recover previous note: %v\n", err)
				fmt.Println("💡 Possible causes:")
				fmt.Println("   • No recovery point available")
				fmt.Println("   • Recovery file corrupted")
				fmt.Println("   • Storage permission issues")
				os.Exit(1)
			}
			return
		}
		if len(args) != 1 {
			fmt.Println("❌ Note identifier required")
			fmt.Println("💡 Usage: clinote note edit \"Note Title\"")
			fmt.Println("   • Use exact note title (case sensitive)")
			fmt.Println("   • Or use note index from: clinote note list")
			return
		}
		if title != "" {
			clinote.ChangeTitle(client.Config.Store(), ns, args[0], title)
		}
		if notebook != "" {
			clinote.MoveNote(client.Config.Store(), ns, args[0], notebook)
		}

		if title == "" && notebook == "" {
			c := clinote.NewClient(client.Config, client.Config.Store(), ns, clinote.DefaultClientOptions)
			err := clinote.EditNote(c, args[0], opts)
			if err != nil {
				fmt.Printf("❌ Failed to edit note: %v\n", err)
				fmt.Println("💡 Troubleshooting:")
				fmt.Println("   • Check if note exists: clinote note list --search \"title\"")
				fmt.Println("   • Verify editor: echo $EDITOR")
				fmt.Println("   • Check permissions and network connectivity")
				if strings.Contains(err.Error(), "not found") {
					fmt.Println("   • Note may have been deleted or moved")
				}
				os.Exit(1)
			}
		}
	},
}

func init() {
	noteCmd.AddCommand(editNoteCmd)
	editNoteCmd.Flags().StringP("title", "t", "", "Change the note title to.")
	editNoteCmd.Flags().StringP("notebook", "b", "", "Move the note to notebook.")
	editNoteCmd.Flags().Bool("raw", false, "Use raw content instead of markdown version.")
	editNoteCmd.Flags().Bool("recover", false, "Recover previous note that failed to save.")
}

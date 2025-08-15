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
	"strings"

	"github.com/TcM1911/clinote"
	"github.com/spf13/cobra"
)

var newNoteCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new note.",
	Long: `
New creates a new note. A title needs to be given for the
note.

If no notebook is given, the default notebook will be used.

The new note can be open in the $EDITOR by using the edit
flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		title, err := cmd.Flags().GetString("title")
		if err != nil {
			fmt.Printf("❌ Failed to parse note title: %v\n", err)
			fmt.Println("💡 Tip: Use --title \"Your Note Title\" or -t \"Your Note Title\"")
			return
		}
		edit, err := cmd.Flags().GetBool("edit")
		if err != nil {
			fmt.Printf("❌ Invalid edit flag value: %v\n", err)
			fmt.Println("💡 Tip: Use --edit or -e (no value needed)")
			return
		}
		if title == "" && !edit {
			fmt.Println("❌ Note title is required when not using edit mode")
			fmt.Println("💡 Options:")
			fmt.Println("   • Add a title: clinote note new --title \"My Note\"")
			fmt.Println("   • Use edit mode: clinote note new --edit")
			return
		}
		notebook, err := cmd.Flags().GetString("notebook")
		if err != nil {
			fmt.Printf("❌ Failed to parse notebook name: %v\n", err)
			fmt.Println("💡 Tip: Use --notebook \"Notebook Name\" or -b \"Notebook Name\"")
			fmt.Println("   • List available notebooks: clinote notebook list")
			return
		}
		raw, err := cmd.Flags().GetBool("raw")
		if err != nil {
			fmt.Printf("❌ Invalid raw flag value: %v\n", err)
			fmt.Println("💡 Tip: Use --raw (no value needed) to edit in XML format")
			return
		}
		createNote(title, notebook, edit, raw)
	},
}

func init() {
	noteCmd.AddCommand(newNoteCmd)
	newNoteCmd.Flags().StringP("title", "t", "", "Note title.")
	newNoteCmd.Flags().StringP("notebook", "b", "", "The notebook to save note to, if not set the default notebook will be used.")
	newNoteCmd.Flags().BoolP("edit", "e", false, "Open note in the editor.")
	newNoteCmd.Flags().Bool("raw", false, "Edit the content in raw mode.")
}

func createNote(title, notebook string, edit, raw bool) {
	c := newClient(clinote.DefaultClientOptions)
	defer c.Store.Close()

	note := new(clinote.Note)
	if title == "" {
		note.Title = "Untitled note"
	} else {
		note.Title = title
	}
	if notebook != "" {
		nb, err := clinote.FindNotebook(c.Store, c.NoteStore, notebook)
		if err != nil {
			fmt.Printf("❌ Notebook '%s' not found: %v\n", notebook, err)
			fmt.Println("💡 Available options:")
			fmt.Println("   • List notebooks: clinote notebook list")
			fmt.Println("   • Create new notebook: clinote notebook new \"Notebook Name\"")
			fmt.Println("   • Use default notebook: omit --notebook flag")
			return
		}
		note.Notebook = nb
	}
	opts := clinote.DefaultNoteOption
	if raw {
		opts |= clinote.RawNote
	}
	if edit {
		if err := clinote.CreateAndEditNewNote(c, note, opts); err != nil {
			fmt.Printf("❌ Failed to create and edit note: %v\n", err)
			fmt.Println("💡 Troubleshooting:")
			fmt.Println("   • Check if $EDITOR environment variable is set")
			fmt.Println("   • Verify editor is installed and accessible")
			fmt.Println("   • Try: export EDITOR=vim (or nano, code, etc.)")
			fmt.Println("   • Check network connection for Evernote sync")
			if strings.Contains(err.Error(), "recovery") {
				fmt.Println("   • Recovery available: clinote note edit --recover")
			}
		}
		return
	}
	err := clinote.SaveNewNote(c.NoteStore, note, raw)
	if err != nil {
		fmt.Printf("❌ Failed to save note: %v\n", err)
		fmt.Println("💡 Troubleshooting:")
		fmt.Println("   • Check network connection")
		fmt.Println("   • Verify authentication: clinote user login")
		fmt.Println("   • Check account quota and permissions")
	}
}

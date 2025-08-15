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

	"github.com/TcM1911/clinote"
	"github.com/spf13/cobra"
)

var deleteNoteCmd = &cobra.Command{
	Use:   "delete \"note title\"",
	Short: "Delete note.",
	Long: `Moves the note into the trash. The note may still be undeleted, unless it is expunged.
To expunge the note you need to use the official client or the web client.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("❌ Note identifier required")
			fmt.Println("💡 Usage: clinote note delete \"Note Title\"")
			fmt.Println("   • Use exact note title (case sensitive)")
			fmt.Println("   • Or use note index from: clinote note list")
			return
		}
		nb, err := cmd.Flags().GetString("notebook")
		if err != nil {
			fmt.Printf("❌ Invalid notebook parameter: %v\n", err)
			fmt.Println("💡 Tip: Use --notebook \"Notebook Name\" to specify source")
			return
		}
		client := defaultClient()
		defer client.Close()
		ns, err := client.GetNoteStore()
		if err != nil {
			return
		}
		err = clinote.DeleteNote(client.Config.Store(), ns, args[0], nb)
		if err != nil {
			fmt.Printf("❌ Failed to delete note: %v\n", err)
			fmt.Println("💡 Possible causes:")
			fmt.Println("   • Note not found or already deleted")
			fmt.Println("   • Network connectivity issues")
			fmt.Println("   • Insufficient permissions")
			fmt.Println("   • Note is being edited elsewhere")
			os.Exit(1)
		}
	},
}

func init() {
	noteCmd.AddCommand(deleteNoteCmd)
	deleteNoteCmd.Flags().StringP("notebook", "b", "", "The notebook of the note.")
}

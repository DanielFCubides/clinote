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

var editNotebookCmd = &cobra.Command{
	Use:   "edit \"notebook name\"",
	Short: "Edit a notebook.",
	Long: `
Edit a notebook. The notebook's name can be changed using the
name flag.

To move the notebook to another stack, use the stack flag to
define the new stack.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("❌ Notebook name required")
			fmt.Println("💡 Usage: clinote notebook edit \"Notebook Name\"")
			fmt.Println("   • List notebooks: clinote notebook list")
			return
		}
		change := false
		notebook := new(clinote.Notebook)
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Printf("❌ Invalid new name parameter: %v\n", err)
			fmt.Println("💡 Tip: Use --name \"New Notebook Name\"")
			return
		}
		if name != "" {
			notebook.Name = name
			change = true
		}

		stack, err := cmd.Flags().GetString("stack")
		if err != nil {
			fmt.Printf("❌ Invalid stack parameter: %v\n", err)
			fmt.Println("💡 Tip: Use --stack \"Stack Name\" to organize notebooks")
			return
		}
		if stack != "" {
			notebook.Stack = stack
			change = true
		}

		if !change {
			fmt.Println("⚠️  No changes specified")
			fmt.Println("💡 Available options:")
			fmt.Println("   • Change name: --name \"New Name\"")
			fmt.Println("   • Change stack: --stack \"Stack Name\"")
			return
		}
		client := defaultClient()
		defer client.Close()
		ns, err := client.GetNoteStore()
		if err != nil {
			return
		}
		err = clinote.UpdateNotebook(client.Config.Store(), ns, args[0], notebook)
		if err != nil {
			fmt.Printf("❌ Failed to update notebook: %v\n", err)
			fmt.Println("💡 Possible causes:")
			fmt.Println("   • New name conflicts with existing notebook")
			fmt.Println("   • Network connectivity issues")
			fmt.Println("   • Insufficient permissions")
			fmt.Println("   • Notebook not found")
			os.Exit(1)
		}
		fmt.Println("✅ Notebook updated successfully")
	},
}

func init() {
	notebookCmd.AddCommand(editNotebookCmd)
	editNotebookCmd.Flags().StringP("name", "n", "", "Change notebook name to.")
	editNotebookCmd.Flags().StringP("stack", "s", "", "Change notebook stack to.")
}

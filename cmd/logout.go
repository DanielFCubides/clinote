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

	"github.com/TcM1911/clinote/evernote"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout user.",
	Long: `
Logs a user out by removing the session token.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := defaultClient()
		defer client.Close()
		err := evernote.Logout(client.GetConfig())
		if err != nil {
			fmt.Printf("❌ Logout failed: %v\n", err)
			fmt.Println("💡 Troubleshooting:")
			fmt.Println("   • Check if you're currently logged in")
			fmt.Println("   • Try: clinote user list")
			fmt.Println("   • Verify config permissions")
			return
		}
		fmt.Println("✅ Successfully logged out")
	},
}

func init() {
	userCmd.AddCommand(logoutCmd)
}

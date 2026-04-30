package profiles

import "github.com/spf13/cobra"

var ProfilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "Manage and explore profiles",
}

func init() {
	ProfilesCmd.AddCommand(listCmd)
	ProfilesCmd.AddCommand(getCmd)
	ProfilesCmd.AddCommand(searchCmd)
	ProfilesCmd.AddCommand(createCmd)
	ProfilesCmd.AddCommand(exportCmd)
}
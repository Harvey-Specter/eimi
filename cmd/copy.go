/*
Copyright © 2024 lihong

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Harvey-Specter/eimi/svc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// copyCmd represents the copy command
// var Append string
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("copy called")
		fmt.Printf("Config: %v\n", viper.AllSettings())
		// fmt.Printf("Source: %v\n", Append)
		allMap := viper.AllSettings()
		// svc.GetRecord(allMap["src"].(map[string]any))
		srcMap := allMap["src"].(map[string]any)
		destMap := allMap["dest"].(map[string]any)
		tables := svc.GetDDL(srcMap)
		_, err := svc.ExecCopy(srcMap, destMap, tables)
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
	// copyCmd.Flags().StringVarP(&printFlag, "file", "f", "", "print flag for local")
	copyCmd.PersistentFlags().StringVarP(&cfgFile, "file", "f", "", "config file (default is $HOME/.cog.yaml)")
	if err := copyCmd.MarkPersistentFlagRequired("file"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//copyCmd.Flags().StringVarP(&Append, "append", "a", "", "append data on exists table")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// copyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// copyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

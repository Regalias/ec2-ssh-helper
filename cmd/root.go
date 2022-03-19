/*
Copyright © 2022 Regalias

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

	"github.com/regalias/ec2-ssh-helper-go/cmd/essh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ec2-ssh-helper-go",
	Short: "AWS EC2 SSH helper tool",
	// 	Long: `A longer description that spans multiple lines and likely contains
	// examples and usage of using your application. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		region := viper.GetString("region")
		useBastion := viper.GetBool("bastion")
		username := viper.GetString("username")
		port := viper.GetUint("port")
		essh.Entrypoint(region, useBastion, username, port)
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ec2-ssh-helper-go.yaml)")

	rootCmd.PersistentFlags().StringP("region", "r", "", "AWS region")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("bastion", "b", false, "Use bastion host")
	rootCmd.Flags().StringP("username", "u", "", "SSH username")
	rootCmd.Flags().UintP("port", "p", 22, "SSH port")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// Default values
	viper.SetEnvPrefix("essh")

	viper.BindEnv("region", "AWS_REGION", "AWS_DEFAULT_REGION", "ESSH_REGION")
	viper.BindPFlag("region", rootCmd.Flags().Lookup("region"))

	viper.BindPFlag("bastion", rootCmd.Flags().Lookup("bastion"))
	viper.BindPFlag("username", rootCmd.Flags().Lookup("username"))
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))

	viper.SetDefault("bastion", false)
	viper.SetDefault("username", "ec2-user")
	viper.SetDefault("port", uint(22))

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".ec2-ssh-helper-go" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ec2-ssh-helper-go")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

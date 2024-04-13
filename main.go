package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var opts struct {
	StructName    string
	InterfaceName string
	OutputFile    string

	Debug bool
}

func main() {
	cmd := &cobra.Command{
		Use:           "ifacegen [dir]",
		Args:          cobra.MaximumNArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) == 0 {
				filename, ok := os.LookupEnv("GOFILE")
				if !ok {
					return fmt.Errorf("must provide package directory")
				}
				dir, err := os.Getwd()
				if err != nil {
					return err
				}
				data, err := os.ReadFile(filepath.Join(dir, filename))
				if err != nil {
					return err
				}
				opts.StructName = getStructName(data)
				args = append(args, dir)
			}

			if opts.StructName == "" {
				return fmt.Errorf("must provide StructName")
			}
			if opts.OutputFile == "" {
				opts.OutputFile = fmt.Sprintf("zz_%s.iface.go", opts.StructName)
			}
			if opts.Debug {
				fmt.Printf("PackageDir: %s\n", args[0])
				fmt.Printf("StructName: %s\n", opts.StructName)
				fmt.Printf("OutputFile: %s\n", opts.OutputFile)
			}

			files, err := ParsePackage(args[0], opts.StructName)
			if err != nil {
				return err
			}
			data, err := GenerateFile(files)
			if err != nil {
				return err
			}
			return os.WriteFile(filepath.Join(args[0], opts.OutputFile), data, 0755)
		},
	}

	cmd.Flags().StringVar(&opts.StructName, "struct", "", "name of struct to generate interface for")
	cmd.Flags().StringVar(&opts.InterfaceName, "iface", "Interface", "name of generated interface")
	cmd.Flags().StringVarP(&opts.OutputFile, "output", "o", "", "name of output file")
	cmd.Flags().BoolVar(&opts.Debug, "debug", false, "")

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

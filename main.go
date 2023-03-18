package main

import (
	"fmt"
	"path/filepath"

	"log"
	"sort"

	"github.com/spf13/viper"
	"github.com/x-mod/build"
	"github.com/x-mod/cmd"
	"github.com/x-mod/dir"

	"github.com/signintech/gopdf"
)

func main() {
	cmd.Version(build.String())
	cmd.Add(
		cmd.Name("images2pdf"),
		cmd.Short("images2pdf - images convert to pdf"),
		cmd.Main(Images2PDF),
	)
	cmd.PersistentFlags().StringP("input-folder", "i", "", "input folder path")
	cmd.PersistentFlags().StringP("output-file", "o", "", "output file path")
	cmd.PersistentFlags().IntP("verbose", "v", 0, "log verbosity level")
	cmd.Execute()
}

func Images2PDF(cmd *cmd.Command, args []string) error {
	input := dir.New(dir.Root(viper.GetString("input-folder")))
	if err := input.Open(); err != nil {
		return fmt.Errorf("failed to open input folder: %w", err)
	}

	files, err := input.Files()
	if err != nil {
		return fmt.Errorf("failed to get files: %w", err)
	}
	sort.Strings(sort.StringSlice(files))

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	for _, file := range files {
		log.Println("file => ", file)

		if !isImage(file) {
			log.Println("file ignored => ", file)
			continue
		}

		pdf.AddPage()
		if err := pdf.Image(input.Path(file), 0, 0, nil); err != nil {
			fmt.Printf("Error reading %s: %v\n", input.Path(file), err)
			continue
		}
	}

	pdf.WritePdf(viper.GetString("output-file"))
	fmt.Printf("PDF file saved as %s.\n", viper.GetString("output-file"))
	return nil
}

func isImage(filename string) bool {
	switch filepath.Ext(filename) {
	case ".png", ".jpg", ".jpeg", ".gif", ".tiff", ".bmp":
		return true
	}
	return false
}

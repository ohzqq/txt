package cmd

import (
	"bytes"
	"log"
	"os"
	"strings"

	"github.com/ohzqq/txt"
	"github.com/spf13/cobra"
)

var (
	font       string
	input      string
	paginate   bool
	dimensions []int
	wrapper    = txt.NewWrapper()
)

// wrapCmd represents the wrap command
var wrapCmd = &cobra.Command{
	Use:   "wrap",
	Short: "wrap text",
	Run: func(cmd *cobra.Command, args []string) {
		switch len(dimensions) {
		case 2:
			wrapper.SetHeight(dimensions[1])
			fallthrough
		case 1:
			wrapper.SetWidth(dimensions[0])
		}
		var ttf []byte
		var err error
		switch {
		case font == "mono":
			//wrapper.WithTTF(txt.GoMono, wrapper.FontSize)
		case font == "regular":
			//wrapper.WithTTF(txt.GoRegular, wrapper.FontSize)
		case font != "":
			ttf, err = os.ReadFile(font)
			if err != nil {
				log.Fatal(err)
			}
			err = wrapper.ParseTTF(ttf, wrapper.FontSize)
			if err != nil {
				log.Fatal(err)
			}
		}
		var buf bytes.Buffer
		if input != "" {
			d, err := os.Open(input)
			if err != nil {
				log.Fatal(err)
			}
			defer d.Close()
			_, err = buf.ReadFrom(d)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			in := os.Stdin
			fi, err := in.Stat()
			if err != nil {
				log.Fatal(err)
			}
			if fi.Mode()&os.ModeNamedPipe == 0 {
				for _, arg := range args {
					f, err := os.Open(arg)
					if err != nil {
						log.Fatal(err)
					}
					defer f.Close()
					_, err = buf.ReadFrom(f)
					if err != nil {
						log.Fatal(err)
					}
				}
			} else {
				_, err := buf.ReadFrom(os.Stdin)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		if paginate {
			pages := wrapper.Paginate(buf.String())
			for i := range pages.TotalPages {
				println(i)
				println(strings.TrimSpace(strings.Join(pages.Current(), "\n")))
				pages.NextPage()
			}
			//fmt.Printf("%s\n", strings.Join(pages.Current(), "\n"))
			//for _, page := range pages.AllPages() {
			//fmt.Printf("%#v\n", strings.Join(page, ""))
			//}
		} else {
			//lines := wrapper.WrapText(buf.String())
			//fmt.Printf("%#v\n", lines)
		}

	},
}

func init() {
	rootCmd.AddCommand(wrapCmd)
	wrapCmd.PersistentFlags().StringVarP(&font, "font", "f", "", "a ttf font to use")
	wrapCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "input file to read")
	wrapCmd.PersistentFlags().IntVarP(&wrapper.FontSize, "font-size", "s", 16, "the font size")
	wrapCmd.PersistentFlags().BoolVarP(&paginate, "paginate", "p", false, "paginate the result")
	wrapCmd.PersistentFlags().IntVarP(&wrapper.Width, "width", "w", 250, "max width for wrapping")
	wrapCmd.PersistentFlags().IntVar(&wrapper.Height, "height", 100, "max height for wrapping")
	wrapCmd.PersistentFlags().IntSliceVarP(&dimensions, "box", "b", []int{}, "dimensions for the wrapping box: w,h")
	wrapCmd.PersistentFlags().IntVarP(&wrapper.MaxLines, "max-lines", "m", 1, "max lines for wrapping")
	wrapCmd.PersistentFlags().IntVar(&wrapper.DPI, "dpi", 72, "the dpi for the font")
}

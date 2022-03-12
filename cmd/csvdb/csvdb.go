package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/stepan2volkov/csvdb/internal/app"
	"github.com/stepan2volkov/csvdb/internal/app/table/formatter"
	"github.com/stepan2volkov/csvdb/internal/app/table/loader"
)

func main() {
	a := app.NewApp()
	f := formatter.DefaultFormatter{}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("~# ")

		openedBuffer := scanner.Scan()
		if !openedBuffer {
			return
		}
		in := scanner.Text()
		if in == "" {
			continue
		}
		switch {
		case in == `\q`:
			return
		case strings.HasPrefix(in, `\load`):
			in = strings.TrimSpace(strings.TrimPrefix(in, `\load`))
			args := strings.Split(in, " ")
			if len(args) != 2 {
				fmt.Println("\tThe right syntax: \n\t\\load <csv-path> <config-path>")
				continue
			}
			t, err := loader.LoadFromCSV(args[0], args[1])
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
				continue
			}
			if err = a.LoadTable(t); err != nil {
				fmt.Printf("ERROR: %v\n", err)
			}
		case in == `\list`:
			fmt.Println(strings.Join(a.TableList(), "\n"))
		case in == `\help`:
		case strings.HasPrefix(in, `\drop`):
			newTableName := strings.TrimSpace(strings.TrimPrefix(in, `\drop`))
			if newTableName == "" {
				fmt.Println("\tThe right syntax: \n\t\\drop <tablename>")
				continue
			}
			if err := a.DropTable(newTableName); err != nil {
				fmt.Printf("\t%v\n", err)
			}
		default:
			res, err := a.Execute(context.Background(), in)
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
				continue
			}
			fmt.Println(f.Format(res))
		}
	}
}

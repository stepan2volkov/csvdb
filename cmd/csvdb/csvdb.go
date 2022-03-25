package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/stepan2volkov/csvdb/internal/app"
	"github.com/stepan2volkov/csvdb/internal/app/table/formatter"
	"github.com/stepan2volkov/csvdb/internal/app/table/loader"
	"github.com/stepan2volkov/csvdb/internal/logger"
	"go.uber.org/zap"
)

func main() {
	log := logger.GetLogger()
	log.Info("starting csv-db")
	a := app.NewApp()
	f := formatter.DefaultFormatter{}

	scanner := bufio.NewScanner(os.Stdin)

	log.Info("csv-db has been ready to accept queries")
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
				fmt.Printf("wrong syntax for \\load: '%s'\n", in)
				continue
			}
			t, err := loader.LoadFromCSV(args[0], args[1])
			if err != nil {
				fmt.Printf("error when loading from csv: %v\n", err)
				continue
			}
			if err = a.LoadTable(t); err != nil {
				fmt.Printf("error when loading table: %v\n", err)
				log.Error("error when loading table",
					zap.Error(err))
			}
		case in == `\list`:
			fmt.Println(strings.Join(a.TableList(), "\n"))
		case in == `\help`:
		case strings.HasPrefix(in, `\drop`):
			newTableName := strings.TrimSpace(strings.TrimPrefix(in, `\drop`))
			if newTableName == "" {
				fmt.Printf("wrong syntax for \\drop: '%s'\n", in)
				continue
			}
			if err := a.DropTable(newTableName); err != nil {
				fmt.Printf("error when dropping table: '%v'\n", err)
				log.Error("error when dropping table",
					zap.Error(err))
			}
		default:
			start := time.Now()
			res, err := a.Execute(context.Background(), in)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				log.Error("error executing query",
					zap.String("query", in),
					zap.Error(err))
				continue
			}
			duration := time.Since(start)
			fmt.Println(f.Format(res))
			log.Info("query has been executed",
				zap.String("query", in),
				zap.String("table", res.Name),
				zap.Duration("duration", duration))
		}
	}
}

package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/stepan2volkov/csvdb/internal/app"
	"github.com/stepan2volkov/csvdb/internal/app/table/formatter"
	"github.com/stepan2volkov/csvdb/internal/app/table/loader"
	"github.com/stepan2volkov/csvdb/internal/logger"
	"go.uber.org/zap"
)

var (
	a   *app.App
	f   formatter.DefaultFormatter
	log *zap.Logger
)

func readStdOut() <-chan string {
	ret := make(chan string, 1)
	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		for {
			openedBuffer := scanner.Scan()
			if !openedBuffer {
				return
			}
			ret <- scanner.Text()
		}
	}()
	return ret

}

func handleInput(ctx context.Context, in string) {
	switch {
	case in == `\q`:
		return
	case strings.HasPrefix(in, `\load`):
		in = strings.TrimSpace(strings.TrimPrefix(in, `\load`))
		args := strings.Split(in, " ")
		if len(args) != 2 {
			fmt.Printf("wrong syntax for \\load: '%s'\n", in)
		}
		t, err := loader.LoadFromCSV(args[0], args[1])
		if err != nil {
			fmt.Printf("error when loading from csv: %v\n", err)
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
		}
		if err := a.DropTable(newTableName); err != nil {
			fmt.Printf("error when dropping table: '%v'\n", err)
			log.Error("error when dropping table",
				zap.Error(err))
		}
	default:
		start := time.Now()
		res, err := a.Execute(ctx, in)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			log.Error("error executing query",
				zap.String("query", in),
				zap.Error(err))
			return
		}
		duration := time.Since(start)
		fmt.Println(f.Format(res))
		log.Info("query has been executed",
			zap.String("query", in),
			zap.String("table", res.Name),
			zap.Duration("duration", duration))
	}
}

func main() {
	log = logger.GetLogger()

	log.Info("starting csv-db")
	a = app.NewApp()
	f = formatter.DefaultFormatter{}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	reader := readStdOut()

	log.Info("csv-db has been ready to accept queries")
	for {
		fmt.Print("~# ")
		select {
		case <-ctx.Done():
			fmt.Println("Bye-bye!")
			log.Info("staring gracefull shutdown")
			return
		case in := <-reader:
			if in == "" {
				continue
			}
			handleInput(ctx, in)
		}
	}
}

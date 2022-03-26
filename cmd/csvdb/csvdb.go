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

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/stepan2volkov/csvdb/internal/app"
	"github.com/stepan2volkov/csvdb/internal/app/table"
	"github.com/stepan2volkov/csvdb/internal/app/table/formatter"
	"github.com/stepan2volkov/csvdb/internal/app/table/loader"
)

const (
	cmdLoadTable  = `\load`
	cmdTableList  = `\list`
	cmdDroupTable = `\drop`
	cmdHelp       = `\help`
	welcomeQuery  = "~# "
)

var (
	helpList = []struct {
		cmd  string
		desc string
	}{
		{cmd: cmdHelp, desc: "Show the help"},
		{cmd: cmdTableList, desc: "Show available loaded tables"},
		{cmd: cmdLoadTable, desc: fmt.Sprintf("Load the table. Format: '%s <csv-path> <yaml-description-path>'", cmdLoadTable)},
		{cmd: cmdDroupTable, desc: fmt.Sprintf("Drop the table. Format: '%s <tablename>'", cmdDroupTable)},
	}
)

func getLogger() *zap.Logger {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	accessEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level",
		TimeKey:     "timestamp",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime:  zapcore.ISO8601TimeEncoder,
	})
	errorEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level",
		TimeKey:     "timestamp",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime:  zapcore.ISO8601TimeEncoder,
	})

	errorFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		msg := fmt.Sprintf("error when createing error.log file: %v", err)
		panic(msg)
	}
	errorSync := zapcore.AddSync(errorFile)

	accessFile, err := os.OpenFile("access.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		msg := fmt.Sprintf("error when createing access.log file: %v", err)
		panic(msg)
	}
	accessSync := zapcore.AddSync(accessFile)

	core := zapcore.NewTee(
		zapcore.NewCore(accessEncoder, accessSync, lowPriority),
		zapcore.NewCore(errorEncoder, errorSync, highPriority),
	)

	return zap.New(core)
}

func readStdOut(ctx context.Context) <-chan string {
	ret := make(chan string, 1)
	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ret)
				return
			default:
				openedBuffer := scanner.Scan()
				if !openedBuffer {
					close(ret)
					return
				}
				ret <- scanner.Text()
			}
		}
	}()
	return ret

}

func handleInput(ctx context.Context, logger *zap.Logger, a *app.App, f table.Formatter, in string) {
	switch {
	case in == `\q`:
		return
	case strings.HasPrefix(in, cmdLoadTable):
		in = strings.TrimSpace(strings.TrimPrefix(in, cmdLoadTable))
		args := strings.Split(in, " ")
		if len(args) != 2 {
			fmt.Printf("wrong syntax for %s: '%s'\n", cmdLoadTable, in)
		}
		t, err := loader.LoadFromCSV(args[0], args[1])
		if err != nil {
			fmt.Printf("error when loading from csv: %v\n", err)
		}
		if err = a.LoadTable(t); err != nil {
			fmt.Printf("error when loading table: %v\n", err)
			logger.Error("error when loading table",
				zap.Error(err))
		}
	case in == cmdTableList:
		fmt.Println(strings.Join(a.TableList(), "\n"))
	case in == cmdHelp:
		fmt.Println("Available command description:")
		for _, helpItem := range helpList {
			fmt.Printf("\t%s\t - %s\n", helpItem.cmd, helpItem.desc)
		}
	case strings.HasPrefix(in, cmdDroupTable):
		newTableName := strings.TrimSpace(strings.TrimPrefix(in, `\drop`))
		if newTableName == "" {
			fmt.Printf("wrong syntax for %s: '%s'\n", cmdDroupTable, in)
		}
		if err := a.DropTable(newTableName); err != nil {
			fmt.Printf("error when dropping table: '%v'\n", err)
			logger.Error("error when dropping table",
				zap.Error(err))
		}
	default:
		start := time.Now()
		res, err := a.Execute(ctx, in)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			logger.Error("error executing query",
				zap.String("query", in),
				zap.Error(err))
			return
		}
		duration := time.Since(start)
		output, err := f.Format(ctx, res)
		if err != nil {
			logger.Error("error when formatting results",
				zap.String("query", in),
				zap.Error(err),
			)
			return
		}
		fmt.Println(output)
		logger.Info("query has been executed",
			zap.String("query", in),
			zap.String("table", res.Name),
			zap.Duration("duration", duration))
	}
}

func main() {
	log := getLogger()

	log.Info("starting csv-db")
	a := app.NewApp(log)
	f := &formatter.DefaultFormatter{}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	reader := readStdOut(ctx)

	log.Info("csv-db has been ready to accept queries")
	for {
		fmt.Print(welcomeQuery)
		select {
		case <-ctx.Done():
			fmt.Println("Bye-bye!")
			log.Info("staring gracefull shutdown")
			return
		case in := <-reader:
			if in == "" {
				continue
			}
			handleInput(ctx, log, a, f, in)
		}
	}
}

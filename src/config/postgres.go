package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
	"time"
)

type PGXOpenTracingLogger struct {
	tracer opentracing.Tracer
}

func formatQueryToMessage(query string) string {
	var message string

	if len(query) > 80 {
		message = query[:80]
	} else {
		message = query
	}

	replacer := strings.NewReplacer("\n", " ", "\t", " ", "   ", " ", "  ", " ")
	return replacer.Replace(message)
}

func (logger *PGXOpenTracingLogger) log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	possibleQuery := data["sql"]
	possibleRowCount := data["rowCount"]
	possibleDuration := data["time"]
	possibleArgs := data["args"]

	query, queryOk := possibleQuery.(string)
	rowCount, rowCountOk := possibleRowCount.(int)
	duration, durationOk := possibleDuration.(time.Duration)
	args, argsOk := possibleArgs.([]interface{})

	if !queryOk || !rowCountOk || !durationOk || !argsOk {
		return
	}

	message := formatQueryToMessage(query)
	startTime := opentracing.StartTime(time.Now().Add(duration * -1))
	span, ctx := opentracing.StartSpanFromContext(ctx, message, startTime)

	constantFieldsCount := 3

	fieldsCount := len(args) + constantFieldsCount

	logs := make([]log.Field, fieldsCount)

	logs[0] = log.String("sql", query)
	logs[1] = log.Int("rows", rowCount)
	logs[2] = log.Int64("duration", duration.Microseconds())

	i := constantFieldsCount

	for k, v := range args {

		key := fmt.Sprintf("%%%d", k+1)
		switch value := v.(type) {
		case string:
			logs[i] = log.String(key, value)
			i++
		case int:
			logs[i] = log.Int(key, value)
			i++
		case int32:
			logs[i] = log.Int32(key, value)
			i++
		case int64:
			logs[i] = log.Int64(key, value)
			i++
		case float32:
			logs[i] = log.Float32(key, value)
			i++
		case float64:
			logs[i] = log.Float64(key, value)
			i++
		}
	}

	span.LogFields(logs...)
	defer span.Finish()
}

func (logger *PGXOpenTracingLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	go logger.log(ctx, level, msg, data)
}

func InitPool(tracer opentracing.Tracer, environment *Environment) *pgxpool.Pool {

	databaseURI := environment.DatabaseURI

	if isTest() {
		databaseURI = databaseURI + "_test"
	}

	config, err := pgxpool.ParseConfig(databaseURI)

	if err != nil {
		panic(err)
	}

	if tracer != nil {
		config.ConnConfig.Logger = &PGXOpenTracingLogger{tracer: tracer}
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)

	if err != nil {
		panic(err)
	}

	return pool
}

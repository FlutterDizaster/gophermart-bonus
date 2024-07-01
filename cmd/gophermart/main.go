// This file is part of the gophermart-bonus project
//
// © 2024 Dmitriy Loginov
//
// Licensed under the MIT License. See the LICENSE.md file in the project root for more information.
//
// https://github.com/FlutterDizaster/gophermart-bonus
package main

import (
	"bufio"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	flag "github.com/spf13/pflag"

	server "github.com/FlutterDizaster/gophermart-bonus/internal"
	keygen "github.com/FlutterDizaster/gophermart-bonus/pkg/key-gen"
)

func main() {
	os.Exit(mainWithCode())
}

func mainWithCode() int {
	// Инициализация slog
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	// загрузка конфига
	settings := loadConfig()
	// Создание сервера
	srv, err := server.New(settings)
	if err != nil {
		slog.Error("Creating server error", slog.String("error", err.Error()))
		return 1
	}
	// Создание контекста отмены
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	// Запуск сервера
	if err = srv.Start(ctx); err != nil {
		slog.Error("Server startup error", slog.String("error", err.Error()))
		return 1
	}

	return 0
}

func loadConfig() server.Settings {
	const (
		defaultAddr    = "localhost:8080"
		defaultDBConn  = ""
		defaultAccrual = ""
	)

	var settings server.Settings

	flag.StringVarP(
		&settings.Addr,
		"addr",
		"a",
		defaultAddr,
		"Server endpoint addres. Default localhost:8080",
	)

	flag.StringVarP(
		&settings.DBConn,
		"dbstring",
		"d",
		defaultDBConn,
		"DB connection string",
	)

	flag.StringVarP(
		&settings.AccrualAddr,
		"accrual",
		"r",
		defaultAccrual,
		"Accrual service endpoint",
	)
	flag.Parse()

	return lookupEnvs(settings)
}

func lookupEnvs(settings server.Settings) server.Settings {
	envAddr, ok := os.LookupEnv("RUN_ADDRESS")
	if ok {
		settings.Addr = envAddr
	}

	envDBConn, ok := os.LookupEnv("DATABASE_URI")
	if ok {
		settings.DBConn = envDBConn
	}

	envAccrual, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	if ok {
		settings.AccrualAddr = envAccrual
	}

	return parseFiles(settings)
}

func parseFiles(settings server.Settings) server.Settings {
	const (
		defaultJWTKeyFile = "./jwt.key"
		defaultSHAKeyFile = "./sha.key"
	)

	settings.JWTSecret = loadKeyFromFile(defaultJWTKeyFile)
	settings.SHASecret = loadKeyFromFile(defaultSHAKeyFile)

	return settings
}

func loadKeyFromFile(path string) string {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		slog.Error("error reading key from file. Creating new key")
		return generateKeyFile(file)
	}
	key := scanner.Text()

	return key
}

func generateKeyFile(file *os.File) string {
	key := keygen.GenerateRandomKey(512)

	writer := bufio.NewWriter(file)
	_, err := writer.WriteString(key)
	if err != nil {
		slog.Error("error writing file")
		os.Exit(1)
	}

	err = writer.Flush()
	if err != nil {
		slog.Error("error flushing file")
		os.Exit(1)
	}

	file.Close()

	return key
}

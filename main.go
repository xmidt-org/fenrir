/**
 * Copyright 2019 Comcast Cable Communications Management, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"fmt"
	"github.com/Comcast/webpa-common/semaphore"
	_ "net/http/pprof"

	"github.com/Comcast/webpa-common/concurrent"

	"github.com/Comcast/codex/db"
	"github.com/Comcast/webpa-common/logging"
	"github.com/goph/emperror"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	//	"github.com/Comcast/webpa-common/secure/handler"
	"os"
	"os/signal"
	"time"

	"github.com/Comcast/webpa-common/server"
)

const (
	applicationName, apiBase = "fenrir", "/api/v1"
	DEFAULT_KEY_ID           = "current"
	applicationVersion       = "0.2.0"
)

type FenrirConfig struct {
	PruneInterval   time.Duration
	PruneRetries    int
	RetryInterval   time.Duration
	MaxPruneWorkers int
	Db              db.Config
}

func fenrir(arguments []string) int {
	start := time.Now()

	var (
		f, v                                = pflag.NewFlagSet(applicationName, pflag.ContinueOnError), viper.New()
		logger, metricsRegistry, codex, err = server.Initialize(applicationName, arguments, f, v, db.Metrics)
	)

	printVer := f.BoolP("version", "v", false, "displays the version number")
	if err := f.Parse(arguments); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse arguments: %s\n", err.Error())
		return 1
	}

	if *printVer {
		fmt.Println(applicationVersion)
		return 0
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to initialize viper: %s\n", err.Error())
		return 1
	}
	logging.Info(logger).Log(logging.MessageKey(), "Successfully loaded config file", "configurationFile", v.ConfigFileUsed())

	/*validator, err := server.GetValidator(v, DEFAULT_KEY_ID)
	 if err != nil {
		 fmt.Fprintf(os.Stderr, "Validator error: %v\n", err)
		 return 1
	 }*/

	config := new(FenrirConfig)
	v.Unmarshal(config)

	dbConn, err := db.CreateDbConnection(config.Db, metricsRegistry, nil)
	if err != nil {
		logging.Error(logger, emperror.Context(err)...).Log(logging.MessageKey(), "Failed to initialize database connection",
			logging.ErrorKey(), err.Error())
		fmt.Fprintf(os.Stderr, "Database Initialize Failed: %#v\n", err)
		return 2
	}

	updater := db.CreateRetryUpdateService(dbConn, config.PruneRetries, config.RetryInterval, metricsRegistry)

	if config.MaxPruneWorkers <= 0 {
		if config.Db.MaxOpenConns > 0 {
			logging.Warn(logger).Log(logging.MessageKey(), "invalid prune worker pool value defaulting to max open connections")
			config.MaxPruneWorkers = config.Db.MaxOpenConns
		} else {
			logging.Warn(logger).Log(logging.MessageKey(), "invalid prune worker pool value defaulting to 5")
			config.MaxPruneWorkers = 5
		}
	}

	pruner := pruner{
		updater:      updater,
		logger:       logger,
		pruneWorkers: semaphore.New(config.MaxPruneWorkers),
	}

	stopPruning := make(chan struct{}, 1)
	if config.PruneInterval > 0 {
		pruner.wg.Add(1)
		go pruner.handlePruning(stopPruning, config.PruneInterval)
	}
	// MARK: Starting the server
	_, runnable, done := codex.Prepare(logger, nil, metricsRegistry, nil)

	waitGroup, shutdown, err := concurrent.Execute(runnable)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start device manager: %s\n", err)
		return 1
	}

	logging.Info(logger).Log(logging.MessageKey(), fmt.Sprintf("%s is up and running!", applicationName), "elapsedTime", time.Since(start))
	signals := make(chan os.Signal, 10)
	signal.Notify(signals)
	for exit := false; !exit; {
		select {
		case s := <-signals:
			if s != os.Kill && s != os.Interrupt {
				logging.Info(logger).Log(logging.MessageKey(), "ignoring signal", "signal", s)
			} else {
				logging.Error(logger).Log(logging.MessageKey(), "exiting due to signal", "signal", s)
				exit = true
			}
		case <-done:
			exit = true
		}
	}

	stopPruning <- struct{}{}
	close(shutdown)
	pruner.wg.Wait()
	waitGroup.Wait()
	err = dbConn.Close()
	if err != nil {
		logging.Error(logger, emperror.Context(err)...).Log(logging.MessageKey(), "closing database threads failed",
			logging.ErrorKey(), err.Error())
	}
	return 0
}

func main() {
	os.Exit(fenrir(os.Args))
}

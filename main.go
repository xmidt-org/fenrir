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
	_ "net/http/pprof"

	"github.com/Comcast/codex/db/retry"

	"github.com/Comcast/codex/db/postgresql"

	"github.com/Comcast/codex/db/batchDeleter"

	"github.com/Comcast/webpa-common/concurrent"

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
	applicationVersion       = "0.5.1"
)

type FenrirConfig struct {
	PruneRetries RetryConfig
	Pruner       batchDeleter.Config
	Shards       []int
	Db           postgresql.Config
}

type RetryConfig struct {
	NumRetries   int
	Interval     time.Duration
	IntervalMult time.Duration
}

func fenrir(arguments []string) int {
	start := time.Now()

	var (
		f, v                                = pflag.NewFlagSet(applicationName, pflag.ContinueOnError), viper.New()
		logger, metricsRegistry, codex, err = server.Initialize(applicationName, arguments, f, v, postgresql.Metrics, dbretry.Metrics)
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

	config := new(FenrirConfig)
	v.Unmarshal(config)

	if len(config.Shards) == 0 {
		logging.Warn(logger).Log(logging.MessageKey(), "no shards given, defaulting to single 0 shard")
		config.Shards = []int{0}
	}

	dbConn, err := postgresql.CreateDbConnection(config.Db, metricsRegistry, nil)
	if err != nil {
		logging.Error(logger, emperror.Context(err)...).Log(logging.MessageKey(), "Failed to initialize database connection",
			logging.ErrorKey(), err.Error())
		fmt.Fprintf(os.Stderr, "Database Initialize Failed: %#v\n", err)
		return 2
	}

	updater := dbretry.CreateRetryUpdateService(
		dbConn,
		dbretry.WithRetries(config.PruneRetries.NumRetries),
		dbretry.WithInterval(config.PruneRetries.Interval),
		dbretry.WithIntervalMultiplier(config.PruneRetries.IntervalMult),
		dbretry.WithMeasures(metricsRegistry),
	)

	stopFuncs := make([]func(), len(config.Shards))

	for _, shard := range config.Shards {
		config.Pruner.Shard = shard
		deleter, err := batchDeleter.NewBatchDeleter(config.Pruner, logger, metricsRegistry, updater)
		if err != nil {
			logging.Error(logger, emperror.Context(err)...).Log(logging.MessageKey(), "Failed to initialize batch deleter",
				logging.ErrorKey(), err.Error(), "shard", shard)
			fmt.Fprintf(os.Stderr, "New Batch Deleter Failed: %#v\n", err)
			return 2
		}
		deleter.Start()
		stopFuncs = append(stopFuncs, deleter.Stop)
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

	close(shutdown)
	waitGroup.Wait()
	for _, stop := range stopFuncs {
		stop()
	}
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

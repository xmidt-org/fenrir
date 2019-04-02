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
	"sync"
	"time"

	"github.com/Comcast/codex/db"
	"github.com/Comcast/webpa-common/logging"
	"github.com/go-kit/kit/log"
	"github.com/goph/emperror"
)

type pruner struct {
	updater db.RetryUpdateService
	logger  log.Logger
	wg      sync.WaitGroup
}

func (r *pruner) handlePruning(quit chan struct{}, interval time.Duration) {
	defer r.wg.Done()
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-quit:
			return
		case <-t.C:
			r.pruneDevice()
		}
	}
}

func (r *pruner) pruneDevice() {
	err := r.updater.PruneRecords(time.Now().Unix())
	if err != nil {
		logging.Error(r.logger, emperror.Context(err)...).Log(logging.MessageKey(),
			"Failed to update event history", logging.ErrorKey(), err.Error())
		return
	}
	logging.Debug(r.logger).Log(logging.MessageKey(), "Successfully pruned events")
	return
}

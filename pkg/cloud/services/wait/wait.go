/*
Copyright 2020 FUJITSU CLOUD TECHNOLOGIES LIMITED. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package wait

import (
	"time"

	"github.com/pkg/errors"
	nferrors "github.com/nifcloud-labs/cluster-api-provider-nifcloud/pkg/cloud/errors"
	"k8s.io/apimachinery/pkg/util/wait"
)

func NewBackoff() wait.Backoff {
	// https://godoc.org/k8s.io/apimachinery/pkg/util/wait#Backoff
	return wait.Backoff{
		Duration: time.Second,
		Factor:   1.71,
		Jitter:   0.5,
		Steps:    10,
	}
}

func WaitForWithRetryable(backoff wait.Backoff, condition wait.ConditionFunc, retryableErrors ...string) error {
	var retErr error
	waitErr := wait.ExponentialBackoff(backoff, func() (bool, error) {
		// clear error from previous iteration
		retErr = nil
		ok, err := condition()
		if ok {
			// finish
			return true, nil
		}
		if err == nil {
			// not done, but no error, then skip waiting
			return false, nil
		}

		// check the error is retryable
		code, ok := nferrors.Code(errors.Cause(err))
		if !ok {
			return false, err
		}

		for _, r := range retryableErrors {
			if code == r {
				// retry
				retErr = err
				return false, nil
			}
		}

		// cannot retry
		return false, err
	})

	// is not timeout error
	if waitErr != wait.ErrWaitTimeout {
		return waitErr
	}

	// retryable error occured, return actual error
	if retErr != nil {
		return retErr
	}

	// timeout
	return waitErr
}

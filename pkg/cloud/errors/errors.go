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

package errors

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

const (
	AuthFailure             = "AuthFailure"
	GroupNotFound           = "InvalidGroup.NotFound"
	PermissionNotFound      = "InvalidPermission.NotFound"
	ResourceNotFound        = "InvalidResourceID.NotFound"
	InvalidInstanceID       = "InvalidInstanceID.NotFound"
	InvalidParameter        = "Client.InvalidParameterNotFound.Instance"
	SecurityGroupProcessing = "Server.ResourceIncorrectState.SecurityGroup.Processing"
)

var _error = &ServerError{}

func Code(err error) (string, bool) {
	if awserr, ok := err.(awserr.Error); ok {
		return awserr.Code(), true
	}
	return "", false
}

func Message(err error) string {
	if awserr, ok := err.(awserr.Error); ok {
		return awserr.Message()
	}
	return ""
}

type ServerError struct {
	err  error
	Code int
}

// Error implements the Error interface.
func (e *ServerError) Error() string {
	return e.err.Error()
}

// NewNotFound returns a new error which indicates that the resource of the kind and the name was not found.
func NewNotFound(err error) error {
	return &ServerError{
		err:  err,
		Code: http.StatusNotFound,
	}
}

// NewConflict returns a new error which indicates that the request cannot be processed due to a conflict.
func NewConflict(err error) error {
	return &ServerError{
		err:  err,
		Code: http.StatusConflict,
	}
}

// NewFailedDependency returns a new error which indicates that a dependency failure status
func NewFailedDependency(err error) error {
	return &ServerError{
		err:  err,
		Code: http.StatusFailedDependency,
	}
}

// IsFailedDependency checks if the error is pf http.StatusFailedDependency
func IsFailedDependency(err error) bool {
	if ReasonForError(err) == http.StatusFailedDependency {
		return true
	}
	return false
}

// IsNotFound returns true if the error was created by NewNotFound.
func IsNotFound(err error) bool {
	if ReasonForError(err) == http.StatusNotFound {
		return true
	}
	return IsInvalidNotFoundError(err)
}

// IsConflict returns true if the error was created by NewConflict.
func IsConflict(err error) bool {
	return ReasonForError(err) == http.StatusConflict
}

// IsSDKError returns true if the error is of type awserr.Error.
func IsSDKError(err error) (ok bool) {
	_, ok = err.(awserr.Error)
	return
}

// IsInvalidNotFoundError tests for common aws not found errors
func IsInvalidNotFoundError(err error) bool {
	if code, ok := Code(err); ok {
		switch code {
		case InvalidInstanceID:
			return true
		case InvalidParameter:
			return true
		}
	}

	return false
}

// ReasonForError returns the HTTP status for a particular error.
func ReasonForError(err error) int {
	switch t := err.(type) {
	case *ServerError:
		return t.Code
	}
	return -1
}

// IsIgnorableSecurityGroupError checks for errors in SG that can be ignored and then return nil.
func IsIgnorableSecurityGroupError(err error) error {
	if code, ok := Code(err); ok {
		switch code {
		case GroupNotFound, PermissionNotFound:
			return nil
		default:
			return err
		}
	}
	return nil
}

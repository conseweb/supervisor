/*
Copyright Mojing Inc. 2016 All Rights Reserved.

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
package protos

import "fmt"

func (this *Error) Error() string {
	return fmt.Sprintf("Error[%v]: %s", this.ErrorType, this.Message)
}

func (this *Error) OK() bool {
	return this.ErrorType == ErrorType_NONE_ERROR
}

func ResponseOK() *Error {
	return &Error{
		ErrorType: ErrorType_NONE_ERROR,
		Message:   "every thing is ok",
	}
}

func NewError(errorType ErrorType, msg string) *Error {
	return &Error{
		ErrorType: errorType,
		Message:   msg,
	}
}

func NewErrorf(errorType ErrorType, format string, args ...interface{}) *Error {
	return &Error{
		ErrorType: errorType,
		Message:   fmt.Sprintf(format, args...),
	}
}
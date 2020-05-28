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

package v1alpha2

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseTags(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want Tag
	}{
		{
			name: "empty input",
			in:   "",
			want: Tag{},
		},
		{
			name: "single tag item",
			in:   "hoge:fuga",
			want: Tag{
				"hoge": "fuga",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseTags(tt.in)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("got[%+v], want[%+v]", got, tt.want)
			}
		})
	}
}

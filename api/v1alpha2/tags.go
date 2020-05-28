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
	"fmt"
	"strings"

	"github.com/aokumasan/nifcloud-sdk-go-v2/nifcloud"
)

const (
	tagSeparator   = ","
	tagSeparatorKV = ":"
)

type BuildParams struct {
	ClusterName string
	// +opitonal
	Name *string
	// +optional
	Role *string
}

type Tag map[string]string

// parse description string to tag
func ParseTags(s string) Tag {
	tags := make(Tag)
	if len(s) == 0 {
		return tags
	}
	ss := strings.Split(s, tagSeparator)
	for _, v := range ss {
		kv := strings.Split(v, tagSeparatorKV)
		tags[kv[0]] = kv[1]
	}
	return tags
}

func BuildTags(params BuildParams) Tag {
	tags := make(Tag)
	tags["cluster"] = params.ClusterName
	if params.Role != nil {
		tags["role"] = *params.Role
	}
	if params.Name != nil {
		tags["Name"] = *params.Name
	}

	return tags
}

// convert tags to string which is set to description
func (t Tag) ConvToString() *string {
	var strTag []string

	for k, v := range t {
		strTag = append(strTag, fmt.Sprintf("%s%s%s", k, tagSeparatorKV, v))
	}
	return nifcloud.String(strings.Join(strTag, tagSeparator))
}

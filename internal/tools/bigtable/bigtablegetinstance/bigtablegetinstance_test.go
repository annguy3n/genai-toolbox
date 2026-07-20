// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bigtablegetinstance_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/mcp-toolbox/internal/server"
	"github.com/googleapis/mcp-toolbox/internal/testutils"
	"github.com/googleapis/mcp-toolbox/internal/tools"
	bigtablegetinstance "github.com/googleapis/mcp-toolbox/internal/tools/bigtable/bigtablegetinstance"
)

func TestParseFromYaml(t *testing.T) {
	ctx, err := testutils.ContextWithNewLogger()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	tcs := []struct {
		desc string
		in   string
		want server.ToolConfigs
	}{
		{
			desc: "basic example",
			in: `
            kind: tool
            name: mock-bigtablegetinstance
            type: bigtable-get-instance
            source: my-bigtable-source
            description: some description
            `,
			want: server.ToolConfigs{
				"mock-bigtablegetinstance": bigtablegetinstance.Config{
					ConfigBase: tools.ConfigBase{
						Name:         "mock-bigtablegetinstance",
						Description:  "some description",
						AuthRequired: []string{},
					},
					Type:   "bigtable-get-instance",
					Source: "my-bigtable-source",
				},
			},
		},
		{
			desc: "with auth required",
			in: `
            kind: tool
            name: mock-bigtablegetinstance-auth
            type: bigtable-get-instance
            source: my-bigtable-source
            description: some description
            authRequired: 
            - my-google-auth-service
            - other-auth-service
            `,
			want: server.ToolConfigs{
				"mock-bigtablegetinstance-auth": bigtablegetinstance.Config{
					ConfigBase: tools.ConfigBase{
						Name:         "mock-bigtablegetinstance-auth",
						Description:  "some description",
						AuthRequired: []string{"my-google-auth-service", "other-auth-service"},
					},
					Type:   "bigtable-get-instance",
					Source: "my-bigtable-source",
				},
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			// Parse contents
			_, _, _, got, _, _, err := server.UnmarshalPrimitiveConfig(ctx, testutils.FormatYaml(tc.in))
			if err != nil {
				t.Fatalf("unable to unmarshal: %s", err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("incorrect parse: diff %v", diff)
			}
		})
	}
}

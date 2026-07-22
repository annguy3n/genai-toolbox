// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bigtablelistschemas_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/googleapis/mcp-toolbox/internal/server"
	"github.com/googleapis/mcp-toolbox/internal/sources"
	"github.com/googleapis/mcp-toolbox/internal/testutils"
	"github.com/googleapis/mcp-toolbox/internal/tools"
	"github.com/googleapis/mcp-toolbox/internal/tools/bigtable/bigtablelistschemas"
)

func TestParseFromYaml(t *testing.T) {
	ctx, err := testutils.ContextWithNewLogger()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	doc := `
kind: tool
name: bigtable-list-schemas-test
type: bigtable-list-schemas
description: my tool description
source: my-instance
`

	_, _, _, got, _, _, err := server.UnmarshalPrimitiveConfig(ctx, testutils.FormatYaml(doc))
	if err != nil {
		t.Fatalf("UnmarshalPrimitiveConfig failed: %v", err)
	}

	actual, ok := got["bigtable-list-schemas-test"].(bigtablelistschemas.Config)
	if !ok {
		t.Fatalf("expected Config, got %T", got["bigtable-list-schemas-test"])
	}

	expected := bigtablelistschemas.Config{
		ConfigBase: tools.ConfigBase{
			Name:         "bigtable-list-schemas-test",
			Description:  "my tool description",
			AuthRequired: []string{},
		},
		Type:   "bigtable-list-schemas",
		Source: "my-instance",
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("parsed config mismatch. got %v, want %v", actual, expected)
	}
}

func TestInitialize(t *testing.T) {
	config := bigtablelistschemas.Config{
		ConfigBase: tools.ConfigBase{Name: "test"},
		Type:       "bigtable-list-schemas",
		Source:     "my-instance",
	}
	tool, err := config.Initialize(context.TODO())
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	if tool == nil {
		t.Fatal("Initialize returned nil tool")
	}
}

func TestToolConfigType(t *testing.T) {
	config := bigtablelistschemas.Config{
		ConfigBase: tools.ConfigBase{Name: "test"},
		Type:       "bigtable-list-schemas",
		Source:     "my-instance",
	}
	if got := config.ToolConfigType(); got != "bigtable-list-schemas" {
		t.Errorf("ToolConfigType() = %v, want bigtable-list-schemas", got)
	}
}

func TestToConfig(t *testing.T) {
	config := bigtablelistschemas.Config{
		ConfigBase: tools.ConfigBase{
			Name:        "test",
			Description: "List all Bigtable schemas, including tables with column family definitions, logical views, and materialized views.",
		},
		Type:   "bigtable-list-schemas",
		Source: "my-instance",
	}
	tool, _ := config.Initialize(context.TODO())

	if got := tool.ToConfig(); !reflect.DeepEqual(got, config) {
		t.Errorf("ToConfig() = %v, want %v", got, config)
	}
}

type mockSourceProvider struct{}

func (m mockSourceProvider) GetSource(name string) (sources.Source, bool) {
	return nil, false
}
func (m mockSourceProvider) GetSources() map[string]sources.Source {
	return nil
}

func TestInvoke(t *testing.T) {
	config := bigtablelistschemas.Config{
		ConfigBase: tools.ConfigBase{Name: "test"},
		Type:       "bigtable-list-schemas",
		Source:     "my-instance",
	}
	tool, _ := config.Initialize(context.TODO())

	_, err := tool.Invoke(context.TODO(), mockSourceProvider{}, nil, "")
	if err == nil {
		t.Error("Invoke() expected error for bad source, got nil")
	}
}

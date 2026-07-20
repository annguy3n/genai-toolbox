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

package bigtablecreatecluster

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/bigtable"
	yaml "github.com/goccy/go-yaml"
	"github.com/googleapis/mcp-toolbox/internal/tools"
	"github.com/googleapis/mcp-toolbox/internal/util"
	"github.com/googleapis/mcp-toolbox/internal/util/parameters"
)

const resourceType string = "bigtable-create-cluster"

func init() {
	if !tools.Register(resourceType, newConfig) {
		panic(fmt.Sprintf("tool type %q already registered", resourceType))
	}
}

func newConfig(ctx context.Context, name string, decoder *yaml.Decoder) (tools.ToolConfig, error) {
	actual := Config{ConfigBase: tools.ConfigBase{Name: name}}
	if err := decoder.DecodeContext(ctx, &actual); err != nil {
		return nil, err
	}
	return actual, nil
}

type compatibleSource interface {
	BigtableInstanceAdminClient() *bigtable.InstanceAdminClient
	BigtableAdminClient() *bigtable.AdminClient
	ProjectID() string
	InstanceID() string
}

type Config struct {
	tools.ConfigBase `yaml:",inline"`
	Type             string                 `yaml:"type" validate:"required"`
	Source           string                 `yaml:"source" validate:"required"`
	Annotations      *tools.ToolAnnotations `yaml:"annotations,omitempty"`
}

var _ tools.ToolConfig = Config{}

func (cfg Config) ToolConfigType() string {
	return resourceType
}

func (cfg Config) Initialize(context.Context) (tools.Tool, error) {
	if cfg.Description == "" {
		cfg.Description = "Create a new Bigtable cluster in an instance."
	}

	allParameters := parameters.Parameters{
		parameters.NewStringParameter("instance_id", "The ID of the instance", parameters.WithStringRequired(true)),
		parameters.NewStringParameter("cluster_id", "The ID of the cluster", parameters.WithStringRequired(true)),
		parameters.NewStringParameter("zone", "The zone for the cluster (e.g. us-central1-b)", parameters.WithStringRequired(true)),
		parameters.NewIntParameter("num_nodes", "The number of nodes to allocate", parameters.WithIntRequired(true)),
	}

	return Tool{
		BaseTool: tools.NewBaseTool(
			cfg,
			tools.GetAnnotationsOrDefault(cfg.Annotations, tools.NewDestructiveAnnotations),
			tools.Manifest{Description: cfg.Description, Parameters: allParameters.Manifest(), AuthRequired: cfg.AuthRequired},
			allParameters,
		),
	}, nil
}

var _ tools.Tool = Tool{}

type Tool struct {
	tools.BaseTool[Config]
}

func (t Tool) ToConfig() tools.ToolConfig {
	return t.Cfg
}

func (t Tool) Invoke(ctx context.Context, resourceMgr tools.SourceProvider, params parameters.ParamValues, accessToken tools.AccessToken) (any, util.ToolboxError) {
	source, err := tools.GetCompatibleSource[compatibleSource](resourceMgr, t.Cfg.Source, t.Cfg.Name, t.Cfg.Type)
	if err != nil {
		return nil, util.NewClientServerError("source used is not compatible with the tool", http.StatusInternalServerError, err)
	}

	paramsMap := params.AsMap()
	_ = paramsMap

	client := source.BigtableInstanceAdminClient()
	conf := &bigtable.ClusterConfig{
		InstanceID: paramsMap["instance_id"].(string),
		ClusterID:  paramsMap["cluster_id"].(string),
		Zone:       paramsMap["zone"].(string),
		NumNodes:   int32(paramsMap["num_nodes"].(int)),
	}
	err = client.CreateCluster(ctx, conf)
	if err != nil {
		return nil, util.NewClientServerError("failed to create cluster", http.StatusInternalServerError, err)
	}
	return map[string]string{"status": "cluster created successfully"}, nil
}

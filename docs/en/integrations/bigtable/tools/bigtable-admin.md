---
title: "bigtable-admin tools"
type: docs
weight: 2
description: >
  A suite of administrative tools for creating, updating, getting, listing, 
  and deleting Bigtable infrastructure assets like Instances, Clusters, 
  Tables, and Logical Views.
---

## About

The Bigtable Administrative MCP tools empower autonomous agents and workflows to natively manage the lifecycle of Bigtable topologies. Instead of being restricted solely to data-plane query tools (`bigtable-sql`), agents can directly provision or inspect underlying infrastructure via the `cloud.google.com/go/bigtable` Admin SDKs.

### Supported Operations

The following suite of tools are provided for complete structural manipulation:

### Instances
- `bigtable-get-instance`
- `bigtable-create-instance`
- `bigtable-update-instance`
- `bigtable-delete-instance`

### Clusters
- `bigtable-get-cluster`
- `bigtable-list-clusters`
- `bigtable-create-cluster`
- `bigtable-update-cluster`
- `bigtable-delete-cluster`

### Tables
- `bigtable-get-table`
- `bigtable-create-table`
- `bigtable-update-table`
- `bigtable-delete-table`

### Logical Views
- `bigtable-get-logical-view`
- `bigtable-list-logical-views`
- `bigtable-create-logical-view`
- `bigtable-update-logical-view`
- `bigtable-delete-logical-view`

## Compatible Sources

{{< compatible-sources >}}

## Example

```yaml
kind: tool
name: list_all_logical_views
type: bigtable-list-logical-views
source: my-bigtable-instance
description: |
  Use this tool to discover all Logical Views currently provisioned on the instance.
parameters:
  - name: instance_id
    type: string
    description: The parent instance ID.
```

### Security & Scoping

Because these tools directly wrap the Bigtable underlying Instance and Admin Clients, they invoke highly destructive operations. Ensure your application's IAM credential has valid `roles/bigtable.admin` bindings on the exact target Google Cloud Project. It is heavily recommended to use these tools sparingly or apply cautious tool-execution-confirmations in agent-driven configurations when triggering `delete` or `update` topologies.

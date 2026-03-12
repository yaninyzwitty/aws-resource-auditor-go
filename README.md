# AWS Resource Auditor

A command-line tool for auditing AWS resources across multiple services to identify security vulnerabilities, compliance issues, and cost optimization opportunities.

## Features

- **EC2 Auditing**: Find unused instances, old unused AMIs, and unencrypted volumes
- **S3 Auditing**: Detect public buckets, unencrypted buckets, and disabled versioning
- **IAM Auditing**: Identify stale access keys, unused roles, and weak password policies
- **Lambda Auditing**: Find old functions and outdated runtimes
- **RDS Auditing**: Discover idle instances, unencrypted databases, and publicly accessible instances
- **Secrets Manager Auditing**: Find unrotated secrets

## Installation

```bash
go install github.com/yaninyzwitty/aws-resource-auditor-go@latest
```

Or build from source:

```bash
git clone https://github.com/yaninyzwitty/aws-resource-auditor-go.git
cd aws-resource-auditor-go
go build -o aws-resource-auditor .
```

## Configuration

A `config.yaml` file is required. By default, the tool looks for `config.yaml` in the current working directory. Use `--config` to specify a custom path.

Create a `config.yaml` file:

```yaml
aws:
  profile: default
  region: us-east-1
  all_regions: false
  role_arn: ""
  external_id: ""

output:
  format: table
  export: ""
  no_color: false
  quiet: false

filter:
  severity: low
  tags: {}

thresholds:
  older_than: 90d
  days: 30

services:
  ec2:
    unused_older_than: 90d
    old_amis_older_than: 90d
  ebs:
    old_snapshots_older_than: 90d
  iam:
    stale_keys_older_than: 90d
    unused_roles_older_than: 90d
  secrets:
    unrotated_older_than: 90d
  lambda: {}
  rds:
    idle_days: 30
```

## Usage

### Global Flags

| Flag            | Description                               |
| --------------- | ----------------------------------------- |
| `-r, --region`  | AWS region                                |
| `--all-regions` | Scan all AWS regions                      |
| `-p, --profile` | AWS named profile                         |
| `--role-arn`    | IAM role ARN to assume                    |
| `--config`      | Path to config file                       |
| `-o, --output`  | Output format: table, json, csv, markdown |
| `--older-than`  | Age threshold (e.g., 90d)                 |

### Commands

#### EC2

```bash
# Run all EC2 audits
aws-resource-auditor ec2

# Find unused instances
aws-resource-auditor ec2 --unused

# Find old AMIs
aws-resource-auditor ec2 --old-amis

# Find unencrypted volumes
aws-resource-auditor ec2 --unencrypted

# Scan all regions
aws-resource-auditor ec2 --all-regions
```

#### S3

```bash
# Run all S3 audits
aws-resource-auditor s3

# Find public buckets
aws-resource-auditor s3 --public

# Find unencrypted buckets
aws-resource-auditor s3 --unencrypted

# Check versioning status
aws-resource-auditor s3 --versioning

# Scan all regions
aws-resource-auditor s3 --all-regions
```

#### IAM

```bash
# Run all IAM audits
aws-resource-auditor iam

# Find stale access keys
aws-resource-auditor iam --stale-keys

# Find unused roles
aws-resource-auditor iam --unused-roles

# Check password policy
aws-resource-auditor iam --password-policy
```

#### Lambda

```bash
# Run all Lambda audits
aws-resource-auditor lambda

# Find old functions
aws-resource-auditor lambda --old-functions

# Find outdated runtimes
aws-resource-auditor lambda --outdated-runtime

# Scan all regions
aws-resource-auditor lambda --all-regions
```

#### RDS

```bash
# Run all RDS audits
aws-resource-auditor rds

# Find idle instances
aws-resource-auditor rds --idle

# Find unencrypted instances
aws-resource-auditor rds --unencrypted

# Find publicly accessible instances
aws-resource-auditor rds --public

# Scan all regions
aws-resource-auditor rds --all-regions
```

#### Secrets

```bash
# Find unrotated secrets
aws-resource-auditor secrets --unrotated
```

## Environment Variables

| Variable                | Description                          |
| ----------------------- | ------------------------------------ |
| `AWS_PROFILE`           | AWS named profile                    |
| `AWS_REGION`            | Default AWS region                   |
| `AWS_ACCESS_KEY_ID`     | AWS access key (standard SDK var)    |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key (standard SDK var)    |
| `AUDIT_ROLE_ARN`        | IAM role to assume                   |
| `AUDIT_OUTPUT`          | Default output format                |
| `AUDIT_SEVERITY`        | Minimum severity level               |
| `NO_COLOR`              | Disable ANSI color output            |
| `AUDIT_EXPORT`          | Write results to file                |
| `AUDIT_LOG_LEVEL`       | Log level (debug, info, warn, error) |

## Output

The tool outputs audit findings with resource identifiers and relevant details. When no issues are found, it prints "No issues found".

Example output:

```text
Checking region: us-east-1
Findings:
  Instance: i-0123456789abcdef0 (web-server) - State: stopped - Launched: 120d ago
  Volume: vol-0123456789abcdef0 - Size: 100 GB - State: available
  AMI: ami-0123456789abcdef0 (my-ami) - Created: 100d ago
```

## License

MIT

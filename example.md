# Compute

audit ec2 --unused --region us-east-1 # Stopped/idle instances
audit ec2 --unattached-eips # Elastic IPs costing money
audit ec2 --old-amis --older-than 180d # Stale AMIs + snapshots
audit ec2 --oversized --cpu-threshold 5 # Instances averaging <5% CPU

# Storage

audit s3 --public-buckets
audit s3 --empty-buckets
audit s3 --no-lifecycle # Buckets with no lifecycle rules
audit s3 --unencrypted
audit ebs --unattached # Volumes not attached to any instance
audit ebs --old-snapshots --older-than 90d

# Networking

audit vpc --empty # VPCs with no resources
audit sg --unused # Security groups attached to nothing
audit sg --wide-open # 0.0.0.0/0 ingress on sensitive ports
audit elb --no-targets # Load balancers with empty target groups
audit nat --utilization # NAT Gateways with low traffic

audit iam --stale-keys --older-than 90d
audit iam --root-usage # Root account last activity
audit iam --no-mfa # Users without MFA
audit iam --unused-roles --older-than 60d # Roles never assumed
audit iam --overpermissioned # Roles/users with \*, wildcard actions
audit iam --inline-policies # Inline vs managed policy hygiene

audit secrets --unused --older-than 30d # Secrets Manager entries never accessed
audit secrets --unrotated --older-than 90d

audit costs --top-services --last 30d # Pull from Cost Explorer
audit costs --idle-resources # Cross-service idle resource report
audit costs --reserved-coverage # On-demand instances that should be RIs
audit costs --savings-plans --gap # Uncovered compute by Savings Plans

audit lambda --never-invoked --older-than 30d
audit lambda --high-error-rate --threshold 10 # Functions with >10% error rate
audit rds --idle --days 14 # DBs with near-zero connections
audit rds --unencrypted
audit rds --no-backup # Automated backups disabled

audit cloudtrail --disabled-regions # Regions with no trail
audit cloudtrail --no-log-validation
audit config --disabled # AWS Config not enabled per region
audit kms --key-rotation-disabled
audit ecr --unscanned --older-than 7d # Images never vulnerability-scanned
audit tags --missing --required Name,Env,Owner # Untagged resources across services

--output table|json|csv|markdown
--region us-east-1 # Single region
--all-regions # Scan everything
--profile myprofile # AWS named profile
--severity critical|high|medium|low
--export ./report-2024.csv
--fix # Dry-run remediation suggestions
--watch --interval 6h # Daemon mode, re-scan on interval

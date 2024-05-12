# HongJot Diagram

```mermaid
sequenceDiagram
	Mobile App ->> HongJot API: Upload Image to S3 Bucket
	activate Mobile App
	HongJot API ->> S3 Bucket: Store Image
	HongJot API ->> Mobile App: success
	deactivate Mobile App
	S3 Bucket ->> Lambda Function: (Event Trigger)
	Lambda Function ->> Amazon Textract: Extract Text from Image
	Amazon Textract ->> Lambda Function: Slip Expense info
	Lambda Function ->> HongJot API: Send Slip Expense info
	HongJot API ->> PostgreSQL Database: Store Expense info
	Mobile App ->> HongJot API: Get Expense Summary (Spender Requests Summary) Retrieve Expense Summary
	activate Mobile App
	HongJot API ->> PostgreSQL Database: Get Expense Summary Data
	PostgreSQL Database ->> HongJot API: Return Summary Data
	HongJot API ->> Mobile App: Show Expense Summary Data
	deactivate Mobile App
```

## Infrastructure

We have created the infrastructure by using Terraform. The infrastructure consists of the following:

- `Kubernetes Cluster (EKS)` - Managed Kubernetes Service
- `PostgreSQL (RDS)` - Relational Database (Managed Service)
- `SonarQube` - Static Code Analysis
- Github (see `.github/workflows` directory)
	- Github Action
	- Github Container Registry
- ArgoCD - Continuous Delivery
- Cloudflare - DNS Management

# AWS Infrastructure Diagram

```mermaid
%%{init: {'theme': 'base', 'themeVariables': { 'primaryColor': '#1a365d', 'primaryTextColor': '#fff', 'primaryBorderColor': '#2c5282', 'lineColor': '#2c5282', 'secondaryColor': '#2d3748', 'tertiaryColor': '#4a5568' } }}%%
flowchart TB
    %% Styling
    classDef vpc fill:#2c5282,stroke:#2a4365,color:white,stroke-width:2px
    classDef fargate fill:#2f855a,stroke:#2c7a5b,color:white,stroke-width:2px
    classDef alb fill:#6b46c1,stroke:#553c9a,color:white,stroke-width:2px
    classDef db fill:#9c4221,stroke:#7b341e,color:white,stroke-width:2px
    classDef internet fill:#4a5568,stroke:#2d3748,color:white,stroke-width:2px
    classDef cw fill:#ff7e47,stroke:#e65100,color:white,stroke-width:2px

    %% Internet
    Internet(("fa:fa-globe Internet"))
    
    %% VPC and Networking
    subgraph VPC["AWS VPC\n10.0.0.0/16"]
        subgraph PublicSubnets["Public Subnets"]
            SubnetA["us-east-1a\n10.0.1.0/24"]
            SubnetB["us-east-1b\n10.0.2.0/24"]
        end
        
        IGW["fa:fa-exchange Internet Gateway"]
        
        subgraph SecurityGroups["Security Groups"]
            ALB_SG["ALB SG\n(Allows HTTP/80 from Internet)"]
            ECS_SG["ECS SG\n(Allows 80, 8080 from ALB SG)"]
        end
        
        subgraph ECSCluster["ECS Fargate Cluster"]
            subgraph APIService["API Service"]
                APIContainer["hillary-api:8080\nFargate Task"]
            end
            
            subgraph ClientService["Client Service"]
                ClientContainer["hillary-client:80\nFargate Task"]
            end
        end
        
        ALB["fa:fa-random Application Load Balancer"]
        APITG["API Target Group\nPort 8080"]
        ClientTG["Client Target Group\nPort 80"]
    end
    
    %% External Services
    DB["fa:fa-database CockroachDB\nExternal"]
    
    %% AWS Services
    IAM["fa:fa-key IAM Role\n(ECS Task Execution)"]
    CloudWatch["fa:fa-chart-line CloudWatch Logs"]
    
    %% Connections
    Internet --> |HTTP/80| ALB
    ALB --> |/api/*| APITG --> APIService
    ALB --> |/*| ClientTG --> ClientService
    
    APIService --> |Database| DB
    IAM -.- |Assumes| ECSCluster
    
    %% Logging
    APIContainer --> |Logs| CloudWatch
    ClientContainer --> |Logs| CloudWatch
    
    %% Styling
    class VPC,SubnetA,SubnetB,IGW vpc
    class ECSCluster,APIContainer,ClientContainer fargate
    class ALB,ALB_SG,APITG,ClientTG alb
    class DB db
    class Internet internet
    class CloudWatch cw
    
    %% Icons
    linkStyle default fill:none,stroke:#4a5568,stroke-width:2px
```

## Infrastructure Components

1. **Networking**
   - VPC with public subnets in 2 AZs (us-east-1a, us-east-1b)
   - Internet Gateway for public internet access
   - Route tables for public subnets

2. **Load Balancing**
   - Application Load Balancer (ALB) with:
     - HTTP (port 80) listener
     - Target groups for API (port 8080) and App (port 80) services
     - Path-based routing (/api/* to API service)

3. **Container Orchestration**
   - ECS Cluster
   - ECS Services:
     - API Service (port 8080)
     - App Service (port 80)
   - IAM roles for ECS task execution

4. **Security**
   - ALB Security Group (allows HTTP/80 from anywhere)
   - ECS Security Group (allows traffic from ALB on ports 8080/80)

5. **Database**
   - External CockroachDB instance (configured via variables)

## Access Patterns
- External users access the application via the ALB
- ALB routes /api/* requests to the API service
- All other requests are routed to the App service
- API service connects to the external CockroachDB database

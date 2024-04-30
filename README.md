# job-scraper
A timed event that once a day scraps relevant jobs links and sends them to discord.
![job-scraper](https://github.com/austin1237/job-scraper/assets/1394341/bf78fb86-e5a0-4399-98c2-cbfa24da8496)

## Deployment
Deployment currently uses [Terraform](https://www.terraform.io/) to set up AWS services.
### Prerequisites
This repo needs a private [Amazon ECR repo](https://us-east-1.console.aws.amazon.com/ecr/repositories?region=us-east-1) to be created in the same region that our container based lambda is deployed to (in our case us-east-1). Name the private repo to headless.

### Setting up remote state
Terraform has a feature called [remote state](https://www.terraform.io/docs/state/remote.html) which ensures the state of your infrastructure to be in sync for mutiple team members as well as any CI system.

This project **requires** this feature to be configured. To configure **USE THE FOLLOWING COMMAND ONCE PER TEAM**.

```bash
cd terraform/remote-state
terraform init
terraform apply
```
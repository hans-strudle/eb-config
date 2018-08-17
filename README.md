# eb-config

Simple command line tool written in Go that queries AWS ElasticBeanstalk for Configuration Settings on environments. By default it will print out all Settings across all of your Environments in all of your Applications. There are arguments to filter down to specific Applications & Environments, as well as filtering for specific Options.

## Usage
# Arguments
```
$ > ./eb-config -h
Usage ./eb-config:
  -a string
    	Application Name
  -e string
    	Environment Name (Must also specify an Application Name if using this flag)
  -p string
    	Comma seperated list of properties to match for (ex. ./eb-config -p conn,Deploy)
  -r string
    	AWS Region (default "us-west-2")
  -v string
    	Comma seperated list of values to match for (ex. ./eb-config -v micro,0.0.0.0)
```
# Examples
Print out all Configuration Settings in Application "my-app" that contain "InstanceType" (InstanseType, InstanceTypeFamily)
```
$ > ./eb-config -a my-app -p instancetype
my-app, env-prod, aws:autoscaling:launchconfiguration, InstanceType, t2.large
my-app, env-prod, aws:cloudformation:template:parameter, InstanceTypeFamily, t2
my-app, env-staging, aws:autoscaling:launchconfiguration, InstanceType, t2.large
my-app, env-staging, aws:cloudformation:template:parameter, InstanceTypeFamily, t2
my-app, env-dev, aws:autoscaling:launchconfiguration, InstanceType, t2.small
my-app, env-dev, aws:cloudformation:template:parameter, InstanceTypeFamily, t2
```
Find every environment that has RollingUpdatesEnabled set to false
```
my-app, env-dev, aws:autoscaling:updatepolicy:rollingupdate, RollingUpdateEnabled, false
my-app-2, env-dev, aws:autoscaling:updatepolicy:rollingupdate, RollingUpdateEnabled, false
my-app-2, env-prod, aws:autoscaling:updatepolicy:rollingupdate, RollingUpdateEnabled, false
my-app-3, env-dev, aws:autoscaling:updatepolicy:rollingupdate, RollingUpdateEnabled, false
```

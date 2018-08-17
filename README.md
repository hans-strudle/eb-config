# eb-config

Simple command line tool written in Go that queries AWS ElasticBeanstalk for Configuration Settings on environments. By default it will print out all Settings across all of your Environments in all of your Applications. There are arguments to filter down to specific Applications & Environments, as well as filtering for specific Options.

## Usage
```
$ > ./eb-config -h
Usage of ./eb-config:
  -a string
    	Application Name
  -e string
    	Environment Name (Must also specify an Application Name if using this flag)
  -p string
    	Comma seperated list of properties to match for (ex. ./eb-config -p conn,Deploy)
```

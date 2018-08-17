package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	props  = flag.String("p", "", "Comma seperated list of properties to match for (ex. ./eb-config -p conn,Deploy)")
	vals   = flag.String("v", "", "Comma seperated list of values to match for (ex. ./eb-config -v micro,0.0.0.0)")
	app    = flag.String("a", "", "Application Name")
	env    = flag.String("e", "", "Environment Name (Must also specify an Application Name if using this flag)")
	region = flag.String("r", "us-west-2", "AWS Region")
)

var logger = log.New(os.Stdout, "", 0)

func fuzzyStrListMatch(l []string, s string) (b bool) {
	// returns true if any string in l is a substring of s, case insensitive
	s = strings.ToLower(s)
	for _, i := range l {
		b = strings.Contains(s, strings.ToLower(i))
		if b == true {
			return true
		}
	}
	return
}

func displayConfig(settings []*elasticbeanstalk.ConfigurationSettingsDescription, props []string, vals []string) {
	var line string
	for _, config := range settings {
		for _, option := range config.OptionSettings {
			if len(props) < 1 || fuzzyStrListMatch(props, *option.OptionName) {
				line = fmt.Sprintf("%s, %s, %s, %s, ", *config.ApplicationName, *config.EnvironmentName, *option.Namespace, *option.OptionName)
				if option.Value != nil {
					if len(vals) > 0 && !fuzzyStrListMatch(vals, *option.Value) {
						continue
					}
					line = line + *option.Value
				} else if len(vals) > 0 {
					continue
				} else {
					line = line + " Nothing"
				}
				logger.Println(line)
			}
		}
	}
	return
}

func getConfigSettings(svc *elasticbeanstalk.ElasticBeanstalk, app, env string) *elasticbeanstalk.DescribeConfigurationSettingsOutput {
	out, err := svc.DescribeConfigurationSettings(&elasticbeanstalk.DescribeConfigurationSettingsInput{
		ApplicationName: aws.String(app),
		EnvironmentName: aws.String(env),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elasticbeanstalk.ErrCodeTooManyBucketsException:
				fmt.Println(elasticbeanstalk.ErrCodeTooManyBucketsException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}
	return out

}

func main() {
	flag.Parse()
	propsList := strings.Split(*props, ",")
	valsList := strings.Split(*vals, ",")
	sess := session.Must(session.NewSession())
	svc := elasticbeanstalk.New(sess, aws.NewConfig().WithRegion(*region))

	appInput := &elasticbeanstalk.DescribeEnvironmentsInput{}
	if *app != "" {
		appInput.ApplicationName = app
	}
	if *env != "" { // we can skip calling for the environments
		out := getConfigSettings(svc, *app, *env)
		displayConfig(out.ConfigurationSettings, propsList, valsList)
		os.Exit(0)
	}
	res, err := svc.DescribeEnvironments(appInput) // for getting Application Names
	if err != nil {
		fmt.Println(err)
	}
	var wg sync.WaitGroup
	wg.Add(len(res.Environments))
	for _, env := range res.Environments {
		go func(env *elasticbeanstalk.EnvironmentDescription) {
			defer wg.Done()
			out := getConfigSettings(svc, *env.ApplicationName, *env.EnvironmentName)
			displayConfig(out.ConfigurationSettings, propsList, valsList)
		}(env)
	}
	wg.Wait()
}

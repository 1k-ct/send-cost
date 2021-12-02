package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// response
// [{
// 	Datapoints: [{
// 		Maximum: 0.01,
// 		Timestamp: 2021-11-28 23:57:00 +0000 UTC,
// 		Unit: "None"
// 	}],
// 	Label: "EstimatedCharges"
// } {
// 	Datapoints: [{
// 		Maximum: 0,
// 		Timestamp: 2021-11-28 23:57:00 +0000 UTC,
// 		Unit: "None"
// 	}],
// 	Label: "EstimatedCharges"
// }
// }]
//
// 非推奨: Use client package's EndpointsID value instead of these ServiceIDs. These IDs are not maintained, and are out of date.
// 非推奨：Service identifiers
// aws services https://pkg.go.dev/github.com/aws/aws-sdk-go@v1.42.16/aws/endpoints#pkg-constants
//
// https://docs.aws.amazon.com/ja_jp/AmazonCloudWatch/latest/APIReference/API_GetMetricStatistics.html
func FetchMetricStatisticsBilling(serviceName string) (metricsStatistic *cloudwatch.GetMetricStatisticsOutput, err error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		return metricsStatistic, err
	}

	svc := cloudwatch.New(sess)

	input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/Billing"),
		MetricName: aws.String("EstimatedCharges"),
		Period:     aws.Int64(86400), // 1 day = 86400 sec
		StartTime:  aws.Time(time.Now().AddDate(0, 0, -1)),
		EndTime:    aws.Time(time.Now()),
		// StartTime: aws.Time(time.Now().AddDate(0, 0, -3)),
		// EndTime:   aws.Time(time.Now().AddDate(0, 0, -2)),
		Statistics: []*string{
			aws.String(cloudwatch.StatisticMaximum),
		},
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("Currency"),
				Value: aws.String("USD"),
			}, {
				Name:  aws.String("ServiceName"),
				Value: aws.String(serviceName),
			},
		},
		Unit: aws.String(cloudwatch.StandardUnitNone),
	}

	metricsStatistic, err = svc.GetMetricStatistics(input)
	if err != nil {
		return metricsStatistic, err
	}
	return metricsStatistic, nil
}
func FetchTotalBilling() (metricsStatistic *cloudwatch.GetMetricStatisticsOutput, err error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		return metricsStatistic, err
	}

	svc := cloudwatch.New(sess)

	input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/Billing"),
		MetricName: aws.String("EstimatedCharges"),
		Period:     aws.Int64(86400), // 1 day = 86400 sec
		StartTime:  aws.Time(time.Now().AddDate(0, 0, -1)),
		EndTime:    aws.Time(time.Now()),
		// StartTime: aws.Time(time.Now().AddDate(0, 0, -3)),
		// EndTime:   aws.Time(time.Now().AddDate(0, 0, -2)),
		Statistics: []*string{
			aws.String(cloudwatch.StatisticMaximum),
		},
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("Currency"),
				Value: aws.String("USD"),
			},
		},
	}

	metricsStatistic, err = svc.GetMetricStatistics(input)
	if err != nil {
		return metricsStatistic, err
	}
	return metricsStatistic, nil
}

// response
// {
// 	Metrics: [{
// 		Dimensions: [{
// 			Name: "ServiceName",
// 			Value: "AmazonCloudWatch"
// 		},{
// 			Name: "Currency",
// 			Value: "USD"
// 		}],
// 		MetricName: "EstimatedCharges",
// 		Namespace: "AWS/Billing"
// 	},{
// 		Dimensions: [{
// 			Name: "Currency",
// 			Value: "USD"
// 		}],
// 		MetricName: "EstimatedCharges",
// 		Namespace: "AWS/Billing"
// 	}]
// }
// awsの使っているサービスを取得する
// [AmazonCloudWatch AWSSecretsManager AmazonRoute53 AmazonS3 AWSCloudTrail AWSLambda AmazonRDS AWSMarketplace AWSELB AmazonEC2 AWSDataTransfer awskms]
//
// https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/go/example_code/cloudwatch/listing_metrics.go
func FetchServices() (services []string, err error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		return nil, fmt.Errorf("session create error : %s", err)
	}
	svc := cloudwatch.New(sess)

	listMetricsOutput, err := svc.ListMetrics(&cloudwatch.ListMetricsInput{
		Dimensions: []*cloudwatch.DimensionFilter{},
		MetricName: aws.String("EstimatedCharges"),
		Namespace:  aws.String("AWS/Billing"),
	})
	if err != nil {
		return nil, fmt.Errorf("get metrics error : %s", err)
	}

	for _, metric := range listMetricsOutput.Metrics {
		name := *metric.Dimensions[0].Name
		value := *metric.Dimensions[0].Value
		if name == "ServiceName" {
			services = append(services, value)
		}
	}
	return services, nil
}

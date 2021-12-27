package main

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"golang.org/x/xerrors"
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
	session := &Sessioner{s: &CloudwatchSession{}}
	metric := &GettedMetricStatisticsOne{ServiceName: serviceName}
	metricsStatistic, err = session.fetchMetricStatistics(metric)
	if err != nil {
		return metricsStatistic, err
	}
	return metricsStatistic, nil
}
func FetchTotalBilling() (metricsStatistic *cloudwatch.GetMetricStatisticsOutput, err error) {
	session := &Sessioner{s: &CloudwatchSession{}}
	metric := &GettedMetricStatisticsAll{}
	metricsStatistic, err = session.fetchMetricStatistics(metric)
	if err != nil {
		return metricsStatistic, err
	}
	return metricsStatistic, nil
}
func (s *Sessioner) fetchMetricStatistics(fetcher Fetcher) (metricsStatistic *cloudwatch.GetMetricStatisticsOutput, err error) {
	svc, err := s.s.newCloudwatchSession()
	if err != nil {
		return metricsStatistic, err
	}

	f := &Fetched{f: fetcher}
	metricsStatistic, err = f.f.fetch(svc)
	if err != nil {
		return metricsStatistic, err
	}
	return metricsStatistic, nil
}

type Sessioner struct {
	s sessioner
}
type sessioner interface {
	newCloudwatchSession() (*cloudwatch.CloudWatch, error)
}
type CloudwatchSession struct{}

func (c *CloudwatchSession) newCloudwatchSession() (*cloudwatch.CloudWatch, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		return nil, xerrors.New(err.Error())
	}
	svc := cloudwatch.New(sess)
	return svc, nil
}

type Fetcher interface {
	fetch(*cloudwatch.CloudWatch) (*cloudwatch.GetMetricStatisticsOutput, error)
}
type Fetched struct {
	f Fetcher
}
type GettedMetricStatisticsOne struct {
	// dimensions []*cloudwatch.Dimension
	ServiceName string
}

func (g *GettedMetricStatisticsOne) fetch(svc *cloudwatch.CloudWatch) (*cloudwatch.GetMetricStatisticsOutput, error) {
	dimensions := []*cloudwatch.Dimension{
		{Name: aws.String("Currency"), Value: aws.String("USD")},
		{Name: aws.String("ServiceName"), Value: aws.String(g.ServiceName)},
	}
	input := setMetricStatistics(dimensions)
	metricsStatistic, err := svc.GetMetricStatistics(input)
	if err != nil {
		return metricsStatistic, err
	}
	return metricsStatistic, nil
}

func setMetricStatistics(dimensions []*cloudwatch.Dimension) *cloudwatch.GetMetricStatisticsInput {
	input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/Billing"),
		MetricName: aws.String("EstimatedCharges"),
		Period:     aws.Int64(86400), // 1 day = 86400 sec
		StartTime:  aws.Time(time.Now().AddDate(0, 0, -1)),
		EndTime:    aws.Time(time.Now()),
		Statistics: []*string{
			aws.String(cloudwatch.StatisticMaximum),
		},
		Dimensions: dimensions,
		Unit:       aws.String(cloudwatch.StandardUnitNone),
	}
	return input
}

type GettedMetricStatisticsAll struct{}

func (g *GettedMetricStatisticsAll) fetch(svc *cloudwatch.CloudWatch) (*cloudwatch.GetMetricStatisticsOutput, error) {
	dimensions := []*cloudwatch.Dimension{
		{Name: aws.String("Currency"), Value: aws.String("USD")},
	}
	input := setMetricStatistics(dimensions)
	metricsStatistic, err := svc.GetMetricStatistics(input)
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
func FetchMetricStatisticServices() (services []string, err error) {
	session := &Sessioner{s: &CloudwatchSession{}}
	listMetric := &ListedMetrics{}
	services, err = session.fetchMetricStatisticServices(listMetric)
	if err != nil {
		return services, err
	}
	return services, nil
}
func (s *Sessioner) fetchMetricStatisticServices(fetcher FetcherService) (services []string, err error) {
	svc, err := s.s.newCloudwatchSession()
	f := &FetchedService{f: fetcher}
	listMetricsOutput, err := f.f.fetch(svc)
	if err != nil {
		return nil, xerrors.New(err.Error())
	}
	if listMetricsOutput == nil {
		return services, xerrors.New("listMetricsOutput is nil")
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

type FetcherService interface {
	fetch(*cloudwatch.CloudWatch) (*cloudwatch.ListMetricsOutput, error)
}
type FetchedService struct {
	f FetcherService
}
type ListedMetrics struct{}

func (l *ListedMetrics) fetch(svc *cloudwatch.CloudWatch) (*cloudwatch.ListMetricsOutput, error) {
	listMetricsOutput, err := svc.ListMetrics(&cloudwatch.ListMetricsInput{
		Dimensions: []*cloudwatch.DimensionFilter{},
		MetricName: aws.String("EstimatedCharges"),
		Namespace:  aws.String("AWS/Billing"),
	})
	if err != nil {
		return listMetricsOutput, err
	}
	return listMetricsOutput, nil
}

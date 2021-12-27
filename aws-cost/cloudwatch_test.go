package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"golang.org/x/xerrors"
)

func TestFetchMetricStatisticsBilling(t *testing.T) {
	// serviceName := "AmazonRoute53"
	session := &Sessioner{s: &CloudwatchSession{}}
	metricsStatistic, err := session.fetchMetricStatistics(&GettedMetricStatisticsOneMock{})
	if err != nil {
		t.Fatalf("%+v\n", err)
	}
	t.Log(metricsStatistic)
}

type GettedMetricStatisticsOneMock struct{}

func (m *GettedMetricStatisticsOneMock) fetch(svc *cloudwatch.CloudWatch) (*cloudwatch.GetMetricStatisticsOutput, error) {
	return &cloudwatch.GetMetricStatisticsOutput{
		Datapoints: []*cloudwatch.Datapoint{{
			Maximum: aws.Float64(0.123),
			Unit:    aws.String("None"),
		}},
		Label: aws.String("EstimatedCharges"),
	}, nil
}
func TestFetchTotalBilling(t *testing.T) {
	session := &Sessioner{s: &CloudwatchSession{}}
	metricsStatistic, err := session.fetchMetricStatistics(&GettedMetricStatisticsAllMock{})
	if err != nil {
		t.Fatalf("%+v\n", err)
	}
	t.Log(metricsStatistic)
}

type GettedMetricStatisticsAllMock struct{}

func (m *GettedMetricStatisticsAllMock) fetch(svc *cloudwatch.CloudWatch) (*cloudwatch.GetMetricStatisticsOutput, error) {
	t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", "2021-11-28 23:57:00 +0000 UTC")
	return &cloudwatch.GetMetricStatisticsOutput{
		Datapoints: []*cloudwatch.Datapoint{{
			Maximum:   aws.Float64(1.55),
			Timestamp: aws.Time(t),
			Unit:      aws.String("None"),
		}},
		Label: aws.String("EstimatedCharges"),
	}, nil
}
func TestFetchMetricStatisticServices(t *testing.T) {
	session := &Sessioner{s: &CloudwatchSession{}}
	services, err := session.fetchMetricStatisticServices(&ListedMetricsMock{})
	if err != nil {
		t.Fatalf("%+v\n", err)
	}
	t.Log(services)
}

type ListedMetricsMock struct{}

func (m *ListedMetricsMock) fetch(svc *cloudwatch.CloudWatch) (*cloudwatch.ListMetricsOutput, error) {
	bfile, err := ioutil.ReadFile("./test/services.json")
	if err != nil {
		return nil, xerrors.New(err.Error())
	}
	ListMetricsOutput := &cloudwatch.ListMetricsOutput{}
	if err := json.Unmarshal(bfile, ListMetricsOutput); err != nil {
		return nil, xerrors.New(err.Error())
	}
	return ListMetricsOutput, nil
}

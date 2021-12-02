参考,引用
https://github.com/ookamiinc/aws-billing-notifier

https://aws.amazon.com/jp/premiumsupport/knowledge-center/cloudwatch-estimatedcharges-alarm/

```json
$ cat metric_data_queries.json
{
    "MetricDataQueries": [
        {
            "Id": "m1",
            "MetricStat": {
                "Metric": {
                    "Namespace": "AWS/Billing",
                    "MetricName": "EstimatedCharges",
                    "Dimensions": [
                        {
                            "Name": "Currency",
                            "Value": "USD"
                        }
                    ]
                },
                "Period": 86400,
                "Stat": "Maximum"
            }
        },
    {
            "Id": "e1",
            "Expression": "RATE(m1) * 86400",
            "Label": "DailyEstimatedCharges",
        "Period": 86400
        }
    ],
    "StartTime": "2020-06-01T00:00:00Z",
    "EndTime": "2020-06-05T00:00:00Z",
    "ScanBy": "TimestampAscending"
}


$ aws cloudwatch get-metric-data --cli-input-json file://metric_data_queries.json
{
    "MetricDataResults": [
        {
            "Id": "m1",
            "Label": "EstimatedCharges",
            "Timestamps": [
                "2020-06-01T00:00:00Z",
                "2020-06-02T00:00:00Z",
                "2020-06-03T00:00:00Z",
                "2020-06-04T00:00:00Z"
            ],
            "Values": [
                0.0,
                22.65,
                34.74,
                46.91
            ],
            "StatusCode": "Complete"
        },
        {
            "Id": "e1",
            "Label": "DailyEstimatedCharges",
            "Timestamps": [
                "2020-06-02T00:00:00Z",
                "2020-06-03T00:00:00Z",
                "2020-06-04T00:00:00Z"
            ],
            "Values": [
                22.65,
                12.090000000000003,
                12.169999999999995
            ],
            "StatusCode": "Complete"
        }
    ],
    "Messages": []
}
```

各サービスの一日の費用
ec2 000
lambda 000

total 000
今月 000

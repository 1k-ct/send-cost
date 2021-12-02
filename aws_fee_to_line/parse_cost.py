from dotenv import load_dotenv
import os
import datetime
# import dateutil
import time
# from dateutil.tz import tzlocal
import boto3
import requests
import logging


logger = logging.getLogger()
logger.setLevel(logging.INFO)

client = boto3.client('cloudwatch', region_name='us-east-1')

# 0.22sec


def parse_services():
    response = client.list_metrics(
        Namespace="AWS/Billing", MetricName="EstimatedCharges"
    )
    services = []

    for iRes in response["Metrics"]:
        value = iRes["Dimensions"][0]["Value"]
        if iRes["Dimensions"][0]["Name"] == "ServiceName":
            services.append(value)
    return services


def parse_total_service_billing():
    response = client.get_metric_statistics(
        Namespace='AWS/Billing',
        MetricName='EstimatedCharges',
        Dimensions=[
            {
                'Name': 'Currency',
                'Value': 'USD'
            }
        ],
        StartTime=datetime.datetime.today() - datetime.timedelta(days=1),
        EndTime=datetime.datetime.today(),
        Period=86400,
        Statistics=['Maximum'])
    return response

# AWSDataTransfer             0.0$
# awskms                     1.07$
# AmazonEC2                   0.0$
# AWSELB                      0.0$
# AWSCloudTrail               0.0$
# AmazonRDS                  0.21$
# AmazonS3                   0.01$
# AWSMarketplace              0.0$
# AmazonCloudWatch            0.0$
# AWSSecretsManager           0.0$
# AmazonRoute53


def parse_service_billing(service_name):
    start = time.time()
    if service_name == "Total":
        return
    response = client.get_metric_statistics(
        Namespace='AWS/Billing',
        MetricName='EstimatedCharges',
        Dimensions=[
            {
                'Name': 'Currency',
                'Value': 'USD'
            }, {
                'Name': 'ServiceName',
                "Value": service_name,
            }
        ],
        StartTime=datetime.datetime.today() - datetime.timedelta(days=1),
        EndTime=datetime.datetime.today(),
        Period=86400,
        Statistics=['Maximum'])
    elapsed_time = time.time() - start
    return response


# print(parse_service_billing("AWSDataTransfer"))
# LINEPOSTURL = os.environ['LINEPostURL']
# LINETOKEN = os.environ['LINEtoken']
# # TODO あとで消す
load_dotenv()
LINEPOSTURL = os.getenv("LINEPostURL")
LINETOKEN = os.getenv("LINEtoken")
headers = {"Authorization": "Bearer " + LINETOKEN}


def make_payload():
    metric_statistics = parse_total_service_billing()
    cost = metric_statistics['Datapoints'][0]['Maximum']
    date = metric_statistics['Datapoints'][0]['Timestamp'].strftime(
        '%Y年%m月%d日')

    message = "\n" + date

    services = parse_services()

    for service in services:
        response = parse_service_billing(service)
        value = str(response["Datapoints"][0]["Maximum"])
        message += "\n"+value+"$"+"     "+service
    message += "\n合計"+str(cost)+"$"
    return message


def lambda_handler(event, context):
    message = make_payload()
    payload = {"message": message}
    try:
        req = requests.post(LINEPOSTURL, headers=headers, params=payload)
    except requests.exceptions.RequestException as e:
        logger.error("Request failed: %s", e)


# TODO あとで消す
# lambda_handler("", "")

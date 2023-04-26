// Package mocks provide mocks for grafana api clients
package mocks

// SampleRulesJSON is a JSON response for grafana rules api
const SampleRulesJSON = `
{
    "Rodrigue": [
        {
            "name": "api gateway",
            "interval": "1m",
            "rules": [
                {
                    "expr": "",
                    "for": "5m",
                    "annotations": {
                        "__dashboardUid__": "S-ZhCPEVk",
                        "__panelId__": "2"
                    },
                    "grafana_alert": {
                        "id": 2,
                        "orgId": 1,
                        "title": "API GW 4xx",
                        "condition": "C",
                        "data": [
                            {
                                "refId": "A",
                                "queryType": "",
                                "relativeTimeRange": {
                                    "from": 21600,
                                    "to": 0
                                },
                                "datasourceUid": "q6wlwvP4z",
                                "model": {
                                    "alias": "",
                                    "datasource": {
                                        "type": "cloudwatch",
                                        "uid": "q6wlwvP4z"
                                    },
                                    "dimensions": {},
                                    "expression": "",
                                    "id": "",
                                    "intervalMs": 1000,
                                    "label": "",
                                    "matchExact": true,
                                    "maxDataPoints": 43200,
                                    "metricEditorMode": 0,
                                    "metricName": "4XXError",
                                    "metricQueryType": 0,
                                    "namespace": "AWS/ApiGateway",
                                    "period": "",
                                    "queryMode": "Metrics",
                                    "refId": "A",
                                    "region": "eu-central-1",
                                    "sqlExpression": "",
                                    "statistic": "Maximum"
                                }
                            },
                            {
                                "refId": "B",
                                "queryType": "",
                                "relativeTimeRange": {
                                    "from": 0,
                                    "to": 0
                                },
                                "datasourceUid": "-100",
                                "model": {
                                    "conditions": [
                                        {
                                            "evaluator": {
                                                "params": [],
                                                "type": "gt"
                                            },
                                            "operator": {
                                                "type": "and"
                                            },
                                            "query": {
                                                "params": [
                                                    "B"
                                                ]
                                            },
                                            "reducer": {
                                                "params": [],
                                                "type": "last"
                                            },
                                            "type": "query"
                                        }
                                    ],
                                    "datasource": {
                                        "type": "__expr__",
                                        "uid": "-100"
                                    },
                                    "expression": "A",
                                    "hide": false,
                                    "intervalMs": 1000,
                                    "maxDataPoints": 43200,
                                    "reducer": "last",
                                    "refId": "B",
                                    "type": "reduce"
                                }
                            },
                            {
                                "refId": "C",
                                "queryType": "",
                                "relativeTimeRange": {
                                    "from": 0,
                                    "to": 0
                                },
                                "datasourceUid": "-100",
                                "model": {
                                    "conditions": [
                                        {
                                            "evaluator": {
                                                "params": [
                                                    0
                                                ],
                                                "type": "gt"
                                            },
                                            "operator": {
                                                "type": "and"
                                            },
                                            "query": {
                                                "params": [
                                                    "C"
                                                ]
                                            },
                                            "reducer": {
                                                "params": [],
                                                "type": "last"
                                            },
                                            "type": "query"
                                        }
                                    ],
                                    "datasource": {
                                        "type": "__expr__",
                                        "uid": "-100"
                                    },
                                    "expression": "B",
                                    "hide": false,
                                    "intervalMs": 1000,
                                    "maxDataPoints": 43200,
                                    "refId": "C",
                                    "type": "threshold"
                                }
                            }
                        ],
                        "updated": "2023-04-22T11:39:06+02:00",
                        "intervalSeconds": 60,
                        "version": 1,
                        "uid": "bRRxjPE4z",
                        "namespace_uid": "8hANjPP4k",
                        "namespace_id": 3,
                        "rule_group": "api gateway",
                        "no_data_state": "NoData",
                        "exec_err_state": "Error"
                    }
                }
            ]
        }
    ],
    "instances": [
        {
            "name": "default",
            "interval": "1m",
            "rules": [
                {
                    "expr": "",
                    "for": "5m",
                    "annotations": {
                        "__dashboardUid__": "AtibwDE4z",
                        "__panelId__": "6"
                    },
                    "grafana_alert": {
                        "id": 1,
                        "orgId": 1,
                        "title": "CPU u",
                        "condition": "C",
                        "data": [
                            {
                                "refId": "A",
                                "queryType": "",
                                "relativeTimeRange": {
                                    "from": 21600,
                                    "to": 0
                                },
                                "datasourceUid": "q6wlwvP4z",
                                "model": {
                                    "alias": "",
                                    "datasource": {
                                        "type": "cloudwatch",
                                        "uid": "q6wlwvP4z"
                                    },
                                    "dimensions": {
                                        "InstanceId": [
                                            "i-09b19cd1749e931c4"
                                        ]
                                    },
                                    "expression": "",
                                    "id": "",
                                    "intervalMs": 1000,
                                    "label": "",
                                    "matchExact": true,
                                    "maxDataPoints": 43200,
                                    "metricEditorMode": 0,
                                    "metricName": "CPUUtilization",
                                    "metricQueryType": 0,
                                    "namespace": "AWS/EC2",
                                    "period": "",
                                    "queryMode": "Metrics",
                                    "refId": "A",
                                    "region": "eu-central-1",
                                    "sqlExpression": "",
                                    "statistic": "Average"
                                }
                            },
                            {
                                "refId": "B",
                                "queryType": "",
                                "relativeTimeRange": {
                                    "from": 0,
                                    "to": 0
                                },
                                "datasourceUid": "-100",
                                "model": {
                                    "conditions": [
                                        {
                                            "evaluator": {
                                                "params": [],
                                                "type": "gt"
                                            },
                                            "operator": {
                                                "type": "and"
                                            },
                                            "query": {
                                                "params": [
                                                    "B"
                                                ]
                                            },
                                            "reducer": {
                                                "params": [],
                                                "type": "last"
                                            },
                                            "type": "query"
                                        }
                                    ],
                                    "datasource": {
                                        "type": "__expr__",
                                        "uid": "-100"
                                    },
                                    "expression": "A",
                                    "hide": false,
                                    "intervalMs": 1000,
                                    "maxDataPoints": 43200,
                                    "reducer": "last",
                                    "refId": "B",
                                    "type": "reduce"
                                }
                            },
                            {
                                "refId": "C",
                                "queryType": "",
                                "relativeTimeRange": {
                                    "from": 0,
                                    "to": 0
                                },
                                "datasourceUid": "-100",
                                "model": {
                                    "conditions": [
                                        {
                                            "evaluator": {
                                                "params": [
                                                    0
                                                ],
                                                "type": "gt"
                                            },
                                            "operator": {
                                                "type": "and"
                                            },
                                            "query": {
                                                "params": [
                                                    "C"
                                                ]
                                            },
                                            "reducer": {
                                                "params": [],
                                                "type": "last"
                                            },
                                            "type": "query"
                                        }
                                    ],
                                    "datasource": {
                                        "type": "__expr__",
                                        "uid": "-100"
                                    },
                                    "expression": "B",
                                    "hide": false,
                                    "intervalMs": 1000,
                                    "maxDataPoints": 43200,
                                    "refId": "C",
                                    "type": "threshold"
                                }
                            }
                        ],
                        "updated": "2023-04-17T09:27:30+02:00",
                        "intervalSeconds": 60,
                        "version": 1,
                        "uid": "1Thz_vEVz",
                        "namespace_uid": "OI7z_DE4z",
                        "namespace_id": 2,
                        "rule_group": "default",
                        "no_data_state": "NoData",
                        "exec_err_state": "Error"
                    }
                }
            ]
        },
        {
            "name": "elb",
            "interval": "1m",
            "rules": [
                {
                    "expr": "",
                    "for": "5m",
                    "annotations": {
                        "__dashboardUid__": "AtibwDE4z",
                        "__panelId__": "8"
                    },
                    "grafana_alert": {
                        "id": 3,
                        "orgId": 1,
                        "title": "ELB errors",
                        "condition": "C",
                        "data": [
                            {
                                "refId": "A",
                                "queryType": "",
                                "relativeTimeRange": {
                                    "from": 21600,
                                    "to": 0
                                },
                                "datasourceUid": "q6wlwvP4z",
                                "model": {
                                    "alias": "",
                                    "datasource": {
                                        "type": "cloudwatch",
                                        "uid": "q6wlwvP4z"
                                    },
                                    "dimensions": {},
                                    "expression": "",
                                    "id": "",
                                    "intervalMs": 1000,
                                    "label": "",
                                    "matchExact": true,
                                    "maxDataPoints": 43200,
                                    "metricEditorMode": 0,
                                    "metricName": "HTTPCode_ELB_5XX_Count",
                                    "metricQueryType": 0,
                                    "namespace": "AWS/ApplicationELB",
                                    "period": "",
                                    "queryMode": "Metrics",
                                    "refId": "A",
                                    "region": "eu-central-1",
                                    "sqlExpression": "",
                                    "statistic": "Maximum"
                                }
                            },
                            {
                                "refId": "B",
                                "queryType": "",
                                "relativeTimeRange": {
                                    "from": 0,
                                    "to": 0
                                },
                                "datasourceUid": "-100",
                                "model": {
                                    "conditions": [
                                        {
                                            "evaluator": {
                                                "params": [],
                                                "type": "gt"
                                            },
                                            "operator": {
                                                "type": "and"
                                            },
                                            "query": {
                                                "params": [
                                                    "B"
                                                ]
                                            },
                                            "reducer": {
                                                "params": [],
                                                "type": "last"
                                            },
                                            "type": "query"
                                        }
                                    ],
                                    "datasource": {
                                        "type": "__expr__",
                                        "uid": "-100"
                                    },
                                    "expression": "A",
                                    "hide": false,
                                    "intervalMs": 1000,
                                    "maxDataPoints": 43200,
                                    "reducer": "last",
                                    "refId": "B",
                                    "type": "reduce"
                                }
                            },
                            {
                                "refId": "C",
                                "queryType": "",
                                "relativeTimeRange": {
                                    "from": 0,
                                    "to": 0
                                },
                                "datasourceUid": "-100",
                                "model": {
                                    "conditions": [
                                        {
                                            "evaluator": {
                                                "params": [
                                                    0
                                                ],
                                                "type": "gt"
                                            },
                                            "operator": {
                                                "type": "and"
                                            },
                                            "query": {
                                                "params": [
                                                    "C"
                                                ]
                                            },
                                            "reducer": {
                                                "params": [],
                                                "type": "last"
                                            },
                                            "type": "query"
                                        }
                                    ],
                                    "datasource": {
                                        "type": "__expr__",
                                        "uid": "-100"
                                    },
                                    "expression": "B",
                                    "hide": false,
                                    "intervalMs": 1000,
                                    "maxDataPoints": 43200,
                                    "refId": "C",
                                    "type": "threshold"
                                }
                            }
                        ],
                        "updated": "2023-04-23T15:57:01+02:00",
                        "intervalSeconds": 60,
                        "version": 1,
                        "uid": "6eqP3wPVz",
                        "namespace_uid": "OI7z_DE4z",
                        "namespace_id": 2,
                        "rule_group": "elb",
                        "no_data_state": "NoData",
                        "exec_err_state": "Error"
                    }
                }
            ]
        }
    ]
}`

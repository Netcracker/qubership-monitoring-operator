non_operator_managed = ['BLACKBOX-EXPORTER', 'CERT-EXPORTER', 'CLOUDWATCH-EXPORTER',
                        'CONFIGURATIONS-STREAMER', 'GRAPHITE-REMOTE-ADAPTER',
                        'NETWORK-LATENCY-EXPORTER', 'PROMETHEUS-ADAPTER',
                        'PROMETHEUS-ADAPTER-OPERATOR', 'PROMITOR-AGENT-SCRAPER',
                        'PROMXY', 'VERSION-EXPORTER', 'GRAFANA', 'JSON-EXPORTER']

def check_that_parameters_are_presented(environ) -> list:
    exclude_tags = [
        variable.lower() for variable in non_operator_managed
        if environ.get(variable) in ('false', None)
    ]

    operator = environ.get('OPERATOR')
    if operator == 'prometheus-operator':
        exclude_tags.append('smoke-test-vm')
    elif operator == 'victoriametrics-operator':
        exclude_tags.append('smoke-test-prometheus')

    return exclude_tags


def get_excluded_tags(environ) -> list:
    return check_that_parameters_are_presented(environ)
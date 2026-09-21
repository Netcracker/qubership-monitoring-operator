import re
import sys
import time
from os import environ
from pathlib import Path

from PlatformLibrary import PlatformLibrary

sys.path.insert(0, str(Path(__file__).resolve().parent / 'lib'))
from CheckJsonObject import describe_incomplete_rollout  # noqa: E402

namespace = environ.get('NAMESPACE')
operator = environ.get('OPERATOR')
grafana_operator = environ.get('GRAFANA')
timeout_before_start = int(environ.get('TIMEOUT-BEFORE-START'))
timeout = 300


def check_deployments_are_ready(service, label):
    """Return 1 when every Deployment labeled ``label=service`` has finished rolling out, else 0.

    ``PlatformLibrary`` counts a Deployment as active as soon as it has no unavailable replicas, which
    is true in the middle of a rolling update while the old pod is still terminating. Waiting for the
    rollout to complete keeps the Robot pod checks from seeing that transient extra pod.
    """
    deployments = [deployment for deployment in k8s_lib.get_deployment_entities(namespace)
                   if (deployment.spec.template.metadata.labels or {}).get(label, '') == service]
    if not deployments:
        print(f'No deployment with label {label}={service} found in {namespace}')
        return 0
    for deployment in deployments:
        reason = describe_incomplete_rollout(deployment)
        if reason:
            print(reason)
            return 0
    return 1


def check_statefulsets_are_ready(service):
    reason = describe_incomplete_rollout(k8s_lib.get_stateful_set(service, namespace))
    if reason:
        print(reason)
        return 0
    return 1


def check_vmagent_targets():
    try:
        pods = k8s_lib.get_pod_names_by_selector(namespace, selector={'app.kubernetes.io/name': 'vmagent'})
    except Exception as e:
        print(f'Failed to get vmagent pods: {e}')
        return 0
    if not pods:
        return 0
    pod_name = pods[0]
    try:
        last_log = k8s_lib.get_pod_logs(pod_name=pod_name, namespace=namespace,
                                        container_name='vmagent', tail_lines=200)
    except Exception as e:
        print(f'Failed to get {pod_name} pod logs: {e}')
        return 0
    matches = re.findall(r'total targets: (\d+)', last_log)
    if matches:
        return int(matches[-1])
    return 0


if __name__ == '__main__':
    try:
        k8s_lib = PlatformLibrary(managed_by_operator='true')
    except Exception as e:
        print(e)
        exit(1)
    print('Checking deployments/StatefulSets are ready')
    enabled_services = dict()
    if operator == 'prometheus-operator':
        print('Checking prometheus-operator')
        enabled_services['prometheus-operator'] = dict(ready=0, label='platform.monitoring.app', kind='deployment')
        print('Checking prometheus-k8s')
        enabled_services['prometheus-k8s'] = dict(ready=0, label='app.kubernetes.io/name', kind='statefulset')
    elif operator == 'victoriametrics-operator':
        print('Checking victoriametrics-operator')
        enabled_services['victoriametrics-operator'] = dict(ready=0, label='app.kubernetes.io/name', kind='deployment')
        print('Checking vmagent-k8s')
        enabled_services['vmagent'] = dict(ready=0, label='app.kubernetes.io/name', kind='deployment')
    else:
        print('Prometheus or victoriametrics operator is not found!')
        exit(1)
    if grafana_operator == 'true':
        print('Checking grafana')
        enabled_services['grafana'] = dict(ready=0, label='app', kind='deployment')

    timeout_start = time.time()

    all_ready = False
    while time.time() < timeout_start + timeout:
        try:
            for service in enabled_services:
                label = enabled_services[service]['label']
                kind = enabled_services[service]['kind']
                if kind == 'deployment':
                    service_is_ready = check_deployments_are_ready(service, label)
                elif kind == 'statefulset':
                    service_is_ready = check_statefulsets_are_ready(service)
                else:
                    service_is_ready = 0
                enabled_services[service]['ready'] = service_is_ready
                if service_is_ready == 0:
                    print(f'{service} deployment/statefulset is not ready')
                    raise Exception
            print('Deployments/statefulsets are ready')
            all_ready = True
            break
        except Exception:
            time.sleep(15)

    if not all_ready:
        print(f'Deployments are not ready at least {timeout} seconds')
        exit(1)
    if operator == 'victoriametrics-operator':
        timeout_start = time.time()
        vmagent_check_interval = 10
        vmagent_targets_installed = False
        while time.time() < timeout_start + timeout:
            targets = check_vmagent_targets()
            print(f'VmAgent total targets: {targets}')
            if targets >= 10:
                print('VmAgent has required amount of targets.')
                print('Sleeping 30s before starting robot tests...')
                time.sleep(30)
                print('Starting robot tests...')
                exit(0)
            print(f'VmAgent does not have required amount of targets yet, retrying in {vmagent_check_interval} seconds...')
            time.sleep(vmagent_check_interval)
        print(f'VmAgent does not have required amount of targets after {timeout} seconds')
        exit(1)

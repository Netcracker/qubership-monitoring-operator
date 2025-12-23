import time
from os import environ
import re

from PlatformLibrary import PlatformLibrary

namespace = environ.get('NAMESPACE')
operator = environ.get('OPERATOR')
grafana_operator = environ.get('GRAFANA')
timeout_before_start = int(environ.get('TIMEOUT-BEFORE-START'))
timeout = 300


def check_deployments_are_ready(service, label):
    deployments_count = k8s_lib.get_deployment_entities_count_for_service(namespace, service, label)
    ready_deployments_count = k8s_lib.get_active_deployment_entities_count_for_service(namespace, service, label)
    if deployments_count == ready_deployments_count and deployments_count != 0:
        return 1
    else:
        return 0


def check_statefulsets_are_ready(service, label):
    statefulsets_count = k8s_lib.get_stateful_set_replicas_count(service, namespace)
    ready_statefulsets_count = k8s_lib.get_stateful_set_ready_replicas_count(service, namespace)
    if statefulsets_count == ready_statefulsets_count and statefulsets_count != 0:
        return 1
    else:
        return 0

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
        last_log = k8s_lib.get_pod_logs(pod_name=pod_name, namespace=namespace, container_name='vmagent', tail_lines=200)
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
    print('Checking deployments/statefulsets are ready')
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
        print(f'Prometheus or victoriametrics operator is not found!')
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
                    service_is_ready = check_statefulsets_are_ready(service, label)
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
                print('VmAgent has required amount of targets. Sleeping 30s before starting robot tests...')
                time.sleep(30)
                print('Starting robot tests...')
                exit(0)
            print(f'VmAgent does not have required amount of targets yet, retrying in {vmagent_check_interval} seconds...')
            time.sleep(vmagent_check_interval)
        print(f'VmAgent does not have required amount of targets after {timeout} seconds')
        exit(1)

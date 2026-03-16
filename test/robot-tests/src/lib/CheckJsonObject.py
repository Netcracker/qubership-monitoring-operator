import base64
import json
import re

import yaml
from PlatformLibrary import PlatformLibrary
from urllib.parse import urlparse


def get_object_data(response):
    return response.json()['data']


def get_prometheus_target(response, target_name):
    json_object = get_object_data(response)
    matched_targets = [
        item for item in json_object.get('activeTargets', []) 
        if target_name in item.get('labels', {}).get('job', '')
    ]
    def get_path(target):
        return urlparse(target.get("scrapeUrl", "")).path
    for target in matched_targets:
        if get_path(target) == "/probe":
            return target
    for target in matched_targets:
        if get_path(target) == "/metrics":
            return target
    return False


def target_state_and_not_empty(target, health):
    scrape_url = str(target.get('scrapeUrl'))
    scrape_path = urlparse(scrape_url).path
    if scrape_path in ["/probe", "/metrics"]:
        if health in target.get('health', ''):
            return True
    return False


def check_labels_in_target(prometheus_json_target, namespace, service, job):
    list_of_labels = prometheus_json_target.get('labels')
    if namespace in list_of_labels.get('namespace'):
        if service in list_of_labels.get('service'):
            if job in list_of_labels.get('job'):
                return True
    return False


def get_metrics_of_test_app(response, job, namespace):
    if 'success' in response.json()['status']:
        json_object = get_object_data(response)
        for item in json_object.get('result'):
            metrics = item.get('metric')
            if job in metrics.get('job'):
                if namespace in metrics.get('namespace'):
                    return item
        return False


def get_metrics_by_job(response, job):
    if 'success' in response.json()['status']:
        json_object = get_object_data(response)
        for item in json_object.get('result'):
            metrics = item.get('metric')
            if job in metrics.get('job'):
                return item
        return False


def check_metrics_is_not_empty(metrics):
    if len(metrics.get('value')) != 0:
        return True


def check_cr_label_updated(response, label):
    labels = response.get('labels')
    if label in labels.get('app.kubernetes.io/component'):
        return True
    return False


def _is_install_enabled(value):
    if isinstance(value, str):
        return value.strip().lower() == "true"
    return bool(value)


def _get_http_route_hostnames(k8s_lib, namespace, base_name):
    name = base_name
    if not name.endswith("-http-route"):
        name = f"{name}-http-route"
    try:
        obj = k8s_lib.get_namespaced_custom_object(
            group="gateway.networking.k8s.io",
            version="v1",
            namespace=namespace,
            plural="httproutes",
            name=name
        )
    except Exception:
        return []
    return obj.get("spec", {}).get("hostnames") or []


def check_cr_service_exists(response, service, parentservice=None):
    if service in response.keys():
        service_child = response.get(service)
        if 'install' in service_child.keys():
            if _is_install_enabled(service_child.get('install')):
                return True
            # Explicitly disabled: do not treat host/hostnames as valid
            return False
        if service == 'route' or service == 'ingress':
            if 'host' in service_child.keys():
                if service_child.get('host') != '':
                    return True
        if service == 'httpRoute':
            hostnames = service_child.get('hostnames') or []
            if len(hostnames) > 0:
                return True
    elif parentservice:
        if parentservice in response.keys():
            service_child = response.get(parentservice)
            if service in service_child.keys():
                return True
    return False


def check_route_or_ingress(response, service_in_cr, service, namespace, parentservice=None):
    k8s_lib = PlatformLibrary()
    sub_service = ''
    if parentservice:
        sub_service = response.get('spec').get(parentservice).get(service_in_cr)
    else:
        sub_service = response.get('spec').get(service_in_cr)
    if not sub_service:
        return ""
    # Prefer HTTPRoute when configured
    http_route = sub_service.get('httpRoute')
    if http_route is not None:
        if 'install' in http_route and not _is_install_enabled(http_route.get('install')):
            return ""
        hostnames = http_route.get('hostnames') or []
        if len(hostnames) > 0:
            return hostnames[0]
        hostnames = _get_http_route_hostnames(k8s_lib, namespace, service)
        if len(hostnames) > 0:
            return hostnames[0]
        # Fallback: use ingress.host if it is provided (even if ingress is disabled).
        ingress_cfg = sub_service.get('ingress') or {}
        ingress_host = ingress_cfg.get('host') or ""
        if ingress_host:
            return ingress_host
        return ""
    if (check_cr_service_exists(sub_service, 'route')):
        return k8s_lib.get_route_url(service, namespace) or ""
    ingress_cfg = sub_service.get('ingress')
    if ingress_cfg is not None and 'install' in ingress_cfg.keys():
        if not _is_install_enabled(ingress_cfg.get('install')):
            return ""
    if (check_cr_service_exists(sub_service, 'ingress')):
        return k8s_lib.get_ingress_url(service, namespace) or ""
    return ""


def get_list_length(list_for_check):
    return len(list_for_check)


def get_dashboard_from_status(dictionary, namespace, name):
    status = dictionary.get('status')
    dashboards = status.get('dashboards')
    for i in dashboards:
        if (i.get('namespace') == namespace) and (i.get('name') == name):
            return i


def parse_yaml_file(file_path):
    return yaml.safe_load(open(file_path))


def get_dashboard_value_from_file(file_path, value):
    dict_dashboard = json.loads(parse_yaml_file(file_path).get('spec').get('json'))
    return dict_dashboard.get(value)


def update_dashboard_parameter(body: dict, key: str, updated_value: str) -> dict:
    body.get('metadata').update({key: updated_value})
    return body


def get_object_in_namespace_by_mask(list, mask):
    masked_objects = []
    pattern = mask + '-operator-*'
    for object in list:
        if mask in object.metadata.name:
            if re.search(pattern, object.metadata.name) is None:
                masked_objects.append(object)
    return masked_objects


def get_value_from_prometheus_query(response):
    metrics = get_object_data(response)
    for i in metrics['result']:
        if i['metric']['service'] == 'autoscaling-example-service':
            return i['value'][1]


def convert_to_json(str_body):
    return json.loads(str_body)


def compare_two_jsons(body_a, body_b):
    body_a, body_b = json.dumps(body_a, sort_keys=True), json.dumps(body_b, sort_keys=True)
    if body_a == body_b:
        return True


def check_tags_contain_value(tags, value):
    if value in tags:
        return True
    return False


def add_security_context_to_deployment(path_to_file, namespace):
    k8s_lib = PlatformLibrary()
    deployment = parse_yaml_file(path_to_file)
    pods = k8s_lib.get_pods(namespace)
    for pod in pods:
        if 'monitoring-operator' in pod.metadata.name:
            fs_group = pod.spec.security_context.fs_group
            run_as_user = pod.spec.security_context.run_as_user
            break
    if fs_group == None and run_as_user == None:
        return deployment
    deployment['spec']['template']['spec']['securityContext'] = dict(fsGroup=fs_group, runAsUser=run_as_user)
    return deployment


def get_pass_from_secret(secret):
    password = base64.b64decode(secret.data.get('password')).decode()
    return password


def get_username_from_secret(secret):
    username = base64.b64decode(secret.data.get('username')).decode()
    return username


def extract_config_map(configmap, key):
    config_data = configmap.to_dict() if hasattr(configmap, "to_dict") else configmap
    raw_content = config_data.get("data", {}).get(key, "")
    try:
        return yaml.safe_load(raw_content) if raw_content else {}
    except yaml.YAMLError:
        return {}

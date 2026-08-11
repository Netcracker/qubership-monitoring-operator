import base64
import json
import re
from urllib.parse import urlparse

import yaml
from PlatformLibrary import PlatformLibrary


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


def _extract_hostname(host_or_url):
    """Return the hostname part of a bare host or http(s) URL."""
    if not host_or_url:
        return ""
    value = str(host_or_url).strip()
    parsed = urlparse(value if "://" in value else f"//{value}", scheme="")
    hostname = parsed.hostname or parsed.netloc or value
    if "@" in hostname:
        hostname = hostname.rsplit("@", 1)[-1]
    if ":" in hostname and not hostname.startswith("["):
        hostname = hostname.split(":", 1)[0]
    return hostname


def _is_wildcard_host(host_or_url):
    """True if the hostname is a Kubernetes Ingress wildcard (*.suffix)."""
    hostname = _extract_hostname(host_or_url)
    return bool(hostname) and hostname.startswith("*.") and "*" not in hostname[2:]


def _hostname_matches_ingress_wildcard(hostname, wildcard_pattern):
    """True if hostname matches a K8s Ingress wildcard host rule.

    Per Kubernetes, ``*.example.com`` matches exactly one DNS label
    (``foo.example.com``) and does not match ``example.com`` or
    ``bar.foo.example.com``.
    """
    if not hostname or not wildcard_pattern:
        return False
    hostname = str(hostname).strip().lower()
    wildcard_pattern = str(wildcard_pattern).strip().lower()
    if not _is_wildcard_host(wildcard_pattern):
        return hostname == wildcard_pattern
    suffix = wildcard_pattern[1:]  # ".example.com"
    if not hostname.endswith(suffix):
        return False
    prefix = hostname[: -len(suffix)]
    return bool(prefix) and "." not in prefix and "*" not in prefix


def _ingress_spec_hosts(ingress):
    """Collect host values from Ingress.spec.rules[].host."""
    hosts = []
    spec = getattr(ingress, "spec", None)
    rules = getattr(spec, "rules", None) or []
    for rule in rules:
        host = getattr(rule, "host", None) or ""
        if host:
            hosts.append(host)
    return hosts


def _select_matching_lb_hostname(ingress, wildcard_pattern):
    """Pick a status.loadBalancer hostname that matches the wildcard pattern.

    Inspects all loadBalancer.ingress entries. IPs and non-matching hostnames
    are ignored so they are never used as HTTP Host / TLS SNI.
    """
    status = getattr(ingress, "status", None)
    if status is None:
        return ""
    load_balancer = getattr(status, "load_balancer", None)
    if load_balancer is None:
        return ""
    entries = getattr(load_balancer, "ingress", None) or []
    for entry in entries:
        hostname = getattr(entry, "hostname", None) or ""
        if hostname and _hostname_matches_ingress_wildcard(hostname, wildcard_pattern):
            return hostname
    return ""


def _matching_ingress_lb_hostname(k8s_lib, name, namespace, wildcard_pattern):
    """Fetch Ingress and return a matching LB hostname, or "".

    Propagates Kubernetes API / client errors from get_ingress. An empty
    return means status.loadBalancer has no matching hostname yet.
    """
    ingress = k8s_lib.get_ingress(name, namespace)
    return _select_matching_lb_hostname(ingress, wildcard_pattern)


def _resolve_ingress_url(url_or_host, k8s_lib, ingress_name, namespace):
    """Replace a wildcard Ingress host/URL with a matching LB hostname.

    Non-wildcard values are returned unchanged. When the host is a wildcard,
    only a status.loadBalancer hostname that matches Kubernetes wildcard
    semantics is accepted (preserving http/https scheme). IPs and unrelated
    hostnames are never used as request authority.

    Raises AssertionError when the host is a wildcard but no matching LB
    hostname is ready yet, so Robot ``Wait Until Keyword Succeeds`` can retry
    instead of treating the check as a soft skip.
    """
    if not url_or_host or not _is_wildcard_host(url_or_host):
        return url_or_host or ""
    pattern = _extract_hostname(url_or_host)
    lb_hostname = _matching_ingress_lb_hostname(
        k8s_lib, ingress_name, namespace, pattern
    )
    if not lb_hostname:
        raise AssertionError(
            f"Ingress {namespace}/{ingress_name} has no status.loadBalancer "
            f"hostname matching wildcard {pattern!r} yet"
        )
    value = str(url_or_host).strip()
    if value.lower().startswith("https://"):
        return f"https://{lb_hostname}"
    if value.lower().startswith("http://"):
        return f"http://{lb_hostname}"
    return lb_hostname


def resolve_ingress_host(host, ingress_name, namespace):
    """Robot-callable: resolve empty/wildcard host from matching Ingress LB hostname.

    Non-empty non-wildcard hosts are returned unchanged. Empty or wildcard
    hosts require a status.loadBalancer hostname that matches the Ingress
    wildcard pattern (from the argument or Ingress.spec.rules). API errors
    from get_ingress are propagated.
    """
    k8s_lib = PlatformLibrary()
    host = host or ""
    if host and not _is_wildcard_host(host):
        return host

    ingress = k8s_lib.get_ingress(ingress_name, namespace)
    if host and _is_wildcard_host(host):
        pattern = _extract_hostname(host)
    else:
        pattern = next(
            (h for h in _ingress_spec_hosts(ingress) if _is_wildcard_host(h)),
            "",
        )
    if not pattern:
        raise AssertionError(
            f"Ingress {namespace}/{ingress_name} has no wildcard host pattern "
            "available to resolve GRAFANA_HOST"
        )
    matched = _select_matching_lb_hostname(ingress, pattern)
    if not matched:
        raise AssertionError(
            f"Ingress {namespace}/{ingress_name} has no status.loadBalancer "
            f"hostname matching wildcard {pattern!r} yet"
        )
    return matched


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
        ingress_url = k8s_lib.get_ingress_url(service, namespace) or ""
        return _resolve_ingress_url(ingress_url, k8s_lib, service, namespace)
    return ""


def get_list_length(list_for_check):
    return len(list_for_check)


def get_dashboard_from_status(dictionary, namespace, name):
    # grafana-operator v5 emits status.dashboards as a list of strings of the
    # form "<namespace>/<name>/<uid>" (v4 used a list of dicts with explicit
    # namespace/name/uid keys). Normalize both shapes to a dict so callers
    # can keep using .get('uid'), .get('namespace'), .get('name').
    status = dictionary.get('status') or {}
    dashboards = status.get('dashboards') or []
    for i in dashboards:
        if isinstance(i, dict):
            if i.get('namespace') == namespace and i.get('name') == name:
                return i
            continue
        if isinstance(i, str):
            parts = i.split('/')
            if len(parts) >= 2 and parts[0] == namespace and parts[1] == name:
                uid = parts[2] if len(parts) >= 3 else None
                return {'namespace': namespace, 'name': name, 'uid': uid}
    return None


def is_resource_ready_from_conditions(obj):
    """True iff the resource is in a ready/success state.

    Checks the grafana-operator v5 schema first (status.conditions[] with at
    least one entry whose status == "True"), then falls back to the v4 schema
    (status.message == "success") for backward compatibility.
    """
    status = (obj or {}).get('status') or {}
    for cond in status.get('conditions') or []:
        if str(cond.get('status')) == 'True':
            return True
    if status.get('message') == 'success':
        return True
    return False


def is_grafana_cr_ready(obj):
    """True iff the Grafana CR reports success.

    v5: status.stageStatus == "success" (with status.stage == "complete").
    v4 fallback: status.message == "success".
    """
    status = (obj or {}).get('status') or {}
    if status.get('stageStatus') == 'success':
        return True
    if status.get('message') == 'success':
        return True
    return False


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
    if fs_group is None and run_as_user is None:
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

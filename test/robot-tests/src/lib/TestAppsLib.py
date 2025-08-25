import os
from os import environ

import urllib3
from PlatformLibrary import PlatformLibrary
from kubernetes import client
from robot.libraries.BuiltIn import BuiltIn


class TestAppsLib(object):

    def __init__(self, managed_by_operator="true"):

        urllib3.disable_warnings()
        self.k8s_lib = PlatformLibrary(managed_by_operator)
        self.namespace = os.environ.get("NAMESPACE")
        self.api_client = self.k8s_lib.k8s_api_client
        self.kubernetes_version_is_new = self.check_kubernetes_version()

    def get_service_monitor(self, name: str):

        return self.k8s_lib.get_namespaced_custom_object(
            group='monitoring.coreos.com',
            version='v1',
            namespace=self.namespace,
            plural='servicemonitors',
            name=name
        )

    def create_service_monitor(self, file_path: str):

        body = self.k8s_lib._parse_yaml_from_file(file_path)

        return self.k8s_lib.create_namespaced_custom_object(
            group='monitoring.coreos.com',
            version='v1',
            namespace=self.namespace,
            plural='servicemonitors',
            body=body
        )

    def patch_service_monitor(self, name: str, body: dict):

        return self.k8s_lib.patch_namespaced_custom_object(
            group='monitoring.coreos.com',
            version='v1',
            namespace=self.namespace,
            plural='servicemonitors',
            name=name,
            body=body
        )

    def delete_service_monitor(self, name: str):

        return self.k8s_lib.delete_namespaced_custom_object(
            group='monitoring.coreos.com',
            version='v1',
            namespace=self.namespace,
            plural='servicemonitors',
            name=name
        )

    def create_custom_metric_rule(self, file_path: str):

        body = self.k8s_lib._parse_yaml_from_file(file_path)

        return self.k8s_lib.create_namespaced_custom_object(
            group='monitoring.qubership.org',
            version='v1alpha1',
            namespace=self.namespace,
            plural='customscalemetricrules',
            body=body
        )

    def delete_custom_metric_rule(self, name: str):

        return self.k8s_lib.delete_namespaced_custom_object(
            group='monitoring.qubership.org',
            version='v1alpha1',
            namespace=self.namespace,
            plural='customscalemetricrules',
            name=name
        )

    def create_horizontal_pod_autoscaler(self, file_path: str):

        if self.kubernetes_version_is_new:
            version = 'v2'
            file_path = file_path + '/horizontal-pod-autoscaler.yaml'
        else:
            version = 'v2beta1'
            file_path = file_path + '/horizontal-pod-autoscaler-old.yaml'

        body = self.k8s_lib._parse_yaml_from_file(file_path)

        return self.k8s_lib.create_namespaced_custom_object(
            group='autoscaling',
            version=version,
            namespace=self.namespace,
            plural='horizontalpodautoscalers',
            body=body
        )

    def delete_horizontal_pod_autoscaler(self, name: str):

        if self.kubernetes_version_is_new:
            version = 'v2'
        else:
            version = 'v2beta1'

        return self.k8s_lib.delete_namespaced_custom_object(
            group='autoscaling',
            version=version,
            namespace=self.namespace,
            plural='horizontalpodautoscalers',
            name=name
        )

    def get_horizontal_pod_autoscaler_status(self, name: str):

        return self.k8s_lib.get_namespaced_custom_object_status(
            group='autoscaling',
            version='v2beta1',
            namespace=self.namespace,
            plural='horizontalpodautoscalers',
            name=name
        )

    def get_kubernetes_version(self):
        api = client.VersionApi(api_client=self.api_client)
        version_info = api.get_code().to_dict()
        git_version = version_info["git_version"]
        if git_version[0] == 'v':
            git_version = git_version[1:]
        return git_version

    def check_kubernetes_version(self):
        KEY = '1.23'
        version = self.get_kubernetes_version()
        arr_key = KEY.split(".")
        arr = version.split(".")
        for i in range(len(arr_key)):
            if arr[i] > arr_key[i]:
                return True
            elif arr[i] < arr_key[i]:
                return False
        return True

    def add_selector_to_cr(self, cr_name, namespace):
        plural = cr_name + 's'
        custom_resource = self.k8s_lib.get_namespaced_custom_object_status(
            group='monitoring.qubership.org',
            version='v1alpha1',
            namespace=namespace,
            plural=plural,
            name=cr_name
        )
        selector = {
            "matchExpressions": [
                {
                    "key": "app.kubernetes.io/component",
                    "operator": "In",
                    "values": [
                        "monitoring"
                    ]
                }
            ]
        }
        operator = environ.get('OPERATOR')
        if operator == 'victoriametrics-operator':
            custom_resource['spec']['victoriametrics']['vmAgent']['serviceMonitorSelector'] = selector
        elif operator == 'prometheus-operator':
            custom_resource['spec']['prometheus']['serviceMonitorSelector'] = selector
        else:
            BuiltIn().run_keyword('log to console',
                                  "Prometheus or victoriametrics operator is not found!")
        self.k8s_lib.patch_namespaced_custom_object(
            group='monitoring.qubership.org',
            version='v1alpha1',
            namespace=namespace,
            plural=plural,
            name=cr_name,
            body=custom_resource
        )

    def delete_selector_from_cr(self, cr_name, namespace):
        plural = cr_name + 's'
        custom_resource = self.k8s_lib.get_namespaced_custom_object_status(
            group='monitoring.qubership.org',
            version='v1alpha1',
            namespace=namespace,
            plural=plural,
            name=cr_name
        )
        operator = environ.get('OPERATOR')
        if operator == 'victoriametrics-operator':
            if custom_resource['spec']['victoriametrics']['vmAgent'].get('serviceMonitorSelector'):
                del custom_resource['spec']['victoriametrics']['vmAgent']['serviceMonitorSelector']
        elif operator == 'prometheus-operator':
            if custom_resource['spec']['prometheus'].get('serviceMonitorSelector'):
                del custom_resource['spec']['prometheus']['serviceMonitorSelector']
        else:
            BuiltIn().run_keyword('log to console',
                                  "Prometheus or victoriametrics operator is not found!")
        self.k8s_lib.replace_namespaced_custom_object(
            group='monitoring.qubership.org',
            version='v1alpha1',
            namespace=namespace,
            plural=plural,
            name=cr_name,
            body=dict(custom_resource)
        )

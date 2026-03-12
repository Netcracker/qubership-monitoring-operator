import io
import paramiko
import json
import requests
import warnings
import base64
import gzip

from kubernetes import client, config
from robot.libraries.BuiltIn import BuiltIn

warnings.simplefilter("ignore")
requests.packages.urllib3.disable_warnings()

class ConfigurationsStreamerLib:
    """
    Library for checking synchronization of dashboards between the cloud,
    external Grafana instances, and FTP server.
    """

    DASHBOARD_GROUP = "grafana.integreatly.org"
    DASHBOARD_VERSION = "v1beta1"
    DASHBOARD_PLURAL = "grafanadashboards"

    PROMETHEUS_GROUP = "monitoring.coreos.com"
    PROMETHEUS_VERSION = "v1"
    PROMETHEUS_PLURAL = "prometheusrules"

    FTP_PATHS = [
        "/var/ftp/grafana/folder1",
        "/var/ftp/grafana/folder2",
        "/var/ftp/grafana/folder3"
    ]

    ALERTMANAGER_PATHS = [
        "/var/ftp/alertmanager/folder1",
        "/var/ftp/alertmanager/folder2",
        "/var/ftp/alertmanager/folder3"
    ]

    PROMETHEUS_RULES_PATH = "/etc/monitoring/alerts"
    REQUIRED_ALERTMANAGER_FILE = "icinga_services.json"

    def __init__(self):
        self.built_in = BuiltIn()
        self._load_variables()
        self._initialize_k8s_client()
        self.built_in.log("Initialized ConfigurationsStreamerLibrary", level="INFO")

    def _load_variables(self):
        """Loads variables from Robot Framework."""
        self.grafana_basic_user = self._get_variable("${GRAFANA_BASIC_USER}")
        self.grafana_basic_pass = self._get_variable("${GRAFANA_BASIC_PASS}")
        self.ip_grafana_creds = self._get_variable("${IP_GRAFANA_CREDS}")
        self.ip_grafana_token = self._get_variable("${IP_GRAFANA_TOKEN}")
        self.ip_ftp_server = self._get_variable("${IP_FTP_SERVER}")
        self.ssh_key = self._get_variable("${SSH_KEY}")
        self.ftp_vm_user = self._get_variable("${FTP_VM_USER}")
        
        self.grafana_basic_url = f"https://{self.ip_grafana_creds}/grafana"
        self.grafana_org_url = f"https://{self.ip_grafana_token}/grafana"
        self.grafana_org_token = ""

    def _get_variable(self, var_name):
        """Retrieves a variable from Robot Framework, fails the test if missing."""
        value = self.built_in.get_variable_value(var_name)
        if not value:
            self.built_in.fail(f"Variable {var_name} is not defined.")
        return value

    def _initialize_k8s_client(self):
        """Initializes the Kubernetes API client."""
        try:
            config.load_incluster_config()
            self.built_in.log("Loaded in-cluster Kubernetes configuration", level="INFO")
        except Exception:
            try:
                config.load_kube_config()
                self.built_in.log("Loaded local Kubernetes configuration", level="INFO")
            except Exception as e:
                self.built_in.fail(f"Error loading Kubernetes configuration: {e}")
        self.custom_api = client.CustomObjectsApi()

    def set_grafana_org_token(self, token):
        """Sets the Grafana Org Token."""
        if not token:
            self.built_in.fail("Provided Grafana Org Token is empty.")
        self.grafana_org_token = token
        self.built_in.log("Grafana Org Token has been set.", level="INFO")

    def get_cloud_dashboards(self, use_cr_name=False):
        """Retrieves the list of dashboards from the cloud. Uses CR name (metadata.name) if use_cr_name=True, otherwise title."""
        try:
            dashboards_cr = self.custom_api.list_cluster_custom_object(
                group=self.DASHBOARD_GROUP,
                version=self.DASHBOARD_VERSION,
                plural=self.DASHBOARD_PLURAL
            )
            return [self._extract_dashboard_title(item, use_cr_name) for item in dashboards_cr.get("items", [])]
        except Exception as e:
            self.built_in.fail(f"Error retrieving dashboards from the cloud: {e}")

    def _extract_dashboard_title(self, item, use_cr_name=False):
        """Extracts either the CR name (metadata.name) or the title from JSON."""
        if use_cr_name:
            return item["metadata"]["name"]

        if "spec" in item:
            try:
                if "json" in item["spec"]:
                    return json.loads(item["spec"]["json"]).get("title", item["metadata"]["name"])
                if "gzipJson" in item["spec"]:
                    decompressed = gzip.decompress(base64.b64decode(item["spec"]["gzipJson"])).decode('utf-8')
                    return json.loads(decompressed).get("title", item["metadata"]["name"])
            except Exception:
                pass
        return item["metadata"]["name"]

    def get_grafana_dashboards(self, url, auth=None, token=None):
        """Retrieves the list of dashboards from Grafana via API."""
        headers = {"Authorization": f"Bearer {token}"} if token else {}
        try:
            response = requests.get(f"{url}/api/search", auth=auth, headers=headers, verify=False)
            response.raise_for_status()
            return [item.get("title") for item in response.json()]
        except requests.RequestException as e:
            self.built_in.fail(f"Error retrieving dashboards from Grafana ({url}): {e}")

    def compare_dashboards(self):
        """Compares dashboards between the cloud and two Grafana instances."""
        cloud_dashboards = self.get_cloud_dashboards()
        basic_dashboards = self.get_grafana_dashboards(self.grafana_basic_url, auth=(self.grafana_basic_user, self.grafana_basic_pass))
        org_dashboards = self.get_grafana_dashboards(self.grafana_org_url, token=self.grafana_org_token)

        missing_in_basic = set(cloud_dashboards) - set(basic_dashboards)
        missing_in_org = set(cloud_dashboards) - set(org_dashboards)

        self.built_in.log(f"Missing in Basic Auth Grafana: {missing_in_basic}", level="INFO")
        self.built_in.log(f"Missing in Org Token Grafana: {missing_in_org}", level="INFO")

        return {"missing_in_basic": list(missing_in_basic), "missing_in_org": list(missing_in_org)}

    def check_dashboards_availability(self):
        """Checks that all dashboards are available in Grafana."""
        result = self.compare_dashboards()
        missing_messages = []
        
        if result["missing_in_basic"]:
            missing_messages.append(f"Missing in Grafana (Basic Auth) {self.grafana_basic_url}: {', '.join(result['missing_in_basic'])}")
        if result["missing_in_org"]:
            missing_messages.append(f"Missing in Grafana (Org Token) {self.grafana_org_url}: {', '.join(result['missing_in_org'])}")
        
        if missing_messages:
            self.built_in.fail("\n".join(missing_messages))
        
        self.built_in.log("All dashboards are available in external Grafana.", level="INFO")
        return "All dashboards are available."

    def create_test_dashboard(self, namespace, body):
        """Creates a new Grafana Dashboard in the specified namespace."""
        try:
            self.custom_api.create_namespaced_custom_object(
                group=self.DASHBOARD_GROUP,
                version=self.DASHBOARD_VERSION,
                namespace=namespace,
                plural=self.DASHBOARD_PLURAL,
                body=body
            )
            self.built_in.log(f"Dashboard {body['metadata']['name']} created successfully in namespace {namespace}.", level="INFO")
        except Exception as e:
            self.built_in.fail(f"Error creating dashboard in namespace {namespace}: {e}")

    def update_test_dashboard(self, namespace, dashboard_name, new_title):
        """Updates the title of an existing Grafana Dashboard in the specified namespace."""
        try:
            # Get dashboard
            dashboard = self.custom_api.get_namespaced_custom_object(
                group=self.DASHBOARD_GROUP,
                version=self.DASHBOARD_VERSION,
                namespace=namespace,
                plural=self.DASHBOARD_PLURAL,
                name=dashboard_name
            )

            # Change title in JSON
            dashboard_json = json.loads(dashboard["spec"]["json"])
            dashboard_json["title"] = new_title
            dashboard["spec"]["json"] = json.dumps(dashboard_json)

            # Send update
            self.custom_api.replace_namespaced_custom_object(
                group=self.DASHBOARD_GROUP,
                version=self.DASHBOARD_VERSION,
                namespace=namespace,
                plural=self.DASHBOARD_PLURAL,
                name=dashboard_name,
                body=dashboard
            )

            self.built_in.log(f"Dashboard '{dashboard_name}' updated successfully with new title '{new_title}'.", level="INFO")

        except Exception as e:
            self.built_in.fail(f"Error updating dashboard '{dashboard_name}': {e}")

    def remove_test_dashboard(self, namespace, dashboard_name):
        """Deletes the Test Grafana Dashboard in the specified namespace."""
        try:
            self.custom_api.delete_namespaced_custom_object(
                group=self.DASHBOARD_GROUP,
                version=self.DASHBOARD_VERSION,
                namespace=namespace,
                plural=self.DASHBOARD_PLURAL,
                name=dashboard_name
            )
            self.built_in.log(f"Dashboard {dashboard_name} deleted successfully from namespace {namespace}.", level="INFO")
        except Exception as e:
            self.built_in.fail(f"Error deleting dashboard {dashboard_name} in namespace {namespace}: {e}")

    def check_dashboard_existence(self, title):
        """Checks if a specific Grafana Dashboard exists in external Grafana instances."""
        basic_dashboards = self.get_grafana_dashboards(self.grafana_basic_url, auth=(self.grafana_basic_user, self.grafana_basic_pass))
        org_dashboards = self.get_grafana_dashboards(self.grafana_org_url, token=self.grafana_org_token)

        missing_in_basic = title not in basic_dashboards
        missing_in_org = title not in org_dashboards

        if missing_in_basic:
            self.built_in.fail(f"Dashboard '{title}' is missing in Grafana (Basic Auth) {self.grafana_basic_url}")
        if missing_in_org:
            self.built_in.fail(f"Dashboard '{title}' is missing in Grafana (Org Token) {self.grafana_org_url}")

        self.built_in.log(f"Dashboard '{title}' exists in both Grafana instances.", level="INFO")
        return True

    def _create_ssh_connection(self, ip):
        """Establishes SSH connection to a given IP address."""
        try:
            pkey = paramiko.RSAKey.from_private_key(io.StringIO(self.ssh_key))
            ssh_client = paramiko.SSHClient()
            ssh_client.set_missing_host_key_policy(paramiko.AutoAddPolicy())

            ssh_client.connect(hostname=ip, username=self.ftp_vm_user, pkey=pkey)
            return ssh_client
        except Exception as e:
            self.built_in.fail(f"Failed to establish SSH connection to {ip}: {e}")

    def _execute_command_on_vm(self, ssh_client, command):
        """Executes a shell command on a remote server using sudo."""
        command = f"sudo {command}"
        stdin, stdout, stderr = ssh_client.exec_command(command)
        return stdout.channel.recv_exit_status(), stdout.read().decode()

    def get_ftp_dashboards(self):
        """Retrieves the list of dashboard files from the FTP server using sudo."""
        ssh_client = self._create_ssh_connection(self.ip_ftp_server)
        ftp_dashboards = {}
        missing_paths = []
        
        try:
            for path in self.FTP_PATHS:
                command = f"ls {path} 2>/dev/null"
                status, output = self._execute_command_on_vm(ssh_client, command)
                
                if status == 0 and output.strip():
                    files = [f.split("]")[-1] for f in output.strip().split("\n")]
                    ftp_dashboards[path] = files
                else:
                    ftp_dashboards[path] = []
                    missing_paths.append(path)
        finally:
            ssh_client.close()
        
        return ftp_dashboards, missing_paths

    def compare_ftp_dashboards(self):
        """Compares dashboards between the cloud and the FTP server using CR name."""
        cloud_dashboards = self.get_cloud_dashboards(use_cr_name=True)
        ftp_dashboards, missing_paths = self.get_ftp_dashboards()
        
        missing_on_ftp = {}
        for path, files in ftp_dashboards.items():
            missing_files = [dashboard for dashboard in cloud_dashboards if f"{dashboard}.json" not in files]
            if missing_files:
                missing_on_ftp[path] = missing_files
        
        if missing_on_ftp:
            error_messages = [f"Missing dashboards in {path}: {', '.join(missing)}" for path, missing in missing_on_ftp.items()]
            self.built_in.fail("\n".join(error_messages))
        
        if missing_paths:
            self.built_in.fail(f"Missing or empty FTP directories: {', '.join(missing_paths)}")
        
        self.built_in.log("All cloud dashboards are present on FTP.", level="INFO")
        return "All dashboards are available on FTP."

    def check_alertmanager_files(self):
        """Checks that icinga_services.json exists and is not empty in all alertmanager directories."""
        ssh_client = self._create_ssh_connection(self.ip_ftp_server)
        missing_or_empty_files = []
        missing_paths = []
        
        try:
            for path in self.ALERTMANAGER_PATHS:
                file_path = f"{path}/{self.REQUIRED_ALERTMANAGER_FILE}"
                command = f"test -s {file_path} || echo {file_path}"
                status, output = self._execute_command_on_vm(ssh_client, command)
                
                if output.strip():
                    missing_or_empty_files.append(output.strip())
                
                command = f"ls {path} 2>/dev/null"
                status, output = self._execute_command_on_vm(ssh_client, command)
                if status != 0 or not output.strip():
                    missing_paths.append(path)
        finally:
            ssh_client.close()
        
        if missing_paths:
            self.built_in.fail(f"Missing or empty Alertmanager directories: {', '.join(missing_paths)}")
        
        if missing_or_empty_files:
            self.built_in.fail(f"Missing or empty Alertmanager files: {', '.join(missing_or_empty_files)}")
        
        self.built_in.log("All required Alertmanager files are present and non-empty.", level="INFO")
        return "All Alertmanager files verified successfully."

    def get_cloud_prometheusrules(self):
        """Retrieves PrometheusRules only from specific namespaces."""
        try:
            monitoring_namespace = self.built_in.get_variable_value("${namespace}")
            target_namespaces = [monitoring_namespace, "cert-manager"]
            all_rules = {}

            for ns in target_namespaces:
                rules_cr = self.custom_api.list_namespaced_custom_object(
                    group=self.PROMETHEUS_GROUP,
                    version=self.PROMETHEUS_VERSION,
                    namespace=ns,
                    plural=self.PROMETHEUS_PLURAL
                )
                for item in rules_cr.get("items", []):
                    all_rules[f"{item['metadata']['name']}.yaml"] = ns

            return all_rules
        except Exception as e:
            self.built_in.fail(f"Error retrieving PrometheusRules from the cloud: {e}")

    def get_prometheus_rules_files(self, ip):
        """Retrieves PrometheusRules files from a remote server, filtering by namespaces."""
        ssh_client = self._create_ssh_connection(ip)
        rules_files = {}
        monitoring_namespace = self.built_in.get_variable_value("${namespace}")
        target_namespaces = [monitoring_namespace, "cert-manager"]

        try:
            command = f"ls {self.PROMETHEUS_RULES_PATH}"
            status, output = self._execute_command_on_vm(ssh_client, command)

            if status == 0:
                for file in output.strip().split("\n"):
                    if file != "vm_alerts.yaml":
                        parts = file.split("_")
                        if len(parts) >= 3:
                            namespace = parts[-2]
                            rule_name = parts[-1]
                            if namespace in target_namespaces:
                                rules_files[rule_name] = namespace
                        else:
                            self.built_in.log(f"WARNING: Unable to parse file name '{file}', skipping.", level="WARN")
            else:
                self.built_in.fail(f"Failed to retrieve PrometheusRules files from {ip}: {output}")
        finally:
            ssh_client.close()

        return rules_files

    def compare_prometheus_rules(self):
        """Compares PrometheusRules from specific namespaces with remote servers."""
        cloud_rules = self.get_cloud_prometheusrules()
        rules_server1 = self.get_prometheus_rules_files(self.ip_grafana_creds)
        rules_server2 = self.get_prometheus_rules_files(self.ip_grafana_token)

        missing_rules = {
            self.ip_grafana_creds: set(cloud_rules.keys()) - set(rules_server1.keys()),
            self.ip_grafana_token: set(cloud_rules.keys()) - set(rules_server2.keys())
        }
        
        errors = [f"Missing on {server}: {', '.join(missing)}" for server, missing in missing_rules.items() if missing]
        
        if errors:
            self.built_in.fail("\n".join(errors))
        
        self.built_in.log("All PrometheusRules from monitoring and cert-manager namespaces are present on both servers.", level="INFO")
        return "All PrometheusRules verified successfully."
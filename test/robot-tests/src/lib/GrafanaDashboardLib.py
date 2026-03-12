# Library for create/delete/replace GrafanaDashboard via Kubernetes API (grafana.integreatly.org/v1beta1).
# Used when kubectl is not available in the test pod.

from kubernetes import client, config


GROUP = "grafana.integreatly.org"
VERSION = "v1beta1"
PLURAL = "grafanadashboards"


class GrafanaDashboardLib:
    ROBOT_LIBRARY_SCOPE = "GLOBAL"

    def __init__(self):
        self._api = None

    def _get_api(self):
        if self._api is None:
            try:
                config.load_incluster_config()
            except Exception:
                config.load_kube_config()
            self._api = client.CustomObjectsApi()
        return self._api

    def create_dashboard(self, namespace, body):
        """Create GrafanaDashboard in namespace. body is dict (e.g. from Parse Yaml File)."""
        api = self._get_api()
        name = body.get("metadata", {}).get("name")
        if not name:
            raise ValueError("GrafanaDashboard body must have metadata.name")
        api.create_namespaced_custom_object(
            group=GROUP,
            version=VERSION,
            namespace=namespace,
            plural=PLURAL,
            body=body,
        )
        return {"status": "Success", "message": "created"}

    def delete_dashboard(self, namespace, name):
        """Delete GrafanaDashboard by name. Ignores 404."""
        api = self._get_api()
        try:
            api.delete_namespaced_custom_object(
                group=GROUP,
                version=VERSION,
                namespace=namespace,
                plural=PLURAL,
                name=name,
            )
        except client.rest.ApiException as e:
            if e.status != 404:
                raise
        return {"status": "Success"}

    def replace_dashboard(self, namespace, body):
        """Replace (update) GrafanaDashboard without dropping operator finalizers."""
        api = self._get_api()
        name = body.get("metadata", {}).get("name")
        if not name:
            raise ValueError("GrafanaDashboard body must have metadata.name")
        existing = api.get_namespaced_custom_object(
            group=GROUP,
            version=VERSION,
            namespace=namespace,
            plural=PLURAL,
            name=name,
        )
        existing_meta = existing.get("metadata", {})
        body_meta = body.setdefault("metadata", {})

        # Preserve metadata that the grafana-operator adds during reconciliation.
        # If replace removes finalizers, the dashboard CR disappears from Kubernetes
        # without Grafana cleanup and the dashboard stays orphaned in Grafana.
        if "resourceVersion" not in body_meta and existing_meta.get("resourceVersion"):
            body_meta["resourceVersion"] = existing_meta["resourceVersion"]
        if "finalizers" not in body_meta and existing_meta.get("finalizers"):
            body_meta["finalizers"] = existing_meta["finalizers"]
        if "annotations" not in body_meta and existing_meta.get("annotations"):
            body_meta["annotations"] = existing_meta["annotations"]
        api.replace_namespaced_custom_object(
            group=GROUP,
            version=VERSION,
            namespace=namespace,
            plural=PLURAL,
            name=name,
            body=body,
        )
        return {"status": "Success"}

from grafana_api.grafana_face import GrafanaFace


class GrafanaApiLib(object):
    def __init__(self, url, g_user, g_password):
        grafana_api = GrafanaFace(
            auth=(g_user, g_password),
            host=url,
            port=None,
            protocol="http",
            verify=False,
            timeout=5.0
        )
        self._grafana_api = grafana_api

    def search_dashboards(self):
        return self._grafana_api.search.search_dashboards(tag='applications')

    def find_dashboard(self, uid):
        dashboard = self._grafana_api.dashboard.get_dashboard(uid)
        if isinstance(dashboard, dict) and dashboard.get('uid') != "" and dashboard.get(
                'title') != "" and dashboard.get('folderName') != "":
            return dashboard
        else:
            return None

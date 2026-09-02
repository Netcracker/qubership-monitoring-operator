from __future__ import annotations

import base64
import json
import socket
import ssl
from typing import Any
from urllib import error, parse, request


class QueryError(RuntimeError):
    """Raised when a backend request cannot be completed."""


class HttpClient:
    def __init__(
        self,
        base_url: str,
        *,
        user: str = "",
        password: str = "",
        insecure_skip_verify: bool = False,
    ) -> None:
        self.base_url = base_url.rstrip("/")
        self.headers: dict[str, str] = {}
        if user or password:
            credentials = base64.b64encode(f"{user}:{password}".encode("utf-8")).decode("ascii")
            self.headers["Authorization"] = f"Basic {credentials}"
        self.context = ssl._create_unverified_context() if insecure_skip_verify else None  # noqa: SLF001

    def call(self, path: str, query: dict[str, str]) -> Any:
        url = self.base_url + path + "?" + parse.urlencode(query)
        req = request.Request(url, headers=self.headers, method="GET")
        try:
            with request.urlopen(req, context=self.context, timeout=60) as response:
                return response.read().decode("utf-8")
        except error.HTTPError as exc:
            body = exc.read().decode("utf-8", errors="replace")
            raise QueryError(f"GET {path} returned HTTP {exc.code}: {body}") from exc
        except error.URLError as exc:
            raise QueryError(f"cannot connect to {self.base_url}: {exc.reason}") from exc
        except (TimeoutError, socket.timeout) as exc:
            raise QueryError(f"GET {path} timed out while connecting to {self.base_url}") from exc


class VictoriaMetricsClient:
    def __init__(self, client: HttpClient, query_time: str) -> None:
        self.client = client
        self.query_time = query_time

    def read_json(self, path: str, query: dict[str, str] | None = None) -> Any:
        raw = self.client.call(path, query or {})
        return json.loads(raw)

    def query(self, expression: str) -> Any:
        response = self.read_json("/api/v1/query", {"query": expression, "time": self.query_time})
        if response.get("status") != "success":
            raise QueryError(f"VictoriaMetrics query failed for {expression}: {response}")
        return response.get("data", {}).get("result", [])

    def safe_query(self, expression: str) -> Any:
        try:
            return self.query(expression)
        except (QueryError, json.JSONDecodeError) as exc:
            return {"error": str(exc)}

    def tsdb_stats(self, limit: int) -> dict[str, Any]:
        response = self.read_json("/api/v1/status/tsdb", {"limit": str(limit)})
        if response.get("status") != "success":
            raise QueryError(f"VictoriaMetrics TSDB stats request failed: {response}")
        data = response.get("data")
        if not isinstance(data, dict):
            raise QueryError(f"VictoriaMetrics TSDB stats returned unexpected payload: {response}")
        return data

    def safe_tsdb_stats(self, limit: int) -> Any:
        try:
            return self.tsdb_stats(limit)
        except (QueryError, json.JSONDecodeError) as exc:
            return {"error": str(exc)}

    def metric_names_stats(self, limit: int) -> Any:
        response = self.read_json("/api/v1/status/metric_names_stats", {"limit": str(limit)})
        if response.get("status") != "success":
            raise QueryError(f"VictoriaMetrics metric names stats request failed: {response}")
        return response.get("data") if "data" in response else response

    def safe_metric_names_stats(self, limit: int) -> Any:
        try:
            return self.metric_names_stats(limit)
        except (QueryError, json.JSONDecodeError) as exc:
            return {"error": str(exc)}

    def top_queries(self, top_n: int, max_lifetime: str) -> Any:
        response = self.read_json("/api/v1/status/top_queries", {"topN": str(top_n), "maxLifetime": max_lifetime})
        if not isinstance(response, dict):
            raise QueryError(f"VictoriaMetrics top queries request returned unexpected payload: {response}")
        if "status" in response:
            if response.get("status") != "success":
                raise QueryError(f"VictoriaMetrics top queries request failed: {response}")
            data = response.get("data")
            if not isinstance(data, dict):
                raise QueryError(f"VictoriaMetrics top queries request returned unexpected payload: {response}")
            return data
        if any(
            key in response
            for key in ("topByCount", "topByAvgDuration", "topBySumDuration", "topByAvgMemoryUsage")
        ):
            return response
        raise QueryError(f"VictoriaMetrics top queries request returned unexpected payload: {response}")

    def safe_top_queries(self, top_n: int, max_lifetime: str) -> Any:
        try:
            return self.top_queries(top_n, max_lifetime)
        except (QueryError, json.JSONDecodeError) as exc:
            return {"error": str(exc)}

    def label_values(self, label_name: str) -> list[str]:
        response = self.read_json(
            f"/api/v1/label/{parse.quote(label_name, safe='')}/values",
            {"start": self.query_time, "end": self.query_time},
        )
        if response.get("status") != "success":
            raise QueryError(f"VictoriaMetrics label values request failed for {label_name}: {response}")
        data = response.get("data")
        if not isinstance(data, list):
            raise QueryError(
                f"VictoriaMetrics label values request returned unexpected payload for {label_name}: {response}"
            )
        return [item for item in data if isinstance(item, str)]

    def safe_label_values(self, label_name: str) -> Any:
        try:
            return self.label_values(label_name)
        except (QueryError, json.JSONDecodeError) as exc:
            return {"error": str(exc)}

    def series(self, match: str, limit: int) -> list[dict[str, str]]:
        query = {"match[]": match, "start": self.query_time, "end": self.query_time}
        if limit > 0:
            query["limit"] = str(limit)
        response = self.read_json("/api/v1/series", query)
        if response.get("status") != "success":
            raise QueryError(f"VictoriaMetrics series request failed for {match}: {response}")
        data = response.get("data")
        if not isinstance(data, list):
            raise QueryError(f"VictoriaMetrics series request returned unexpected payload for {match}: {response}")
        return [item for item in data if isinstance(item, dict)]

    def safe_series(self, match: str, limit: int) -> Any:
        try:
            return self.series(match, limit)
        except (QueryError, json.JSONDecodeError) as exc:
            return {"error": str(exc)}

import sys
import types
import unittest

# PlatformLibrary is provided by the integration-tests base image, not requirements.txt.
_fake_platform = types.ModuleType("PlatformLibrary")


class _FakePlatformLibrary:
    def __init__(self, *args, **kwargs):
        pass


_fake_platform.PlatformLibrary = _FakePlatformLibrary
sys.modules["PlatformLibrary"] = _fake_platform

import CheckJsonObject as cjo  # noqa: E402


class _Obj:
    def __init__(self, **fields):
        self.__dict__.update(fields)


# Class names match the kubernetes client models that describe_incomplete_rollout dispatches on.
class V1Deployment(_Obj):
    pass


class V1StatefulSet(_Obj):
    pass


class V1DaemonSet(_Obj):
    pass


def _deployment(replicas=1, generation=2, observed=2, updated=1, total=1, available=1):
    return V1Deployment(
        metadata=_Obj(name="grafana-deployment", generation=generation),
        spec=_Obj(replicas=replicas),
        status=_Obj(observed_generation=observed, updated_replicas=updated,
                    replicas=total, available_replicas=available),
    )


def _stateful_set(replicas=1, generation=1, observed=1, updated=1, ready=1,
                  current_revision="rev-1", update_revision="rev-1"):
    return V1StatefulSet(
        metadata=_Obj(name="prometheus-k8s", generation=generation),
        spec=_Obj(replicas=replicas),
        status=_Obj(observed_generation=observed, updated_replicas=updated, ready_replicas=ready,
                    current_revision=current_revision, update_revision=update_revision),
    )


def _daemon_set(desired=3, generation=1, observed=1, updated=3, available=3):
    return V1DaemonSet(
        metadata=_Obj(name="node-exporter", generation=generation),
        spec=_Obj(),
        status=_Obj(observed_generation=observed, desired_number_scheduled=desired,
                    updated_number_scheduled=updated, number_available=available),
    )


class TestDescribeIncompleteRollout(unittest.TestCase):
    def test_complete_deployment(self):
        self.assertEqual(cjo.describe_incomplete_rollout(_deployment()), "")

    def test_deployment_generation_not_observed(self):
        reason = cjo.describe_incomplete_rollout(_deployment(generation=3, observed=2))
        self.assertIn("generation 3 is not observed yet", reason)

    def test_deployment_surge_pod_still_present(self):
        # New pod is up, old pod is terminating: total exceeds the desired count.
        reason = cjo.describe_incomplete_rollout(_deployment(total=2))
        self.assertEqual(reason, "V1Deployment grafana-deployment: total 2 of 1 replicas")

    def test_deployment_new_pod_not_available(self):
        reason = cjo.describe_incomplete_rollout(_deployment(available=None))
        self.assertEqual(reason, "V1Deployment grafana-deployment: available 0 of 1 replicas")

    def test_deployment_old_replica_set_still_serving(self):
        reason = cjo.describe_incomplete_rollout(_deployment(updated=0))
        self.assertEqual(reason, "V1Deployment grafana-deployment: updated 0 of 1 replicas")

    def test_deployment_none_counters_read_as_zero(self):
        deployment = _deployment(observed=None, updated=None, total=None, available=None)
        self.assertIn("generation 2 is not observed yet (observed 0)",
                      cjo.describe_incomplete_rollout(deployment))

    def test_complete_stateful_set(self):
        self.assertEqual(cjo.describe_incomplete_rollout(_stateful_set()), "")

    def test_stateful_set_revision_rolling_out(self):
        reason = cjo.describe_incomplete_rollout(_stateful_set(update_revision="rev-2"))
        self.assertEqual(reason, "V1StatefulSet prometheus-k8s: revision rev-2 is rolling out over rev-1")

    def test_stateful_set_replica_not_ready(self):
        reason = cjo.describe_incomplete_rollout(_stateful_set(replicas=2, updated=2, ready=1))
        self.assertEqual(reason, "V1StatefulSet prometheus-k8s: ready 1 of 2 replicas")

    def test_complete_daemon_set(self):
        self.assertEqual(cjo.describe_incomplete_rollout(_daemon_set()), "")

    def test_daemon_set_node_not_updated(self):
        reason = cjo.describe_incomplete_rollout(_daemon_set(updated=2))
        self.assertEqual(reason, "V1DaemonSet node-exporter: updated 2 of 3 replicas")

    def test_daemon_set_node_not_available(self):
        reason = cjo.describe_incomplete_rollout(_daemon_set(available=2))
        self.assertEqual(reason, "V1DaemonSet node-exporter: available 2 of 3 replicas")


class TestExcludeTerminatingPods(unittest.TestCase):
    def test_drops_pods_with_deletion_timestamp(self):
        running = _Obj(metadata=_Obj(name="grafana-deployment-new", deletion_timestamp=None))
        terminating = _Obj(metadata=_Obj(name="grafana-deployment-old", deletion_timestamp="2026-09-18T14:33:00Z"))
        self.assertEqual(cjo.exclude_terminating_pods([running, terminating]), [running])

    def test_empty_list(self):
        self.assertEqual(cjo.exclude_terminating_pods([]), [])


if __name__ == "__main__":
    unittest.main()

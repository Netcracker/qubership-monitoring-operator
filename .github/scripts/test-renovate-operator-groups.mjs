import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';

const [logFile] = process.argv.slice(2);
assert.ok(logFile, 'Usage: node test-renovate-operator-groups.mjs <renovate-debug-log>');

const renovateConfig = JSON.parse(readFileSync('renovate.json', 'utf8'));
const records = readFileSync(logFile, 'utf8')
  .split('\n')
  .flatMap((line) => {
    try {
      return [JSON.parse(line)];
    } catch {
      return [];
    }
  });
const extraction = records.findLast((record) => record.msg === 'Extracted dependencies');
assert.ok(extraction?.packageFiles, 'Renovate did not emit an Extracted dependencies record');

const dependencies = Object.values(extraction.packageFiles)
  .flat()
  .flatMap((packageFile) => packageFile.deps.map((dependency) => ({ ...dependency, packageFile: packageFile.packageFile })));

function normalizeVersion(version) {
  return version.replace(/^config-reloader-v|^v/, '');
}

function assertGroup(groupName, groupDependencies, packageNames, expectedFiles) {
  const rule = renovateConfig.packageRules.find((candidate) => candidate.groupName === groupName);
  assert.ok(rule, `${groupName}: package rule is missing`);
  assert.equal(
    rule.minimumGroupSize,
    groupDependencies.length,
    `${groupName}: minimumGroupSize must cover every extracted reference`,
  );
  assert.deepEqual(rule.matchPackageNames, packageNames, `${groupName}: package matchers changed without updating this contract`);
  assert.deepEqual(
    groupDependencies.map(({ packageFile }) => packageFile).sort(),
    expectedFiles.sort(),
    `${groupName}: extracted files changed without updating this contract`,
  );
  assert.equal(
    new Set(groupDependencies.map(({ currentValue }) => normalizeVersion(currentValue))).size,
    1,
    `${groupName}: extracted references use different versions`,
  );
}

const prometheusDependencies = dependencies.filter(
  ({ depName }) =>
    depName === 'github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring' ||
    depName === 'prometheus-operator/prometheus-operator',
);

assert.equal(prometheusDependencies.length, 4, 'Prometheus Operator: expected four extracted references');
assert.equal(
  prometheusDependencies.filter(({ datasource }) => datasource === 'github-releases').length,
  3,
  'Prometheus Operator: image and CRD-source references must use GitHub release timestamps',
);
const prometheusImageReferences = prometheusDependencies.filter(
  ({ packageFile }) => packageFile === 'charts/qubership-monitoring-operator/templates/_images.tpl',
);
assert.equal(
  prometheusImageReferences.filter(({ replaceString }) => replaceString.includes('prometheus-operator:v')).length,
  1,
  'Prometheus Operator: expected one operator image reference',
);
assert.equal(
  prometheusImageReferences.filter(({ replaceString }) => replaceString.includes('prometheus-config-reloader:v')).length,
  1,
  'Prometheus Operator: expected one config-reloader image reference',
);
assertGroup(
  'Prometheus Operator',
  prometheusDependencies,
  ['github.com/prometheus-operator/prometheus-operator{/,}**', 'prometheus-operator/prometheus-operator'],
  [
    'go.mod',
    'Makefile',
    'charts/qubership-monitoring-operator/templates/_images.tpl',
    'charts/qubership-monitoring-operator/templates/_images.tpl',
  ],
);

const victoriaMetricsDependencies = dependencies.filter(
  ({ depName }) => depName === 'github.com/VictoriaMetrics/operator/api' || depName === 'VictoriaMetrics/operator',
);

assert.equal(victoriaMetricsDependencies.length, 9, 'VictoriaMetrics Operator: expected nine extracted references');
assert.equal(
  victoriaMetricsDependencies.filter(({ datasource }) => datasource === 'github-releases').length,
  8,
  'VictoriaMetrics Operator: image and CRD-source references must use GitHub release timestamps',
);
const victoriaMetricsTemplateReferences = victoriaMetricsDependencies.filter(
  ({ packageFile }) => packageFile === 'charts/qubership-monitoring-operator/templates/_images.tpl',
);
assert.equal(
  victoriaMetricsTemplateReferences.filter(({ replaceString }) => replaceString.includes('operator:v')).length,
  1,
  'VictoriaMetrics Operator: expected one operator image reference',
);
assert.equal(
  victoriaMetricsTemplateReferences.filter(({ replaceString }) => replaceString.includes('operator:config-reloader-v')).length,
  5,
  'VictoriaMetrics Operator: expected five config-reloader image references',
);
assertGroup(
  'VictoriaMetrics Operator',
  victoriaMetricsDependencies,
  ['github.com/VictoriaMetrics/operator{/,}**', 'VictoriaMetrics/operator'],
  [
    'go.mod',
    'Makefile',
    'charts/qubership-monitoring-operator/templates/_images.tpl',
    'charts/qubership-monitoring-operator/templates/_images.tpl',
    'charts/qubership-monitoring-operator/templates/_images.tpl',
    'charts/qubership-monitoring-operator/templates/_images.tpl',
    'charts/qubership-monitoring-operator/templates/_images.tpl',
    'charts/qubership-monitoring-operator/templates/_images.tpl',
    'controllers/victoriametrics/vmoperator/assets/deployment.yaml',
  ],
);

console.log('Managed operator Renovate groups are complete and version-aligned.');

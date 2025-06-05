Grafana removes support of plugins based on AngularJS since v11, which means that some plugins will stop working from
then on. This document describes how to migrate from the old AngularJS plugins to the new, mostly core plugins.

# Compatibility of Grafana versions with new core plugins

If you are not using the latest major version of Grafana 10.x, you may be missing some of the core plugins
described below. This means that you will only be able to migrate after upgrading to a new version of Grafana.

If your current version of Grafana does not contain the plugin you would like to migrate to, you will have to recreate
the panel in the new version of Grafana.

List of **some new core plugins** that were added in the latest major Grafana releases (v7.x or higher):

| Plugin         | Minimal Grafana version |
|----------------|-------------------------|
| Canvas         | 9.2                     |
| Geomap         | 8.1                     |
| Bar chart      | 8.0                     |
| State timeline | 8.0                     |
| Status history | 8.0                     |
| Histogram      | 8.0                     |
| Time series    | 7.4                     |
| Table          | 7.0                     |

**NOTE 1:** Most plugins initially appear in Grafana as beta versions.
When compiling the table above, we focused specifically on the first appearance of beta versions.

**NOTE 2:** Please keep in mind that plugins are updated and expanded throughout Grafana updates. This means that some
features on the new panels can be not supported in your version of Grafana, even if the certain type of the panel is
already presented.

**NOTE 3:** Grafana does not have a single list of plugins with minimum supported versions of Grafana,
so **the information above** was taken from various sources (mostly Release Notes) and **is not 100% accurate**.
The table above is provided to demonstrate the absence of some plugins in older versions of Grafana.

# AngularJS plugins and their possible alternatives

You can find a list of deprecated plugins in
[the official Grafana docs](https://grafana.com/docs/grafana/latest/developers/angular_deprecation/angular-plugins/).

Often panels based on new plugins look little different from old ones, so **we recommend checking the type of each panel
on your dashboard**. You can do this by going into `Edit` mode for the panel and checking its type in
the top right corner. This especially applies to core plugins such as `Graph (old)` and `Table (old)`.

## Core plugins with automatic migration

Certain legacy Grafana panel plugins automatically update or can be migrated by pressing the single `Migrate` button
on the panel settings to their React-based replacements when Angular support is disabled.

List of plugins that support this feature:

| Plugin      | Migration target |
|-------------|------------------|
| Graph (old) | Time Series      |
| Singlestat  | Stat             |
| Stat (old)  | Stat             |
| Table (old) | Table            |
| Worldmap    | Geomap           |

Some of the plugins (e.g. Singlestat) were migrated automatically since the old plugin is no longer available in the
current version of Grafana.

Some of the plugins (e.g. old Graph, old Table) require semi-automatic migration:

1. Pick the panel based on the old plugin and press `Edit`
2. Press `Migrate` button in `Display` section of settings

Please **pay attention to the warnings directly below the `Migrate` button**, which contain information about possible
changes on the panel after migration.

For example, you can see the following notes about migration to the new `Table` plugin:

* Sorting is not persisted after migration.
* Links that specify cell values will need to be updated manually after migration.

## Third-party plugins

Third-party plugins require a custom approach.

### Pie chart (old)

**Plugin ID:** grafana-piechart-panel

**Alternatives:**

* Pie chart (core)

**Migration path:**

1. Open your panel in `Edit` mode
2. Find `Options` section in panel settings and click on `Migrate to the new pie chart plugin` button

### Status Panel

**Plugin ID:** vonage-status-panel

**Ref:** <https://github.com/Vonage/Grafana_Status_panel>

**Alternatives:**

* Stat (core)
* [Polystat](https://grafana.com/grafana/plugins/grafana-polystat-panel/)
* [Status Overview Panel](https://github.com/serrrios/Status-Overview-Panel)

**Migration path:**

The plugin doesn't support the latest versions of Grafana officially, and it can work incorrectly,
so it's hard to create a precise migration path.

You can replace most of the functionality of your Status Panel with in-built `Stat` type of panel.
`Stat` supports multiple queries, so you can keep them unchanged. The following settings can be approximately copied
into the new panel:

* `Alias` for each query in the panel options -> `Legend` in the query options
* `Aggregation` for each query in the panel options -> `Calculation` in `Value options` section of the panel options
* Any threshold settings -> General threshold settings in the panel options or overriding thresholds for each metric
* `Metric URL` -> `Data links` (general or overridden)

### FlowCharting

**Plugin ID:** agenty-flowcharting-panel

**Ref:** <https://github.com/algenty/grafana-flowcharting>

**Alternatives:**

* Canvas (core)
* [Diagram](https://github.com/jdbranham/grafana-diagram)

**Migration path:**

There is no other plugin with the same capabilities at this moment, so the only recommendation is to try
to repeat your Draw.io chart in the in-built Canvas type of panel. Canvas plugin has an animated
[documentation](https://grafana.com/docs/grafana/latest/panels-visualizations/visualizations/canvas/).

Also, you can check the progress of the moving flowcharting features to the Canvas core plugin in
["Canvas: Request flowcharting features to be
implemented/migrated into canvas plugin"](https://github.com/grafana/grafana/issues/79874)
issue in the Grafana GitHub repository.

### Singlestat Math

**Plugin ID:** blackmirror1-singlestat-math-panel

**Ref:** <https://github.com/black-mirror-1/singlestat-math>

**Alternatives:**

* Stat (core)
* Gauge (core)

**Migration path:**

Panels cannot be automatically migrated from `Singlestat Math` to `Stat` or `Gauge` automatically, but most of the
options can be moved to the new panel manually. If your Singlestat panel has `Gauge.Show` option as enabled, you should
choose a `Gauge` type, otherwise make the new panel as `Stat`.

Setting mappings from `Singlestat Math` to `Stat` or `Gauge` that can be migrated:

* `Stat` setting from `Value` section -> `Calculation` setting from `Value options` section
* `Unit` and `Decimals` settings from `Value` section -> `Unit` and `Decimals` settings from `Standard options` section
* `Font size` setting from `Value` section -> `Text size` section
* `Coloring` -> `Thresholds`
* `Value mappings` -> `Value mappings`

`Query Math` should be copied by changing query on the new panel. Prometheus/VictoriaMetrics query languages support
any type of math operators and functions.

### Statusmap

**Plugin ID:** flant-statusmap-panel

**Ref:** <https://github.com/flant/grafana-statusmap>

**Alternatives:**

* Status history (core)
* Heatmap (core)

**Migration path:**

`Statusmap` panel can be migrated to `Status history` manually with good approximation.

If you want to display statuses with different colors and the legend correctly on `Status history` panel,
you should set Values and Colors in `Value mappings`, enable `Legend.Visibility` option and then remove all thresholds.
Without the last step the legend can be displayed incorrectly.

Size of the cells and opacity can be configured in `Status history` section of settings. `Items` options can be replaced
with `Data links` in some cases.

### SVG

**Plugin ID:** marcuscalidus-svg-panel

**Ref:** <https://github.com/MarcusCalidus/marcuscalidus-svg-panel>

**Alternatives:**

* Canvas (core)
* Colored SVG

**Migration path:**

`Canvas` provides good opportunities for creating convenient visualizations for any processes. This type of panel
appeared in the latest versions of Grafana core and is actively being improved.

There is no any practical advices about migration to the `Canvas`, you can simply try to recreate your panel with
available instruments and elements.

Unfortunately, Canvas does not currently allow you to use your custom SVG code. Adding individual SVG icons is
possible by specifying a URL, that is, the icon must first be published in an address accessible to Grafana.
Another way to add icons is to copy them inside the Grafana container inside a certain directory, but we do
not officially support this method in monitoring-operator.

At the moment, adding animated SVGs is also difficult, since Grafana does not provide tools in the UI for this.
The only way to add an animated icon now is to modify the list of tools in `Experimental element types` by manipulating
the Grafana source code, however monitoring-operator does not support this method.

You can find this and other useful information about process of Canvas' improvement in
[this discussion](https://github.com/grafana/grafana/discussions/56835).

However, `Canvas` already offers a large set of tools, such as `Server` for a schematic display of your server rack
or data center, or the interactive `Button` element (Button and some other tools can be found in the list of elements
if you enabled the `Experimental element types` option in the panel settings).

### Multistat

**Plugin ID:** michaeldmoore-multistat-panel

**Ref:** <https://github.com/michaeldmoore/michaeldmoore-multistat-panel>

**Alternatives:**

* Bar chart (core)
* Bar gauge (core)

**Migration path:**

The closest type of panel to `Multistat` is `Bar chart`.

`Multistat` has a lot of settings, so one-to-one match for all `Bar chart` options would be difficult.
Let's take a look at the most important settings that can be moved from the old panel to the new one:

* `Label col` in the `Columns` section -> `X Axis` in the `Bar chart` section
* `Aggregation` in the `Columns` section -> there is no option allows calculation in general, but you can add several
  results of calculations into the legend using `Values` in the `Legend` section
* `Group col` in the `Grouping` section -> you can stack values by using `Stacking` option in the `Bar chart` section,
  but usually grouping of values occurs automatically
* `Horizontal` in the `Layout` section -> `Orientation` in the `Bar chart` section
* `Legend` in the `Layout` section -> `Visibility` in the `Legend` section
* Any coloring settings -> `Color scheme` in the `Standard options` section, color of values on Y axis can be changed
  by switching `Color` in the `Axis` section to `Series`. Also, you can configure `Thresholds`
* `Bar links` -> `Data links`

Almost all settings in the new panel can be configured both in general and for each metric independently by using
overrides.

### Discrete

**Plugin ID:** natel-discrete-panel

**Ref:** <https://github.com/NatelEnergy/grafana-discrete-panel>

**Alternatives:**

* State timeline (core)

**Migration path:**

If you simply change the type of your panel from the old `Discrete` to the new `State timeline`, most of the options
will migrate automatically. But there are some specific settings that require manual steps after this.

Color mappings and value mappings in `Discrete` can be customized independently and make some sort of pipeline.
You can match value 1 to "OK" phrase and then match "OK" to green color. But `State timeline` has a different system
of mappings, where values and colors combined into the one table, where mappings work at the same time and not as the
pipeline. That means that automatically created mappings won't work as expected.

To fix this problem, you should open `Value mappings` option in the new panel and set colors for each value-display text
mapping. Then you can remove useless value-color mappings. Range mappings can be configured here too.

Matched values can be displayed on the panel a little bit incorrectly (e.g. text may extend beyond the colored row).
In this case, you can increase `Line width` option or change `Align values`.

### Cal-HeatMap

**Plugin ID:** neocat-cal-heatmap-panel

**Ref:** <https://github.com/NeoCat/grafana-cal-heatmap-panel>

**Alternatives:**

* Heatmap (core)

**Migration path:**

There is no automated migration between `Cal-HeatMap` and `Heatmap`, but these types of panel are very similar.
If your data from query is suitable for `Cal-HeatMap`, it should be displayed well when you simply change the type
of the panel to `Heatmap`.

After that you can customize the color scheme by changing settings in the `Colors` section, cell parameters can be
changed in the `Cell display` section. The `Link template` feature of the old panel can be reproduced by using
`Data links` option on the new one.

`Domain` and `Sub domain` will be configured automatically based on the time range in `Heatmap`, there's no other option
to configure them manually.
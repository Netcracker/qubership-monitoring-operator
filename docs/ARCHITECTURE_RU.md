# Архитектура Qubership Monitoring Operator

## Обзор

**Qubership Monitoring Operator** - это Kubernetes оператор, который автоматизирует развертывание и управление полным стэком мониторинга в Kubernetes кластере. Оператор упрощает управление сложной инфраструктурой мониторинга, предоставляя единый интерфейс через Custom Resource Definition (CRD) `PlatformMonitoring`.

## Принцип работы

Оператор работает по паттерну **Kubernetes Operator Pattern**:

1. **Создание CRD** - Пользователь создает ресурс `PlatformMonitoring`, описывающий желаемое состояние мониторинга
2. **Watching** - Оператор отслеживает изменения ресурса `PlatformMonitoring`
3. **Reconciliation** - Оператор сравнивает желаемое состояние с фактическим и приводит систему к желаемому состоянию
4. **Результат** - Все компоненты мониторинга развернуты и настроены

## Основная архитектура

```
┌─────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                        │
│                                                              │
│  ┌────────────────────────────────────────────────────┐    │
│  │         Qubership Monitoring Operator              │    │
│  │  (Deployment в namespace monitoring)                │    │
│  └──────────────────┬─────────────────────────────────┘    │
│                     │ Watch                                 │
│                     ▼                                       │
│  ┌────────────────────────────────────────────────────┐    │
│  │         PlatformMonitoring CR                      │    │
│  │  (Custom Resource - желаемое состояние)            │    │
│  └────────────────────────────────────────────────────┘    │
│                                                              │
│  ┌────────────────────────────────────────────────────┐    │
│  │          Развернутые компоненты                     │    │
│  │  • Prometheus/VictoriaMetrics (БД метрик)          │    │
│  │  • Grafana (визуализация)                          │    │
│  │  • AlertManager/VMAlert (алертинг)                 │    │
│  │  • Exporters (сбор метрик)                         │    │
│  │  • Operators (управление компонентами)             │    │
│  └────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

## Что устанавливает оператор

### 1. Операторы (Orchestrators)

Эти компоненты управляют другими компонентами через CRD:

- **Prometheus Operator** - управляет Prometheus, AlertManager, ServiceMonitors, PodMonitors, PrometheusRules
- **VictoriaMetrics Operator** - управляет компонентами VictoriaMetrics (VMSingle, VMAgent, VMAlert, etc.)
- **Grafana Operator** - управляет Grafana инстансами, дашбордами и источниками данных

### 2. Базы данных временных рядов (TSDB)

**Вариант 1: Prometheus Stack**
- Prometheus Server - собирает и хранит метрики
- AlertManager - обработка и отправка алертов

**Вариант 2: VictoriaMetrics Stack** (по умолчанию, более эффективный)
- VMSingle - одиночный инстанс для хранения метрик
- VMAgent - агент для сбора метрик (альтернатива Prometheus)
- VMAlert - система алертинга
- VMAlertManager - менеджер алертов
- VMAuth - аутентификация для VictoriaMetrics
- VMUser - пользователи для VictoriaMetrics

### 3. Визуализация

- **Grafana** - веб-интерфейс для визуализации метрик
- Предустановленные дашборды для Kubernetes и приложений

### 4. Экспортеры метрик

**Инфраструктура:**
- **node-exporter** - метрики с физических нод (CPU, память, диск, сеть)
- **kube-state-metrics** - метрики состояния Kubernetes ресурсов (поды, сервисы, деплойменты и т.д.)

**Кластерные метрики:**
- ServiceMonitors для Kubernetes компонентов (API Server, Scheduler, Controller Manager, Kubelet, etcd, CoreDNS, nginx-ingress)

**Опциональные экспортеры:**
- blackbox-exporter - мониторинг доступности endpoints
- cert-exporter - мониторинг TLS сертификатов
- json-exporter - сбор метрик из REST API
- version-exporter - версии приложений
- cloud-events-exporter - события из облака

### 5. Интеграции с облачными провайдерами

- **AWS CloudWatch Exporter** - метрики из AWS
- **Azure Monitor** (Promitor Agent) - метрики из Azure
- **Google Cloud Operations** (Stackdriver Exporter) - метрики из GCP

### 6. Дополнительные компоненты

- **Pushgateway** - прием метрик через push (для краткосрочных задач)
- **Promxy** - прокси для высокой доступности Prometheus
- **Graphite Remote Adapter** - интеграция с Graphite
- **Prometheus Adapter** - адаптер для HPA (Horizontal Pod Autoscaling)

## Процесс развертывания

### Шаг 1: Установка оператора

Оператор устанавливается через Helm chart:

```bash
helm install monitoring-operator charts/qubership-monitoring-operator \
  --namespace monitoring --create-namespace
```

Это создает:
- Deployment с оператором
- ServiceAccount, Role, ClusterRole для оператора
- CRD для `PlatformMonitoring`

### Шаг 2: Создание PlatformMonitoring ресурса

Пользователь создает ресурс `PlatformMonitoring`:

```yaml
apiVersion: monitoring.qubership.org/v1alpha1
kind: PlatformMonitoring
metadata:
  name: monitoring-stack
  namespace: monitoring
spec:
  # Конфигурация компонентов
```

### Шаг 3: Reconciliation цикл

Оператор выполняет reconciliation в следующем порядке:

```
1. Prometheus Operator
   └─> Создает CRDs для Prometheus, ServiceMonitor, PodMonitor, etc.
   
2. etcd Monitor
   └─> Создает ServiceMonitor для etcd (если доступен)
   
3. Kubernetes Monitors
   └─> Создает ServiceMonitors/PodMonitors для компонентов кластера
   
4. VictoriaMetrics Operator
   └─> Создает CRDs для VictoriaMetrics компонентов
   
5. VMSingle / VMCluster
   └─> Развертывает базу данных для метрик
   
6. VMUser
   └─> Создает пользователей для доступа к VictoriaMetrics
   
7. VMAgent
   └─> Развертывает агент сбора метрик
   
8. VMAuth
   └─> Настраивает аутентификацию
   
9. Prometheus (если включен)
   └─> Развертывает Prometheus через Prometheus Operator
   
10. VMAlertManager / AlertManager
    └─> Развертывает систему алертинга
    
11. VMAlert
    └─> Настраивает правила алертинга
    
12. Exporters
    ├─> kube-state-metrics
    └─> node-exporter
    
13. Grafana Operator
    └─> Создает CRDs для Grafana
    
14. Grafana
    └─> Развертывает Grafana через Grafana Operator
    
15. Prometheus Rules
    └─> Применяет правила алертинга
    
16. Pushgateway (если включен)
    └─> Развертывает Pushgateway
```

## Оркестрация компонентов

### Уровень 1: Qubership Monitoring Operator

**Управляет:**
- Все компоненты мониторинга через `PlatformMonitoring` CR
- Координация между различными операторами
- Создание базовых ресурсов (ServiceAccounts, Services, etc.)

**Не управляет напрямую:**
- Внутренние компоненты Prometheus/VictoriaMetrics/Grafana (это делают их операторы)

### Уровень 2: Prometheus Operator

**Управляет:**
- Deployment Prometheus
- Deployment AlertManager
- Конфигурация scraping через ServiceMonitor/PodMonitor CRs
- Правила алертинга через PrometheusRule CRs

### Уровень 3: VictoriaMetrics Operator

**Управляет:**
- VMSingle/VMCluster (хранение метрик)
- VMAgent (сбор метрик)
- VMAlert (алертинг)
- VMAlertManager (управление алертами)
- VMAuth (аутентификация)
- VMUser (пользователи)

### Уровень 4: Grafana Operator

**Управляет:**
- Deployment Grafana
- GrafanaDataSource CRs (источники данных)
- GrafanaDashboard CRs (дашборды)

## Настройка

### Способ 1: Через PlatformMonitoring CR

Основной способ - редактирование `PlatformMonitoring` ресурса:

```yaml
spec:
  victoriametrics:
    vmOperator:
      install: true
    vmSingle:
      install: true
      retentionPeriod: "7d"
  grafana:
    install: true
  kubeStateMetrics:
    install: true
  nodeExporter:
    install: true
```

### Способ 2: Через Helm values

При установке через Helm можно передать конфигурацию через values.yaml:

```bash
helm install monitoring-operator charts/qubership-monitoring-operator \
  -f custom-values.yaml
```

Helm chart создаст `PlatformMonitoring` ресурс с указанными параметрами.

### Что настраивается

- **Установка/удаление компонентов** - через флаг `install: true/false`
- **Образы контейнеров** - кастомные версии
- **Ресурсы** - CPU/память для каждого компонента
- **Хранилище** - PersistentVolumeClaims для данных
- **Аутентификация** - OAuth, LDAP, Basic Auth
- **Ingress** - внешний доступ к компонентам
- **Сетевые настройки** - порты, TLS
- **Secrets/ConfigMaps** - конфигурационные данные

## Схема потока данных

```
┌─────────────────────────────────────────────────────────────┐
│                    Источники метрик                          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Node       │  │ Kubernetes   │  │ Applications │     │
│  │  Exporter    │  │  Components  │  │   Services   │     │
│  │              │  │              │  │              │     │
│  │ /metrics     │  │ /metrics     │  │ /metrics     │     │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘     │
│         │                  │                  │             │
│         └──────────────────┼──────────────────┘             │
│                            │                                │
│                            ▼                                │
│  ┌────────────────────────────────────────────────────┐    │
│  │         VMAgent / Prometheus                       │    │
│  │         (Сбор метрик)                              │    │
│  └──────────────────┬─────────────────────────────────┘    │
│                     │                                        │
│                     ▼                                        │
│  ┌────────────────────────────────────────────────────┐    │
│  │      VMSingle / Prometheus Server                  │    │
│  │      (Хранение метрик)                             │    │
│  └──────────────────┬─────────────────────────────────┘    │
│                     │                                        │
│         ┌───────────┼───────────┐                          │
│         ▼           ▼           ▼                          │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐                  │
│  │ Grafana  │ │ VMAlert  │ │ Cloud    │                  │
│  │          │ │          │ │ Export   │                  │
│  │ Визуализ │ │ Алертинг │ │          │                  │
│  └──────────┘ └──────────┘ └──────────┘                  │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## Особенности архитектуры

### 1. Многоуровневая оркестрация

- **Qubership Operator** управляет общим стэком
- **Промetheus/VictoriaMetrics/Grafana Operators** управляют своими компонентами
- Каждый уровень независим и может работать отдельно

### 2. Выбор стэка

- Можно выбрать **Prometheus** или **VictoriaMetrics**
- По умолчанию - **VictoriaMetrics** (более эффективный)
- Можно использовать оба одновременно

### 3. Автоматическое обнаружение

- Оператор автоматически создает ServiceMonitors для компонентов Kubernetes
- Grafana автоматически обнаруживает дашборды через GrafanaDashboard CRs
- Правила алертинга применяются автоматически

### 4. Расширяемость

- Пользователи могут создавать свои ServiceMonitors/PodMonitors для своих приложений
- Можно добавлять кастомные GrafanaDashboard CRs
- Поддержка кастомных PrometheusRule для алертов

## По умолчанию устанавливается

✅ **Включено:**
- VictoriaMetrics Operator
- VMSingle (24 часа retention)
- Grafana
- Grafana Operator
- kube-state-metrics
- node-exporter
- Common Dashboards
- Prometheus Rules

❌ **Выключено:**
- AlertManager (используется VMAlert)
- Все облачные экспортеры (AWS, Azure, GCP)
- Опциональные экспортеры (blackbox, cert, json и т.д.)
- Prometheus Adapter для HPA
- Интеграции (Graphite, Promxy)
- Pushgateway

## Пространство имен

Все компоненты устанавливаются в namespace, указанный в `PlatformMonitoring` ресурсе (по умолчанию `monitoring`).

Оператор может работать в том же или другом namespace, но он должен иметь доступ к namespace, где находятся компоненты мониторинга.

## Итоговая схема архитектуры

```
                    ┌─────────────────────────┐
                    │  Helm Chart Deployment  │
                    └──────────┬──────────────┘
                               │
                    ┌──────────▼──────────────┐
                    │  Monitoring Operator    │
                    │  (Deployment)           │
                    └──────────┬──────────────┘
                               │ Watches
                    ┌──────────▼──────────────┐
                    │ PlatformMonitoring CR   │
                    └──────────┬──────────────┘
                               │ Orchestrates
        ┌──────────────────────┼──────────────────────┐
        │                      │                      │
        ▼                      ▼                      ▼
┌───────────────┐    ┌───────────────┐    ┌───────────────┐
│ Prometheus    │    │ VictoriaMetrics│    │ Grafana       │
│ Operator      │    │ Operator       │    │ Operator      │
│               │    │                │    │               │
│ • Prometheus  │    │ • VMSingle     │    │ • Grafana     │
│ • AlertMgr    │    │ • VMAgent      │    │ • Dashboards  │
│ • ServiceMon  │    │ • VMAlert      │    │ • DataSources │
└───────┬───────┘    └───────┬────────┘    └───────┬───────┘
        │                    │                      │
        └────────────────────┼──────────────────────┘
                             │
                ┌────────────▼────────────┐
                │  Monitoring Components  │
                │                         │
                │ • TSDB (VM/Prom)        │
                │ • Grafana UI            │
                │ • Alerting              │
                │ • Exporters             │
                └─────────────────────────┘
```

## Ключевые моменты

1. **Один CR для всего** - `PlatformMonitoring` описывает весь стэк мониторинга
2. **Автоматизация** - оператор сам создает все необходимые ресурсы
3. **Гибкость** - можно включать/выключать компоненты через флаги
4. **Production-ready** - предустановленные конфигурации для продакшена
5. **Масштабируемость** - поддержка как небольших, так и крупных кластеров
6. **Стандарты** - использует стандартные Kubernetes ресурсы и CRDs

---

**Примечание:** Это high-level описание архитектуры. Для деталей конкретных компонентов и их настройки смотрите документацию в директории `docs/`.


@helm
@k8s
@filebeat
Feature: The Helm chart is following product recommended configuration for Kubernetes

Scenario: The Filebeat chart will create recommended K8S resources
  Given a cluster is running
  When the "filebeat" Elastic's helm chart is installed
  Then a pod will be deployed on each node of the cluster by a DaemonSet
    And a "ConfigMap" resource contains the "filebeat.yml" content
    And a "ServiceAccount" resource manages RBAC
    And a "ClusterRole" resource manages RBAC
    And a "ClusterRoleBinding" resource manages RBAC
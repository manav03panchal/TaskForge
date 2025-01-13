terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.24.0"
    }
  }
}

provider "kubernetes" {
  config_path = "~/.kube/config"
}

resource "kubernetes_namespace" "taskforge" {
  metadata {
    name = var.namespace
  }
}

resource "kubernetes_deployment" "zookeeper" {
  metadata {
    name      = "zookeeper"
    namespace = kubernetes_namespace.taskforge.metadata[0].name
  }

  spec {
    replicas = 1  # Single node for dev

    selector {
      match_labels = {
        app = "zookeeper"
      }
    }

    template {
      metadata {
        labels = {
          app = "zookeeper"
        }
      }

      spec {
        container {
          name  = "zookeeper"
          image = "bitnami/zookeeper:latest"
          
          port {
            container_port = 2181
          }

          env {
            name  = "ALLOW_ANONYMOUS_LOGIN"
            value = "yes"  # For development only
          }
        }
      }
    }
  }
}

# Zookeeper Service
resource "kubernetes_service" "zookeeper" {
  metadata {
    name      = "zookeeper"
    namespace = kubernetes_namespace.taskforge.metadata[0].name
  }

  spec {
    selector = {
      app = "zookeeper"
    }

    port {
      port        = 2181
      target_port = 2181
    }
  }
}


resource "kubernetes_deployment" "kafka" {
  metadata {
    name      = "kafka"
    namespace = kubernetes_namespace.taskforge.metadata[0].name
  }

  spec {
    replicas = 1  # Single node for development

    selector {
      match_labels = {
        app = "kafka"
      }
    }

    template {
      metadata {
        labels = {
          app = "kafka"
        }
      }

      spec {
        container {
          name  = "kafka"
          image = "bitnami/kafka:latest"

          port {
            container_port = 9092
            name          = "kafka"
          }

          env {
            name  = "KAFKA_CFG_ZOOKEEPER_CONNECT"
            value = "zookeeper:2181"
          }

          env {
            name  = "ALLOW_PLAINTEXT_LISTENER"
            value = "yes"
          }

          env {
            name  = "KAFKA_CFG_LISTENERS"
            value = "PLAINTEXT://:9092"
          }

          env {
            name  = "KAFKA_CFG_ADVERTISED_LISTENERS"
            value = "PLAINTEXT://kafka:9092"
          }

          resources {
            limits = {
              memory = "1Gi"
              cpu    = "500m"
            }
            requests = {
              memory = "512Mi"
              cpu    = "250m"
            }
          }
        }
      }
    }
  }
}

# Kafka Service
resource "kubernetes_service" "kafka" {
  metadata {
    name      = "kafka"
    namespace = kubernetes_namespace.taskforge.metadata[0].name
  }

  spec {
    selector = {
      app = "kafka"
    }

    port {
      port        = 9092
      target_port = 9092
    }
  }
}


resource "kubernetes_deployment" "redis" {
  metadata {
    name      = "redis"
    namespace = kubernetes_namespace.taskforge.metadata[0].name
  }

  spec {
    replicas = 1  # Single replica for dev

    selector {
      match_labels = {
        app = "redis"
      }
    }

    template {
      metadata {
        labels = {
          app = "redis"
        }
      }

      spec {
        container {
          name  = "redis"
          image = "redis:latest"

          port {
            container_port = 6379
          }

          resources {
            limits = {
              memory = "512Mi"
              cpu    = "200m"
            }
            requests = {
              memory = "256Mi"
              cpu    = "100m"
            }
          }
        }
      }
    }
  }
}

# Redis Service
resource "kubernetes_service" "redis" {
  metadata {
    name      = "redis"
    namespace = kubernetes_namespace.taskforge.metadata[0].name
  }

  spec {
    selector = {
      app = "redis"
    }

    port {
      port        = 6379
      target_port = 6379
    }
  }
}

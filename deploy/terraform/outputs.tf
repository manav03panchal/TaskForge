output "namespace" {
  value = kubernetes_namespace.taskforge.metadata[0].name
}

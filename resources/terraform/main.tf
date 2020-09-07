resource "null_resource" "e1" {
  provisioner "local-exec" {
    command = "date"
  }
}

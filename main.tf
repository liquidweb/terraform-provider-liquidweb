variable "storm_config_path" {
  type = "string"
}

provider "storm" {
  config_path = "${var.storm_config_path}"
}
  
resource "storm_server" "test_server" {
  config_id = 123
  image_id = 114
  domain = "terraform.dev"
  public_ssh_key = "${file("./devkey.pub")}"
  zone = 0
}

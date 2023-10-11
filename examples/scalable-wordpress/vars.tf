
# recommended changing with file `./terraform.tfvars` with content
# top_domain = "example.com"
# site_domain = "wordpress.example.com"

variable "site_name" {
  type = string
  default = "scaling.hostbaitor.com"
}

variable "top_domain" {
  type = string
  default = "hostbaitor.com"
}

variable "username" {
  type = string
  default = "wordpress"
}

variable "wordpress_dbname" {
  type = string
  default = "wordpress"
}

variable "wordpress_dbuser" {
  type = string
  default = "wordpress"
}

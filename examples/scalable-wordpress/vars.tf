
# recommended changing with file `./terraform.tfvars` with content
# top_domain = "example.com"
# site_domain = "wordpress.example.com"

variable "site_name" {
  type = string
  default = "scaling.example.com"
}

variable "top_domain" {
  type = string
  default = "example.com"
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

variable "virl_user_name" {
 description = "user name for Virl Apiu"
 default = "guest"
}

variable "virl_password" {
 description = "password for Virl Apiu"
 default = "guest"
}


resource "virl_server" "dev" {
    address = "10.204.106.107"
    port="19399"
    user_name="${var.virl_user_name}"
    password="${var.virl_password}"
}
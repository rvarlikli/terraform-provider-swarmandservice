variable "virl_user_name" {
 description = "user name for Virl Api"
 default = "guest"
}

variable "virl_password" {
 description = "password for Virl Api"
 default = "guest"
}


resource "virl_server" "dev" {
    address = "10.204.106.107"
    port="19399"
    user_name="${var.virl_user_name}"
    password="${var.virl_password}"
    virl_file="example.virl"
 	simulation_name="kadir0002"
}
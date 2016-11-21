provider "terraform-provider-mocker" {

}

variable "virl_user_name" {
  description = "user name for Virl Api"
  default = "guest"
}

variable "virl_password" {
  description = "password for Virl Api"
  default = "guest"
}


resource "mocker_server" "dev" {
  api_address = "192.168.99.100"
  port="19399"
  image_name="alpine"
  service_name="alpinex"
}
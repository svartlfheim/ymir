module "mys3bucket" {
  source  = "ymir.local/akirk/tf-ymir-mod-1/aws"
  version = "1.0.0"

  name = "mybucket"
}

# module "mys3bucket" {
#   source = "git::https://github.com/svartlfheim/tf-ymir-mod-1.git"

#   name = "mybucket"
# }

module "mydynamodb" {
  source = "git::https://github.com/svartlfheim/tf-ymir-mod-2.git"

  name     = "mydynamotable"
  hash_key = "some_key"
}

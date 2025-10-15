resource "distribution_permission_target" "my_permission" {
  name        = "my-permission"
  resource_type = "destination"
  distribution_destinations = [
    {
      site_name     = "*"
      city_name     = "*"
      country_codes = ["*"]
    }
  ]

  principals = {
    users = {
      "test1" = ["x","d"]
    }
    groups = {
      "grp1" = ["x","d"]
    }
  }
}
---
page_title: "Scaleway: scaleway_baremetal_option"
description: |-
Gets information about a baremetal option.
---

# scaleway_baremetal_option

Gets information about a baremetal option.
For more information, see [the documentation](https://developers.scaleway.com/en/products/baremetal/api).

## Example Usage

```hcl
# Get info by option name 
data "scaleway_baremetal_option" "by_name" {
  name = "Remote Access"
}

# Get info by option id
data "scaleway_baremetal_option" "by_id" {
  option_id = "931df052-d713-4674-8b58-96a63244c8e2"
}
```

## Argument Reference

- `name` - (Optional) The option name. Only one of `name` and `option_id` should be specified.
- `option_id` - (Optional) The option id. Only one of `name` and `option_id` should be specified.
- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the option exists.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the option.
- `name` - The name of the option.
- `manageable` - Is false if the option could not be added or removed.
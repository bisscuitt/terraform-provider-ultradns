---
subcategory: "HTTP-Probe"
layout: "ultradns"
page_title: "ULTRADNS: ultradns_probe_http"
description: |-
  Manages the HTTP probe records in UltraDNS.
---

# Resource: ultradns_probe_http

Use this resource to manage the HTTP probe records in UltraDNS.

## Example Usage

### Create HTTP Probe for SB Pool

```terraform
resource "ultradns_probe_http" "http_sb" {
	zone_name = "example.com."
	owner_name = "sb.example.com."
	interval = "HALF_MINUTE"
	agents = ["EUROPE_WEST", "SOUTH_AMERICA", "PALO_ALTO", "NEW_YORK"]
	threshold = 4
	total_limit{
		fail = 10
	}
	transaction{
		method = "POST"
		protocol_version = "HTTP/1.0"
		url = "https://www.ultradns.com/"
		transmitted_data = "foo=bar"
		follow_redirects = true
		expected_response = "3XX"
		search_string {
			fail = "Failure"
		}
		connect_limit{
			fail = 11
		}
		run_limit{
			fail = 12
		}
	}
}
```

### Create HTTP Probe for TC Pool

```terraform
resource "ultradns_probe_http" "http_tc" {
	zone_name = "example.com."
	owner_name = "tc.example.com."
	interval = "HALF_MINUTE"
	agents = ["EUROPE_WEST", "SOUTH_AMERICA", "PALO_ALTO", "NEW_YORK"]
	threshold = 4
	total_limit{
		warning = 10 
		critical = 15 
		fail = 20
	}
	transaction{
		method = "POST"
		protocol_version = "HTTP/1.0"
		url = "https://www.ultradns.com/"
		transmitted_data = "foo=bar"
		follow_redirects = true
		expected_response = "3XX"
		search_string {
			warning = "Warning"
			critical = "Critical"
			fail = "Failure"
		}
		connect_limit{
			warning = 5 
			critical = 8 
			fail = 10
		}
		avg_connect_limit{
			warning = 5 
			critical = 8 
			fail = 10
		}
		run_limit{
			warning = 5 
			critical = 8 
			fail = 10
		}
		avg_run_limit{
			warning = 5 
			critical = 8 
			fail = 10
		}
	}
}
```

## Argument Reference

The following arguments are supported:

* `zone_name` - (Required) (String) Name of the zone.
* `owner_name` - (Required) (String) The domain name of the owner of the RRSet. Can be either a fully qualified domain name (FQDN) or a relative domain name. If provided as a FQDN, it must be contained within the zone name's FQDN.
* `interval` - (Optional) (String) Length of time between probes in minutes. Valid values are `HALF_MINUTE`, `ONE_MINUTE`, `TWO_MINUTES`, `FIVE_MINUTES`, `TEN_MINUTES`, and `FIFTEEN_MINUTES`.</br>Default set to `FIVE_MINUTES`.
* `agents` - (Required) (String List) Locations that will be used for probing. One or more values must be specified.
Valid values are `NEW_YORK`, `PALO_ALTO`, `DALLAS`, and `AMSTERDAM`.
* `threshold` - (Required) (Integer) Number of agents that must agree for a probe state to be changed.
* `pool_record` - (Optional) (String) The pool record associated with this probe. Specified when creating a record-level probe.
* `transaction` - (Required) (Block Set, Min:1) List of nested blocks describing the http requests sent for a single probe. The structure of this block is described below.
* `total_limit` - (Optional) (Block Set, Max:1) Nested block describing the total amount of time spent on all http transactions. The structure of this block follows the same structure as the <a href="#nested-limit-block-has-the-following-structure">`limit`</a> block described below.

### Nested `transaction` block has the following structure:

* `method` - (Required) (String) HTTP method. Valid values are `GET` or `POST`.
* `protocol_version` - (Required) (String) HTTP protocol version. Valid values are `HTTP/1.0`, `HTTP/1.1`, `HTTP/2`.

-> HTTP probes will only correctly work if the indicated server supports the configured HTTP protocol version, otherwise the probe will fail.

* `url` - (Required) (String) URL to probe.
* `transmitted_data` - (Optional) (String) Data to send to URL.
* `follow_redirects` - (Optional) (Boolean) Indicates whether or not to follow redirects. Default set to false.
* `expected_response` - (Optional) (String) The Expected Response code for probes to be returned as Successful. Valid values are</br>
`2XX`: Probe will pass for any code between 200-299.</br>
`3XX`: Probe will pass for any code between 300-399.</br>
`2XX|3XX`: Probe will pass for any code between 200-399.</br>
Any combination of HTTP codes between 100-599 separated by "|" </br>For example:</br>
`201|302`
* `search_string` - (Optional) (Block Set, Max:1) Nested block describing the strings need to be search on probes successful response. It does not search the status line and headers. The structure of this block is described below.
* `connect_limit` - (Optional) (Block Set, Max:1) Nested block describing how long the probe stays connected to the resource. The structure of this block follows the same structure as the <a href="#nested-limit-block-has-the-following-structure">`limit`</a> block described below.
* `avg_connect_limit` - (Optional) (Block Set, Max:1) Nested block describing the mean connect time over the five most recent probes run on each agent. Only used for Traffic Controller Pools. The structure of this block follows the same structure as the <a href="#nested-limit-block-has-the-following-structure">`limit`</a> block described below.
* `run_limit` - (Optional) (Block Set, Max:1) Nested block describing how long the probe should run. The structure of this block follows the same structure as the <a href="#nested-limit-block-has-the-following-structure">`limit`</a> block described below.
* `avg_run_limit` - (Optional) (Block Set, Max:1) Nested block describing the mean run time over the five most recent probes run on each agent. Only used for Traffic Controller Pools. The structure of this block follows the same structure as the <a href="#nested-limit-block-has-the-following-structure">`limit`</a> block described below.

### Nested `limit` block has the following structure:

* `warning` - (Optional) (Integer) Time for the HTTP transactional probe for a warning to be generated.
* `critical` - (Optional) (Integer)  Time for the HTTP transactional probe for a critical warning to be generated.
* `fail` - (Optional) (Integer)  Time for the HTTP transactional probe for the probe to fail.

-> `warning` and `critical` are only used for Traffic Controller Pools.

### Nested `search_string` block has the following structure:

* `warning` - (Optional) (String) If the probe does not find the string within the response, or does not match it as a regular expression, a warning will be generated. 
* `critical` - (Optional) (String) If the probe does not find the string within the response, or does not match it as a regular expression, a critical warning will be generated.
* `fail` - (Optional) (String) If the probe does not find the string within the response, or does not match it as a regular expression, the probe will fail.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `guid` - (Computed) (String) The internal id for this probe.


## Import

HTTP probe records can be imported by combining their `owner_name`, `zone_name`, `record_type`, and `guid`, separated by a colon.<br/>
Example : `www.example.com.:example.com.:A (1):06084A729D56C85C`.


-> For import, the `owner_name` and `zone_name` must be a FQDN, and `record_type` should have the type, as well as the corresponding number, as shown in the example below.

Example:
```terraform
$ terraform import ultradns_probe_http.example "www.example.com.:example.com.:A (1):06084A729D56C85C" 
```
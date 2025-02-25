# pagerduty-jobs

Project for jobs related to pagerduty.

![Go container packaging](https://github.com/GeoNet/pagerduty-jobs/actions/workflows/build.yml/badge.svg)

## pd-reassign-all

Replaces a no longer working python script.
Reassigns all incidents currently assigned to one user, either to another user or re-escalate to a given level.
For example:
```
# ./bin/pd-reassign-all -subdomain="pd-subdomain" -apikey="123apikeyxyz" -from-user=PD123XYZ -to-level=1
2015/08/21 15:42:41 Found from-user: Joe Bloggs, id: PD123XYZ
```

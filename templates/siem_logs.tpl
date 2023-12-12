{
  "hostname": "{{random_v_from_list "ipAddress_list"}}",
  "action": "{{randoms "allow|deny"}}",
  "l4": "tcp",
  "access_group": "{{randoms "group-1|group-2|group-3|group-4|admin"}}",
  "source": {
    "ip": "{{random_v_from_list "ipAddress_list"}}",
    "port": {{integer 10000 99999}}
  },
  "destination": {
    "ip": "{{random_v_from_list "ipAddress_list"}}",
    "port": {{integer 10000 99999}}
  }
}
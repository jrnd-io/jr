{
  "activity_id" : {{counter "activity_id" 1 1}},
  "customer_id" : {{atoi (random_v_from_list "customer_id")}},
  "activity_type" : "{{randoms "web_open|mobile_open|new_account"}}",
  "propensity_to_churn" : {{integer 0 2}}.{{integer 0 10}}{{integer 0 10}},
  "ip_address" : "{{random_v_from_list "insurrance_c_ip"}}"
}
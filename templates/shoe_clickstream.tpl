{
  "product_id": "{{random_v_from_list "shoes_id_list"}}",
  "user_id": "{{random_v_from_list "customers_id_list"}}",
  "view_time": {{integer 10 120}},
  "page_url": "https://www.acme.com/product/{{random_string 4 5}}",
  "ip": "{{ip "10.1.0.0/16"}}",
  "ts": {{counter "ts" 1609459200000 10000 }}
}
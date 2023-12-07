{{$customer_id :=  counter "customer_id" 1000 1}}{{add_v_to_list "customer_id"  (itoa $customer_id) }}
{{$insurrance_c_ip := ip "0.0.0.0/1" }}{{add_v_to_list "insurrance_c_ip"  $insurrance_c_ip }}
{
  "customer_id": {{$customer_id}},
  "first_name": "{{name}}",
  "last_name": "{{surname}}",
  "email": "{{email}}",
  "gender": "{{gender}}",
  "income": {{integer 100000 1500000}},
  "fico": {{integer 300 850}},
  "years_active":  {{integer 1 60}}
}
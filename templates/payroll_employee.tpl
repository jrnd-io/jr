{{$employee_id :=  counter "employee_id" 1000 1}}{{add_v_to_list "employee_id"  (itoa $employee_id) }}
{
  "employee_id": {{$employee_id}},
  "first_name": "{{name}}",
  "last_name": "{{surname}}",
  "age": {{integer 18 65}},
  "ssn": "{{regex "\\d{3}-\\d{2}-\\d{4}"}}",
  "hourly_rate": {{integer 8 25}},
  "gender": "{{gender}}",
  "email": "{{email}}"
}
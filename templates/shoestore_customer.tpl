{{$id:=uuid}}{{add_v_to_list "customers_id_list" $id}}{
  "id": "{{$id}}",
  "first_name": "{{name}}",
  "last_name": "{{surname}}",
  "email": "{{email}}",
  "phone_number": "{{phone}}",
  "street_address": "{{city}}, {{street}} {{building 2}}, {{zip}}",
  "state": "{{state}}",
  "zip_code": "{{zip}}",
  "country": "United States",
  "country_code": "US"
}

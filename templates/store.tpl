{{$store_id :=  counter "store_id" 1 1}}{{add_v_to_list "store_id"  (itoa $store_id) }}
{
  "store_id": {{$store_id}},
  "city": "{{city}}",
  "state": "{{state_short}}"
}
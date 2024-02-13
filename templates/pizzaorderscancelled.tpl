{{$storeId := atoi (random_v_from_list "storeId") }}
{
  "store_id": {{$storeId}},
  "store_order_id": {{counter  (print $storeId "_store_order_id") 1001 2}},
  "date":  {{integer 18000 19000}},
  "status": "cancelled"
}
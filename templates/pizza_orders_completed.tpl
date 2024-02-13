{{$storeId := atoi (random_v_from_list "storeId") }}
{
    "store_id": {{$storeId}},
    "store_order_id": {{counter  (print $storeId "_store_order_id") 1000 2}},
    "date":  {{integer 18000 19000}},
    "status" : "completed",
    "rack_time_secs" : {{integer 130 230}},
    "order_delivery_time_secs" :  {{integer 1100 2000}}
}
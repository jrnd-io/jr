{{$id:=uuid}}{{add_v_to_list "shoes_id_list" $id}}{
  "id": "{{$id}}",
  "sale_price": "{{amount 200 2000 ""}}",
  "brand": "{{from "sport_brand"}}",
  "name": "{{randoms "Pro|Cool|Soft|Air|Perf"}} {{from "cool_name"}} {{integer 1 20}}",
  "rating": {{format_float "%.2f" (floating 1 5)}}
}

{{$id:=integer 1000 9999}}{{add_v_to_list "vehicle_id" (itoa $id)  }}{
  "vehicle_id": {{$id}},
  "engine_temperature": {{integer 150 250}},
  "average_rpm": {{integer 1800 5000}} 
}
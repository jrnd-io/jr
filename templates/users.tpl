{{$userid := (print "user_" (counter "user_id" 1 1 ) )}}{{add_v_to_list "userId"  $userid }}
{
    "registertime": {{integer64 1487715775521 1519273364600}},
    "userid": "{{$userid}}",
    "regionid": "Region_{{integer 1 9}}",
    "gender": "{{randoms "MALE|FEMALE|OTHER"}}"
}
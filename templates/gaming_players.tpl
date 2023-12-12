{{$player_id :=  counter "player_id" 1000 1}}{{add_v_to_list "player_id"  (itoa $player_id) }}
{{$gamingRoomId :=  (integer 1000 5000) }}{{add_v_to_list "gaming_room_ids"  (itoa $gamingRoomId) }}
{
  "player_id": {{$player_id}},
  "player_name": "{{username (name) (surname) }}",
  "ip": "{{ip "0.0.0.0/1" }}"
}
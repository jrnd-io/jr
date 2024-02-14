{{$ip:=ip "111.2.0.0/8"}}{{add_v_to_list "ips" $ip  }}
{{$userid:=  (itoa (counter "user_id" 1 1)) }}{{add_v_to_list "userId"  $userid }}
{
  "user_id" : {{$userid}},
  "username" : "{{randoms "akatz1022|bobk_43|alison_99|k_robertz_88|Ferd88|Reeva43|Antonio_0966|ArlyneW8ter|DimitriSchenz88|Oriana_70|AbdelKable_86|Roberto_123|AlanGreta_GG66|Nathan_126|AndySims_345324|GlenAlan_23344|LukeWaters_23|BenSteins_235"}}",
  "registered_at" : {{integer 1407645330000 1502339792000}},
  "first_name" : "{{randoms "Elwyn|Curran|Hanson|Woodrow|Ferd|Reeva|Antonio|Arlyne|Dimitri|Oriana|Abdel|Greta"}}",
  "last_name" : "{{randoms "Vanyard|Vears|Garrity|Trice|Tomini|Jushcke|De Banke|Pask|Rockhill|Romagosa|Adicot|Lalonde"}}",
  "city" : "{{randoms 	"Palo Alto|San Francisco|Raleigh|London|FrankfurtNew York"}}",
  "level" : "{{randoms "Gold|Silver|Platinum"}}"
}
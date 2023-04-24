{{$NAME := name}}{{$SURNAME := surname}}{{$COMPANY := company}}{{$CITY_INDEX := random_index "city"}}{
  "guid": "{{uuid}}",
  "isActive": {{bool}},
  "balance": {{amount 100 10000 "€"}},
  "picture": "http://placehold.it/32x32",
  "age": {{integer 20 60}},
  "eyeColor": "{{randoms "blue|brown|green"}}",
  "name": "{{$NAME}} {{$SURNAME}}",
  "company": "{{$COMPANY}}",
  "email": "{{lower $NAME}}.{{lower $SURNAME}}@{{firstword (lower $COMPANY)}}.com",
  "alt_email": "{{first (lower $NAME)}}.{{lower $SURNAME}}@{{squeeze (lower $COMPANY)}}.com",
  "about": "{{lorem 20}}",
  "address": "{{city_at (atoi $CITY_INDEX)}}, {{street}} {{building 2}}, {{zip_at (atoi $CITY_INDEX)}}"
  "phone_number": "{{land_prefix_at (atoi $CITY_INDEX)}} {{regex "[0-9]{7}"}}"
  "latitude": {{latitude}},
  "longitude": {{longitude}}
}

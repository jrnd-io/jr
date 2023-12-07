{{/* TODO: This template should produce also value City_ and State_ without numbers */}}
{
  "ordertime" : {{integer 1487715775521 1519273364600}},
  "orderid" : {{counter "orderid" 0 1}},

  "itemid" : "Item_{{integer 1 99}}{{trimchars (random_string_vocabulary 1 2 "1234567890               ") " " }}",
  "orderunits" : {{integer 0 10}}.{{integer 1 99}},
  "address" : {
    "city" : "City_{{integer 1 10}}{{trimchars (random_string_vocabulary 1 2 "1234567890               ") " " }}",
    "state" : "State_{{integer 1 10}}{{trimchars (random_string_vocabulary 1 2 "1234567890               ") " " }}",
    "zipcode" : {{integer 10000 99999}}
  }
}
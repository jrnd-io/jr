{
  "guid": "{{uuid}}",
  "isActive": {{bool}},
  "balance": "{{amount 100 10000 "€"}}",
  "picture": "http://placehold.it/32x32",
  "age": {{integer 20 60}},
  "eyeColor": "{{randoms "blue|brown|green"}}",
  "name": "{{name}} {{surname}}",
  "gender": "{{gender}}",
  "company": "{{company}}",
  "email": "{{email}}",
  "about": "{{lorem 20}}",
  "address": "{{city}}, {{street}} {{building 2}}, {{zip}}",
  "phone_number": "{{land_prefix}} {{regex "[0-9]{7}"}}",
  "latitude": {{latitude}},
  "longitude": {{longitude}}
}

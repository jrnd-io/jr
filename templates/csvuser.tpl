{
  "age": {{integer 20 60}},
  "eyeColor": "{{randoms "blue|brown|green"}}",
  "name": "{{fromcsv "NAME"}}",
  "surname": "{{fromcsv "SURNAME"}}",
  "company": "{{company}}",
  "email": "{{lower (fromcsv "NAME") }}.{{lower (fromcsv "SURNAME") }}@emeraldcity.oz"
}
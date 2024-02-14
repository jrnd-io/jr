{
  "ip" : "{{random_v_from_list "ips"}}",
  "userid" : {{atoi (random_v_from_list "userId")}},
  "remote_user" : "-",
  "time" : "{{integer 1 10}}",
  "_time" : {{integer 1 10}},
  "request" : "{{randoms "GET /index.html HTTP/1.1|GET /site/user_status.html HTTP/1.1|GET /site/login.html HTTP/1.1|GET /site/user_status.html HTTP/1.1|GET /images/track.png HTTP/1.1|GET /images/logo-small.png HTTP/1.1"}}",
  "status" : "{{randoms "200|302|404|405|406|407"}}",
  "bytes" : "{{randoms "278|1289|2048|4096|4006|4196|14096"}}",
  "referrer" : "-",
  "_logtime": {{counter "logtime" 1 10}},
  "agent" : "{{randoms "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)|Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36"}}"
}



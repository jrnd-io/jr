{{ $cardId      := itoa (counter "cardId" 1 1)}}{{add_v_to_list "cardId" $cardId  }}{{ $cardNumber  := card (randoms "visa|mastercard|amex|discover") }}{{add_v_to_list "cardNumber" $cardNumber  }}{{ $cardCVV     := cardCVV 3 }}{{ add_v_to_list "cardCVV" $cardCVV  }}
{{ $cardExpire  := (print (regex "^0[1-9]|1[012]$" )  "/" (integer 22 30) )}}{{ add_v_to_list "cardExpire" $cardExpire }}{{ $cardExpire  := (print (regex "^0[1-9]|1[012]$" )  "/" (integer 22 30) )}}{{ add_v_to_list "cardExpire" $cardExpire }}{
     "card_id": {{$cardId}},
     "card_number":  "{{$cardNumber}}",
     "cvv": "{{$cardCVV}}",
     "expiration_date": "{{$cardExpire}}"
}
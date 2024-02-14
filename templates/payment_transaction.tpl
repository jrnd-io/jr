{{$index:= integer 1 30 }}
{
    "transaction_id": {{counter "transaction_id" 1 1}},
    "card_id": {{atoi (get_v_from_list_at_index "cardId" $index)}},
    "user_id": "{{get_v_from_list_at_index "userId" $index}}",
    "purchase_id": {{counter "purchase_id" 0 1 }},
    "store_id":  {{integer 1 7 }}
}
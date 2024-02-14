{{$productId :=counter "productId_list" 1 1 }}{{add_v_to_list "productId_list" (itoa $productId) }}
{
  "id": {{$productId}},
  "name": "Product-{{counter "name" 0 1}}",
  "description": "Item n. {{counter "description" 0 1}}",
  "price": {{integer 1 100}}
}
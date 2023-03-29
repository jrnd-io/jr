jr run  --template '{{regex "123[0-2]+.*\w{3}"}}'



echo '123[0-2]+.*\w{3}' | { printf "%s" "Enter regex: "; read regex ; jr run  --template "{{regex \"$regex\"}}" ; }

MY_REGEX='123[0-2]+.*\w{3}' ;  jr run  --template "{{regex \"${MY_REGEX}\"}}" 

echo '123[0-2]+.*\w{3}' | { read regex; jr run  --template "{{regex \"$regex\"}}" ; }
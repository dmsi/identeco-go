#TOKEN=eyJhbGciOiJSUzI1NiIsImtpZCI6IjVlOWU4YWJmYWZhZDg4ZTRkMThhY2FmMTNkNjJhZDc5IiwidHlwIjoiSldUIn0.eyJleHAiOjE2ODg0OTI4NjQsImlzcyI6Imh0dHBzOi8vZ2l0aHViLmNvbS9kbXNpL2lkZW50ZWNvIiwidG9rZW5fdXNlIjoiYWNjZXNzIiwidXNlcm5hbWUiOiJib3NzIn0.CRNRq0F__NddCf_9VXkmYjZOl_UMpVCaoXrQgIvvZ7LTthmttgG7GoVkwvluxQHcmJCjQOKgKy0E0C3i7wGe0wIlRqHbyMvxE3iwveQdy9IJY7jj-sdzjRcWerP95dtH06rDYT7VVmo6FVLnEyT63LpCBnMdX9EwL-zRPEfH1l3ZoVX3M2AfyIoeKv3AMWR87cf0FGyQgRwnz51Ecz5ZXXaDG0QmtDAB1-BSetO9TP5nXnKdG6mVjYRTjBET08cwSzCI3BD1IJYjGD_lPY-M0YimkcYmYVtzJl4x6XcLmwblqNtKh-0uLH7FVmmGOgucrHMzZDhB7Gsgwi2BWs8OIw
TOKEN=eyJhbGciOiJSUzI1NiIsImtpZCI6IjVlOWU4YWJmYWZhZDg4ZTRkMThhY2FmMTNkNjJhZDc5IiwidHlwIjoiSldUIn0.eyJleHAiOjE2OTEwNzg5OTMsImlzcyI6Imh0dHBzOi8vZ2l0aHViLmNvbS9kbXNpL2lkZW50ZWNvIiwidG9rZW5fdXNlIjoicmVmcmVzaCIsInVzZXJuYW1lIjoiYm9zcyJ9.EurjtrTfjjYrMWdRWLgUPCLMbdARKgc5n7YRNTEMB_HUFLA96lB8jGBOZNe_s8SqX_246_HaVqYvY-xuI8fLMk65KT3liK-Fj6Q9dF9T4Lbd19tDCUL75Svf0jMs5TVifo6Z5qlPGxPSI_pvjnnjdhQPqTSLbb8bKTK8WxkbctHQ7j3NygZ0yLmFV9HmNtBxDLnCV_HI7f9jnNjAk4DnZkzdULcqhkysu4E5bcBtGh1paT1XKG56LZUk5aBgDW27QNut2HXYeDlpguYk7GFG_InxwrIFxewwFqHsMv7DBBrxSUmhuQPQ_hsWa-0QbYwJB1KP6RsvUnKEukVVk1ZD5A
#https://u10mhwvp7i.execute-api.eu-west-1.amazonaws.com/dev/.well-known/jwks.json
#aws-ido --region eu-west-1 dynamodb get-item --table-name identeco-dev-users --key '{"username":{"S":"boss"}}'
ENDPOINT=u10mhwvp7i.execute-api.eu-west-1.amazonaws.com/dev
#curl -i -X GET https://$ENDPOINT/.well-known/jwks.json
#curl -i -X POST -H "Content-Type: application/json" https://$ENDPOINT/register -d '{"username":"boss", "password":"boss"}'
#curl -i -X GET -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" https://$ENDPOINT/refresh
curl -i -X POST -H "Content-Type: application/json" https://$ENDPOINT/login -d '{"username":"boss", "password":"boss"}'

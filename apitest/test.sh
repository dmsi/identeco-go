#aws-ido --region eu-west-1 dynamodb get-item --table-name identeco-dev-users --key '{"username":{"S":"boss"}}'
ENDPOINT=nqf9s1b2zi.execute-api.eu-west-1.amazonaws.com/dev
curl -i -X POST -H "Content-Type: application/json" https://$ENDPOINT/register -d '{"username":"boss", "password":"boss"}'

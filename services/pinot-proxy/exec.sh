# set env vars

# set credentials
aws configure set aws_access_key_id "$TEMP_PINOT_AWS_ACCESS_KEY_ID"
aws configure set aws_secret_access_key "$TEMP_PINOT_AWS_SECRET_ACCESS_KEY"
aws configure set region "$TEMP_PINOT_AWS_REGION"

# update kubeconfig
aws eks --region "$TEMP_PINOT_AWS_REGION" update-kubeconfig --name pinot-on-eks

# run portforward
AWS_SECRET_ACCESS_KEY=$TEMP_PINOT_AWS_SECRET_ACCESS_KEY AWS_ACCESS_KEY_ID=$TEMP_PINOT_AWS_ACCESS_KEY_ID AWS_REGION=$TEMP_PINOT_AWS_REGION kubectl port-forward --address=0.0.0.0 service/pinot-broker 8099:8099 -n pinot
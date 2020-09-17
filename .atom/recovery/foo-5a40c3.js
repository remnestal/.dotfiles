tests["Status is OK"] = pm.response.to.be.ok;

// Response content validation
var data = JSON.parse(responseBody);
tests["ID is correct"] = data.id === 3;
tests["Name is correct"] = data.first_name === "Anton";
tests["Name is correct"] = data.last_name === "Janer";
tests["StripeID is set"] = data.stripe_id === "cus_FrXnwKw6CcWA2w";
tests["Payment status is not verified"] = data.has_payment === false;

// Schema check
tests["JSON response schema validation"] = tv4.validate(data, {
  "required": ["id", "name", "stripe_id", "has_payment"],
  "properties": {
    "id": {
      "type": "integer"
    },
    "name": {
      "type": "string"
    },
    "stripe_id": {
      "anyOf": [
        { "type": "string" },
        { "type": "null" },
      ]
    },
    "has_payment": {
      "type": "boolean"
    },
  }
});

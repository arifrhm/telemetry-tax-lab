package testutil

import _ "embed"

var SamplePayload = []byte(`{
  "event_id": "evt-tax-01a2b3c4-d5e6-7890-abcd-ef0123456789",
  "tenant_id": "tenant-corp-enterprise-us-east-1",
  "timestamp": 1731724800000000000,
  "service_name": "checkout-order-dispatcher",
  "route": "/api/v1/orders/checkout",
  "client_ip": "192.168.1.105",
  "attributes": {
    "cart_size": 7,
    "currency": "USD",
    "total_amount": 349.99,
    "payment_method": "credit_card",
    "is_vip": true
  },
  "tags": ["ecommerce", "checkout", "production", "high-priority", "tier-1"]
}`)

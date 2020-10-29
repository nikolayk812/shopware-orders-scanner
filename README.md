#  Shopware orders scanner

Extensible tool to scan orders in <a href=https://www.shopware.com/en/products/shopware-6/>Shopware 6</a> online shop system. 
It finds orders in inconsistent states and generates a report.

## Motivation

An order in Shopware has multiple attributes which change during the order lifecycle:
- Payment status
- Delivery status
- Order status
- Tracking code
- Documents, i.e. delivery slip
- Custom fields, i.e. external order id, external line item ids, warehouse order id, partial cancellations, returns and refunds

![order statuses](order-statuses.png)

One could think of interdependencies between above mentioned attributes.

- If an order has delivery status: Shipped, then normally it should have a tracking code and might have a delivery slip upload as a document.
- If an order is Cancelled, then it should be refunded. (In the case it has been ever paid by a customer)

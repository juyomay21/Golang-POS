Login

✔ Functional: Succesfully get cashier list
✔ Functional: Succesfully get cashier list with limit and skip
✔ Functional: Succesfully to get passcode for login
✔ Functional: Succesfully to login with correct passcode
✔ Functional: Succesfully to get respons with Authorization header
✔ Negative case: Failed to get passcode with wrong Cashier ID
✔ Negative case: Failed to login with wrong passcode
✔ Negative case: Failed to get respon without Authorization header
✔ Functional: Succesfully to create cashier
✔ Functional: Succesfully to update cashier
✔ Negative case: Failed to create, incomplete cashier data
✔ Negative case: Failed to update, cashier data does not match
✔ Functional: Successfully to delete cashier
✔ Negative case: Failed to update, cashier ID does not exist
✔ Negative case: Failed to delete, cashier ID does not exist
List Produk

✔ Functional: Succesfully to get product list
✔ Functional: Succesfully to get product list with limit and skip
✔ Functional: Succesfully to get product List with ID category
✔ Functional: Succesfully to get product list with pocari query keyword
✔ Functional: Succesfully to create product
✔ Functional: Succesfully to update product
✔ Negative case: Failed to create, incomplete product data
✔ Negative case: Failed to update, product data does not match
✔ Functional: Succesfully to delete product
✔ Negative case: Failed to update, product ID does not exist
✔ Negative case: Failed to delete, produk ID does not exist
✔ Functional: Succesfully to get category list
✔ Functional: Succesfully to create category
✔ Functional: Succesfully to update category
✔ Negative case: Failed to create, incomplete category data
✔ Negative case: Failed to update, category data does not match
✔ Functional: Succesfully to delete category
✔ Negative case: Failed to update, category ID does not exist
✔ Negative case: Failed to delete, category ID does not exist

Confirm Order

Functional: Succesfully to get subtotal
Negative case: Failed to get subtotal, empty product

Input Payment

✔ Functional: Succesfully to get list payment method
✔ Functional: Cash card validation with grand total 5000
✔ Functional: Cash card validation with grand total 10000
Functional: Succesfully to create order
Negative case: Failed to create, empty order data
Negative case: Failed to create, empty product data

Proof of Payment

Functional: Succesfully to get order detail
Functional: Succesfully to get pdf download

Report

Functional: Succesfully to get cashier revenues today
Functional: Succesfully to get product solds today

Performance

✔ Scalability: Docker image file size optimization should be less than 300MB
Scalability: Create Product - succesfully to run 1000 request on 1 concurrency
Scalability: Product Search - succesfully to run 1000 request on 1 concurrency
Scalability: Product search - succesfully to run 1000 request on 10 concurrency
Scalability: Get Subtotal - succesfully to run 1000 request on 1 concurrency
Scalability: Create Order - succesfully to run 1000 request on 1 concurrency
Scalability: Get Search - Average response should be les than 100ms from stress test


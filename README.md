## Mobile Wallet

### databases

1. user table
	| id          | phone_no | name      | password | salt          | balance        |transter_pin|create_date|
	| ----------- | -------------- | --------- | -------- | ------------- | -------------- |--| ---|
	| primary_key | user unique id | user name | password | password salt | balance amount |transter confirm pin|
	
2. transcation table

   | id   | no | from_user | to_user | amount | request_date | confirm_date| status|
   | ---- |----| --------- | ------- | ------ | ---------------- |---|--|
   |      |random 64 characters strings| user_no   | user_no | amout  | datetime         |

​      * `status` *1*: created; *2*: confirmed

### Constants

#### transcation type defines

| value | meaning |
| ----- | ------- |
| 1 | got|
| 2 | sent|

#### status code

| value | meaning | description |
| ----- | :------ | ----------- |
| 0     | success |             |
| 1001 | user not exists| value start with 1 relates to user logins|
| 1002 | error password||
|1003|account not exists| query by other user check receiver whether exists |
|1004| need login||
| 2001 | cannot find receiver account | values starts with 2 relates to transter transcation errors|
| 2002 | balance not enough||

### APIs

1. `/user/login`

   **params:**

   * phone： user login no。mostly phone no or other unique IDs.

   * password

   **return:**

   ```json
   {
     "status": 0,  // 0 success. other error code. 
     "body": {
     	"user_info": {
       	"transcation_no": "",
     		"name": "userNname",
       	"balance": "12.30"  
   		}
     }
   }
   ```

2. `/user/get_account_intro`

   check account whether exists. If exists, return the user base info with publish format.

3. `/transaction/list`

   get transaction history

   **params:**

   * `page`: *optional*, default 1. from 1.
   * `psize`: *optional*, default 10.
   * `transcation_type`: [reference type defines](transcation type defines) *optional*, default 0.
   * `time_range_start:`  *future implementes*
   * `time_range_end:` *future implementes*
   * `user_info`: *future implementes*

   **return:**

   ```json
   {
     "status": 0, // 0 is success. other is error code.
     "body": {
       "total": 88, // total transcation count
       "page": 3, // current page
       "psize": 10, // current page size
       "transcation_list": [
          { 
           "transcation_no": "ajixiyasx", 
           "user_name": "user_name1", // transaction user
           "transcation_type": 1, // ref transaction type defines below
           "amount": 20.30  // transter amount
         },
          { 
           "transcation_no": "joijoijioa", 
           "user_name": "user_name2", // transaction user
           "transcation_type": 1, // ref transaction type defines below
           "amount": 20.30  // transter amount
         },
       ]
     }
   }
   ```

4. `/transcation/detail`

   get specify transcation

5. `/transaction/create_order`

   **params:**

   * `to_user`

   * `amount`

   **return**

   ```json
   {
    	status: 0, // success. error defines: ref error code  
     transcation_no: "jijijoaijixea" // encoded transcation id.
   }

6. `/transcation/confim`

   **params:**

   * `tranfer_password` user setted transfer pin. future can be sms or email.
   * `transcation_no`

   **return**

     ```json
   {
    	status: 0, // success. error defines: ref error code  
     transcation_no: "jijijoaijixea" // encoded transcation id.
   }
     ```

7. `/transcation/cancel`
   **params**

   ​	* `transcation_no`: 

   **return**

     ```json
   {
    	status: 0, // success. error defines: ref error code  
     transcation_no: "jijijoaijixea" // encoded transcation id.
   }
     ```


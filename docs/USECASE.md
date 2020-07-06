# List of usecase

## User Login
Title: User login <br/>
Description: Actor want to login into system <br/>
Actors:
- Customer

Input: Username/email and password<br/>
Output: JWT Token<br/>
Pre-Conditions:<br/>
- Customer already registered in system

Basic Flow:
1. Actor input username/email and password
2. Validate input:
    - Business rule: customer email must valid
    - Business rule: password not empty
3. Check user already registered
4. If user not exists return error Not Found
5. Compare password with hashed password
6. Generate JWT Token
7. Return JWT Token

Post-Conditions: -

## User Logout
Title: User logout <br/>
Description: Actor want to logout from system <br/>
Actors:
- Customer

Input: JWT Token<br/>
Output: Success or fail<br/>
Pre-Conditions:<br/>
- Token already registered in system

Basic Flow:
1. Actor provide JWT Token
3. Revoke JWT Token
4. Return succeed or failed

Post-Conditions: -

## Create User
Title: Create user<br/>
Description: Actor want to create user into system<br/>
Actors:
- Customer

Input: Username, Email, Mobile Number, Password<br/>
Pre-conditions:<br/>
- Customer not registered in system

Basic Flow:
1. Actor input username, email, mobile number, and password
2. Validate input:
    - Business rule: customer email must valid
    - Business rule: username not empty
    - Business rule: password not empty
3. Check user not registered in system
    - Business rule: email must unique
    - Business rule: username must unique
    - Business rule: mobile number must unique
4. If user already registered in system then return error Conflict
5. Generate hashed password
6. Save data into system
7. Return user data

Post-Conditions: -

## Get User
Title: Get user<br/>
Description: Actor want to get user data<br/>
Input: User id<br/>
Actor:
- Customer

Pre-conditions:
- Customer already registered in system

Basic Flow:
1. Actor provide user id
2. Check user in system by user id
3. If user not exists return error Not Found
4. Return id, username, and email

Post-Conditions: -

## Update User
Title: Update user<br/>
Description: Actor want to update user data<br/>
Input: User id, email and password
Actor:
- Customer

Pre-conditions:
- Customer already registered in system

Basic Flow:
1. Actor provide user id
2. Check user in system by user id
3. If user not exists return error Not Found
4. Update email and password into system
5. Return id, username, and email

Post-Conditions: -

## Delete User
Title: Delete user<br/>
Description: Actor want to delete user data<br/>
Input: User id<br/>
Actor:
- Customer

Pre-conditions:
- Customer already registered in system

Basic Flow:
1. Actor provide user id
2. Check user in system by user id
3. If user not exists return error Not Found
4. Delete user, balance, and balance history in system
5. Return deleted true

Post-Conditions: -

## Top up Balance
Title: Top up balance<br/>
Description: Actor want to top up balance into system<br/>
Input: User id, nominal<br/>
Actor:
- Customer

Pre-conditions:
- Customer already registered in system

Basic Flow:
1. Actor provide user id and nominal
2. Check user in system by user id
3. If user not exists return error Not Found
4. Update balance and insert history
5. Return balance

Post-Conditions: -

## Get Balance
Title: Get balance<br/>
Description: Actor want to get balance from system<br/>
Input: User id<br/>
Actor:
- Customer

Pre-conditions:
- Customer already registered in system

Basic Flow:
1. Actor provide user id 
2. Check user in system by user id
3. If user not exists return error Not Found
5. Return balance

Post-Conditions: -

## Get Balance Histories
Title: Get balance histories<br/>
Description: Actor want to get balance histories from system<br/>
Input: User id<br/>
Actor:
- Customer

Pre-conditions:
- Customer already registered in system

Basic Flow:
1. Actor provide user id 
2. Check user in system by user id
3. If user not exists return error Not Found
5. Return balance histories

Post-Conditions: -

## Transfer Balance
Title: Transfer balance<br/>
Description: Actor want to transfer balance to other user<br/>
Input: Sender user id, receiver user id, nominal<br/>
Actor:
- Customer

Pre-conditions:
- Customer already registered in system

Basic Flow:
1. Actor provide sender user id, receiver user id and nominal
2. Check user in system by user id
3. If user not exists return error Not Found
4. Reduce balance and insert history into sender account
5. Add balance and insert history into receiver account
5. Return sender balance

Post-Conditions: -

# E-Wallet Example

Table of content:
1. Description
2. System Design

---

## Description
E-Wallet is a system used to store money balances but not real money and use it to facilitate online transactions. In this repository two systems will be distinguished, E-Wallet system and Bank system. E-Wallet system will not record any transactions from bank, vice versa. Users who will make transactions using E-Wallet must register to the system and top up balance through the bank system.

![Diagram](docs/assets/diagram.png)

## System Design

### Usecase Diagram
![Diagram](docs/assets/ewallet-ewallet usecase diagram.png)
 
 [See the details](docs/USECASE.md)

### Database Design
![Diagram](docs/assets/database-design.png)
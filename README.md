Микросервис для работы с балансом пользователей.
Основное задание:
Метод начисления средств на баланс. Принимает id пользователя и сколько средств зачислить.
Метод списания средств с баланса. Принимает id пользователя и сколько средств списать.
Метод перевода средств от пользователя к пользователю. Принимает id пользователя с которого нужно списать средства, id пользователя которому должны зачислить средства, а также сумму.
Метод получения текущего баланса пользователя. Принимает id пользователя. Баланс всегда в рублях.

EndPoints list:

GET "/app/{userId}" - get user balance

GET "/app/history/{userId}" - get history of user transanctions

POST "/app/{userId}?incr= " - add funds to balance (if userId does not exist - create user and then add funds)

POST "/app/{userId}?decr= " - withdrawal of funds from the balance

POST "/app/{userId1}/{userId2}?sum= " - transfer of funds between users (userId1 - sender, userId2 - reciver)
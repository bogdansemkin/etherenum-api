# etherenum-api

`
https://nubip-app.herokuapp.com/api/v1/transactions/filter?page=pageNumber
`
Use this end-point for gettin' transactions by filter:
-transaction id
-blocknumber
-from 
-to
-timestamp

`
https://nubip-app.herokuapp.com/api/v1/transactions/?page=pageNumber
`
Use this end-point for get all transactions 

P.S: 

en:
Not implemented in this specification:
Incrementing incoming blocks, as well as outputting data on the amount and block number as a number.

About the second point: the data is flying in the form of a hash and I have no idea how to display it as a number, unless decrypted.

There is also a transfer of Service to pkg, which is an application-level error (well, this is more of a minor flaw).

It was also worth adding a swagger, but there is no time for this. Due to problems with the light, I already stretched out the delivery of the TK, and I don’t see the point in pulling further.

ru:
В данном техническом задании не реализовано:
Инкрементация входящий блоков, а так же вывод данных о сумме и номера блока как число.

О втором пункте: данные летят в виде хэша и я не представляю как его вывести числом, если только не дешифровать. 

Так же есть передача Service в pkg, что есть ошибкой уровней приложения(ну, это скорее минорный недочёт).

Так же стоило добавить сваггер, но на это нету времени. Из-за проблем я и так расстянул сдачу ТЗ, да и тянуть дальше уже как-то нечестно.

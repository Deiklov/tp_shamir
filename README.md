# test_channels
1. если данные отправили то читаем только в другой горутине
2. каналы это просто структура данных
3. данные нельзя отправить после закртия канала
3. после закрытия из канала можно прочитать все что там было, после фул чтения будем default значения принимать
те кто был заблочен из-за канала сразу отблочаться либо чтением либо паникой
4. проверка закрытия val, opened:= <-intCh
5. пустая структура не занимает память, необходима для синхронизации горутин
6. range по каналу идет пока close не придет
7. потоковое получение даннных, можно заменить на range
```go
for { 
    num, opened := <- intCh 
    if !opened { 
        break
    } 
    fmt.Println(num)
}
```
8. for + select блокируется пока не получит данные из канала любого, не грузит проц
9. если меняются данные, то используем Lock() если только читаем то RLock()
пусть нужно прочитать их мапы двум потокам и одному написать, те кто читают вызывают Rlock() тогда,
 заблочится только вызов Lock(), RLock() пройдет нормально, тем самым меньше будет очередь горутни
10. есть sync.Map() мапа для thread-safing
11. из канала можно читать в той же горутине, если буфер не заполнен
12. range считает весь буфер до закрытия канала
13. если в select есть ветвь default, то она выполняется и не идет switch context 
14. нужно получить любой ответ за отведенное время
```go
 select {
    case res := <-chan1:
        fmt.Println("Response from service 1", res, time.Since(start))
    case res := <-chan2:
        fmt.Println("Response from service 2", res, time.Since(start))
    case <-time.After(30 * time.Second):
        fmt.Println("No response received", time.Since(start))
    }
```
15. пустой select вызовет deadlock
16. можно органищовать пул воркеров
17. обычные каналы синхронные, буферезированные асинхронные


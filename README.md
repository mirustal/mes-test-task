Тестовое задание
Set запрос отправляет сообщение на сервер, сервер записывает его в базу данных и после успешной записи отправляет сообщение в брокер.

На сервере крутится второй сервис где consumer кафки обрабатывает входящие сообщение и помечает его в базе данных как прочитанное.

тестировался на машине MacOS с arm процессором

Локальный запуск: 
```
git clone https://github.com/mirustal/mes-test-task.git
cd ./mes-test-task
docker-compose up -d
```


Пример работы:

![record](https://github.com/user-attachments/assets/08ebfb53-d2ac-48fc-8bd7-d13259d42676)


Доступ к api
get запрос

![image](https://github.com/user-attachments/assets/479b4bf0-1ef9-49d5-b00b-f434f6550624)

set запрос

![image](https://github.com/user-attachments/assets/978ac039-7173-4c1e-b539-ef3b92884ed8)

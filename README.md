# **Тестовое задание**

**Описание задачи:**
* Set-запрос отправляет сообщение на сервер.
* Сервер записывает его в базу данных.
* После успешной записи отправляет сообщение в брокер.

На сервере работает второй сервис, где **Kafka consumer** обрабатывает входящие сообщения и помечает их в базе данных как прочитанные.

**Тестирование:**
Тестировался на машине MacOS с ARM процессором.

## **Локальный запуск:**

1. Клонируйте репозиторий:
    ```bash
    git clone https://github.com/mirustal/mes-test-task.git
    ```
2. Перейдите в директорию проекта:
    ```bash
    cd ./mes-test-task
    ```
3. Запустите Docker Compose:
    ```bash
    docker-compose up -d
    ```
**Документация:** [Swagger](http://http://176.109.107.236/:8081/swagger/)

## **Пример работы:**

![record](https://github.com/user-attachments/assets/08ebfb53-d2ac-48fc-8bd7-d13259d42676)


## **Доступ к API:**

**GET-запрос:**

![image](https://github.com/user-attachments/assets/479b4bf0-1ef9-49d5-b00b-f434f6550624)

**SET-запрос:**

![image](https://github.com/user-attachments/assets/978ac039-7173-4c1e-b539-ef3b92884ed8)

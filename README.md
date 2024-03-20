[Jump to English](#English)

<a name="Russain"></a>
# Русский
<p id="ru"><h3>ТЗ Lv0 на стажировке Wildberries</h3></p>
<p>Необходимо разработать демонстрационный сервис с простейшим интерфейсом, отображающий данные о заказе</p>
<p>Модель данных в формате JSON прилагается к заданию</p>
<p>Что нужно сделать:</p>
<ol>
  <li>
    Развернуть локально PostgreSQL
    <ul>
      <li>Создать свою БД</li>
      <li>Настроить своего пользователя</li>
      <li>Создать таблицы для хранения полученных данных</li>
    </ul>
  </li>
  <li>
    Разработать сервис
    <ul>
      <li>Реализовать подключение и подписку на канал в nats-streaming</li>
      <li>Полученные данные записывать в БД</li>
      <li>Реализовать кэширование полученных данных в сервисе (сохранять in memory)</li>
      <li>В случае падения сервиса необходимо восстанавливать кэш из БД</li>
      <li>Запустить http-сервер и выдавать данные по id из кэша</li>
    </ul>
  </li>
  <li>Разработать простейший интерфейс отображения полученных данных по id заказа</li>
</ol>
<p>Советы:</p>
<ul>
  <li>Данные статичны, исходя из этого подумайте насчет модели хранения в кэше и в PostgreSQL. Модель в файле model.json</li>
  <li>Подумайте как избежать проблем, связанных с тем, что в канал могут закинуть что-угодно</li>
  <li>Чтобы проверить работает ли подписка онлайн, сделайте себе отдельный скрипт, для публикации данных в канал</li>
  <li>Подумайте как не терять данные в случае ошибок или проблем с сервисом</li>
  <li>Nats-streaming разверните локально (не путать с Nats)</li>
</ul>

<hr>

[Перейти к русскому](#Russian)
<a name="English"></a>
# English

<p><h3>Specifications for Lv0 at the Wildberries internship</h3></p>
<p>It is necessary to develop a demo service with a simple interface that displays order data</p>
<p>The data model in JSON format is attached to the assignment</p>
<p>What to do:</p>
<ol>
   <li>
     Deploy PostgreSQL locally
     <ul>
       <li>Create your own database</li>
       <li>Set up your user</li>
       <li>Create tables to store the received data</li>
     </ul>
   </li>
   <li>
     Develop a service
     <ul>
       <li>Implement connection and subscription to a channel in nats-streaming</li>
       <li>Write the received data into the database</li>
       <li>Implement caching of received data in the service (save in memory)</li>
       <li>If the service fails, you need to restore the cache from the database</li>
       <li>Start the http server and output data by id from the cache</li>
     </ul>
   </li>
   <li>Develop a simple interface for displaying received data by order id</li>
</ol>
<p>Tips:</p>
<ul>
   <li>The data is static, based on this, think about the storage model in the cache and in PostgreSQL. Model in the model.json file</li>
   <li>Think about how to avoid problems associated with the fact that anything can be thrown into the channel</li>
   <li>To check if the online subscription works, make yourself a separate script to publish data to the channel</li>
   <li>Think about how not to lose data in case of errors or problems with the service</li>
   <li>Deploy Nats-streaming locally (not to be confused with Nats)</li>
</ul>

# Uyga vazifa: RabbitMQ va REST API yordamida Task Management System

## Maqsad
Ushbu uyga vazifa `RabbitMQ` va `Go` yordamida asinxron vazifalarni qayta ishlashni o'rganish va `REST API` yaratish bo'yicha tajriba orttirishdan iborat.

## Loyiha Tavsifi
Task Management System yaratish. `REST API` yordamida vazifalarni yaratish, o'qish, yangilash va o'chirish imkonini berish. `RabbitMQ` asinxron vazifalarni qayta ishlash uchun ishlatiladi.

## Loyihaning tarkibi
1. **REST API**: Vazifalarni yaratish, o'qish, yangilash va o'chirish uchun endpointlar taqdim etilsin.
2. **RabbitMQ producer**: `RabbitMQ` ga vazifalar bo'yicha xabarlarni yuboradi.
3. **RabbitMQ consumer**: `RabbitMQ` dan vazifalar bo'yicha xabarlarni qayta ishlaydi. Consumer xabarni oganidan so'ng MongoDB ga saqlaydi
4. **Vazifalarni statusini o'zgartisrish**: Vazifalarni statusini o'zgarishini qayta ishlash.

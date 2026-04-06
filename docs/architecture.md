# Proje Mimari Detayları
Bu dosya opsiyonel ek dökümantasyon için ayrılmıştır.

Genel Akış:
1. İstemci NestJS üzerinde REST veya GraphQL (Port: 3000) aracılığıyla istekte bulunur.
2. Etkinlik (Event) oluşturulurken Prisma üzerinden PostgreSQL'e kayıt atılır.
3. Arka planda NestJS asenkron olarak (Axios) Go servisine (Port: 8080) bir webhook ateşler.
4. Go servisi, eşzamanlı goroutine'ler (concurrency) barındırarak gelen webhook isteğini API Key onayından ve Rate Limit filtresinden geçirdikten sonra analiz veya bildirim logu olarak tutar. (Bonus B uygulanmıştır)

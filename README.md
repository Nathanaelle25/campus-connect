# CampusConnect
Öğrenci: NATHANAELLE 
Okul No: 2408****150

Üniversite etkinlik platformu — NestJS + Go polyglot backend.

## Teknolojiler
- NestJS (TypeScript) — REST + GraphQL
- Go — Notification + Analytics
- PostgreSQL — Veritabanı

## Kurulum
### Gereksinimler
- Node.js 18+, Go 1.21+, PostgreSQL 15+

### Projenin İndirilmesi
```bash
git clone <repo-url>
cd campus-connect
```

### Veritabanının Çalıştırılması
Projeyle birlikte gelen Docker dosyasını kullanarak veritabanını başlatabilirsiniz:
```bash
docker-compose up -d db
```

### NestJS Service
```bash
cd nestjs-service
npm install
cp .env.example .env   # Değerleri kontrol edin ve düzenleyin
npx prisma db push     # Migration yerine db push kullanarak şemayı aktif edebilir veya migration atabilirsiniz
npm run start:dev
```

### Go Service
```bash
cd go-service
go mod download
cp .env.example .env   # Değerleri kontrol edin ve düzenleyin
go run main.go
```

## Environment Variables
- `DATABASE_URL`: Ortak PostgreSQL bağlantı dizesi (.env.example içinde yer alır)
- `JWT_SECRET`: NestJS üzerinde Auth işlemleri sırasında tokenları imzalamak için kullanılır
- `GO_WEBHOOK_URL`: NestJS'ten Go webhook receiver'ına HTTP POST atmak için URL
- `API_KEY`: Go Servisi API rotaları için kimlik doğrulama anahtarı

## API Endpoints

### NestJS Service (REST & GraphQL) :3000
| Metot | Path | Açıklama |
|-------|------|----------|
| POST  | `/api/v1/users/register` | Yeni kullanıcı oluştur (register) |
| POST  | `/api/v1/users/login` | Kullanıcı girişi (login), JWT token döner |
| GET   | `/api/v1/users/me` | Giriş yapmış kullanıcının profil kontrolü |
| POST  | `/api/v1/events` | Etkinlik oluşturur (Webhook ile Go servisi tetiklenir) |
| GET   | `/api/v1/events` | Etkinlikleri listeler |
| DELETE| `/api/v1/events/:id` | Etkinliği siler (Sadece Admin yetkisi olanlar) |
| POST  | `/graphql` | GraphQL üzerinden event ve user sorgulamaları/mutasyonları |

### Go Service (REST & Webhook) :8080
| Metot | Path | Açıklama |
|-------|------|----------|
| GET   | `/analytics` | Bildirim ve event istatistiklerini getirir (API Key gerekir) |
| GET   | `/notifications` | Alınan bildirim loglarını listeler (API Key gerekir) |
| POST  | `/webhook` | NestJS tarafından event eylemleri sonucunda tetiklenen webhook listener |

## Örnek Request / Response

### 1. Kullanıcı Kaydı (NestJS)
```bash
curl -X POST http://localhost:3000/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"email":"nathy@mail.com","password":"password123","role":"USER"}'
```
**Başarılı Yanıt (201):**
```json
{"id": 1, "email": "nathy@mail.com", "role": "USER", "createdAt": "..."}
```
**Hata Yanıtı (400 - Duplicate Email):**
```json
{"statusCode": 400, "message": "Email already exists"}
```

### 2. Login (NestJS)
```bash
curl -X POST http://localhost:3000/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"nathy@mail.com","password":"password123"}'
```
**Başarılı Yanıt (200):**
```json
{"access_token": "eyJhbGciOiJIUz...<jwt>"}
```

### 3. Etkinlik Oluşturma (NestJS)
```bash
curl -X POST http://localhost:3000/api/v1/events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt_token>" \
  -d '{"title":"Bahar Şenliği","description":"Kampüs ana meydanında bahar festivali."}'
```
**Başarılı Yanıt (201):**
```json
{"id": 1, "title": "Bahar Şenliği"}
```

### 4. Bildirimleri Alma (Go Service)
```bash
curl -X GET http://localhost:8080/notifications \
  -H "X-API-Key: mysecretapikey"
```
**Başarılı Yanıt (200):**
```json
[{"id": 1, "message": "Etkinlik başarıyla oluşturuldu.", "date": "2026-04-06T14:15:00Z"}]
```
**Hata Yanıtı (401 - Unauthorized):**
```json
{"error": "Unauthorized"}
```

### 5. Analytics Alma (Go Service)
```bash
curl -X GET http://localhost:8080/analytics \
  -H "X-API-Key: mysecretapikey"
```
**Başarılı Yanıt (200):**
```json
{"total_events": 15, "total_webhooks_received": 15}
```

## Mimari Kararlar
- **NestJS + Go Polyglot Yaklaşımı:** İş geliştirme hızı, sağlam modüler ekosistemi ve GraphQL ile esnekliği nedeniyle ana API ve CRUD operasyonlarında NestJS kullanıldı. Yüksek performanslı asenkron I/O gerektiren webhook receiver operasyonları için ise goroutines esnekliğinden tam faydalanmak adına Go tercih edildi.
- **Servisler Arası İletişim (Webhook):** Uygulamalar arası güçlü bağlılıkları kırmak (loose coupling) için HTTP Webhook protokolü seçildi. Böylece NestJS doğrudan veriye kilitlenmeden asenkron iletişim sağlar.
- **Go Concurrency (Goroutine):** Go webhook receiver'ına bir istek geldiğinde; DB yazma arayüzü, loglama operasyonları ve analitik sayaçların artırılması işlemleri Go Concurrency modeli (`WaitGroup` vb.) eşliğinde goroutine ile paralel yapılarak maksimum throughput amaçlandı.
- **Güvenlik & Rate Limiting (Go Servisi):** Dışarıdan gelebilecek bruteforce / DDOS tipli saldırıları minimize etmek ve kaynak yükünü kısıtlamak amacıyla memory üzerinde çalışan bir Token Bucket limit mantığı ve gizli `X-API-Key` kimlik doğrulaması eklendi.

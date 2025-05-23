## Proje Hakkında
Stok kontrol servisidir. 

## Temel Özellikler
- Ürün sku'sunu ve bedenini alıp depoda stok var mı kontrol ediyor
- İlgili deponun yeteri kadar adedi var mı, önceliği ne vb, kontrolleri yapıyor

# Kurulum
## Gereksinimler
- PHP 8.1 veya üstü
- Laravel
- Composer
- Docker ve Docker Compose (isteğe bağlı)
- MySQL
- Go (stok kontrol servisi için)
- Ajax
- Ngnix
- [github](https://github.com/yCalbaz/Laravel_e-ticaret_sitesi)

## Proje Yapısı
- app/: Uygulamanın tüm iş mantığını barındırır.
- resources/views/: Blade şablon dosyalarını içerir.
- routes/: Laravel yönlendirme dosyaları.
- tests/: Unit testlerinin bulunduğu yer.
- go/: Go ile yazılmış stok kontrol servisi.(Ayrı bir proje dosyası olarak)

## Kullanılan Teknolojiler
- Laravel: PHP framework
- Go: Stok kontrol servisi
- MySQL: Veritabanı yönetimi
- Docker: Proje konteynerlemesi ve yönetimi
- Vite: JavaScript derleyicisi
- Html / Css : Sayfa görünüşü için
- Ajax: Uyarıları vermek ve sayfayı yenilemek için

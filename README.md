# secondHand GO-Lang

Repository ini berisi source code API Second Hand yang dibangun menggunakan bahasa pemrograman Go-lang.

SecondHand Merupakan website e-commerce tempat berjual beli barang bekas, Dimana user dapat menjual barang bekas mereka dan membuat penawaran untuk barang bekas milik user lain. Jika penjual menerima tawaran dari pembeli, maka mereka dapat melanjutkan transaksi menggunakan WhatsApp. Project ini merupakan pengembangan dari project SecondHand Sebelumnya yang merupakan final project dari program bootcamp di Binar Academy, dimana pada project ini terdapat fitur baru yaitu user role permission. Pada project ini saya mempelajari pengembangan menggunakan repository pattern dan bagaimana cara melakukan dependecy injection pada bahasa pemrograman Go-lang

Pada API SecondHand terdapat dua role yaitu User dan Admin, dimana admin mempunyai privilege untuk menghapus data produk dan transaksi. selain itu juga Admin dapat menambahkan dan menghapus data kota dan data kategeori.

## Teknologi Yang Digunakan
![Teknologi Yang Digunakan](./Backend-Stack.png)

## ERD SECOND HAND
![ERD SECOND HAND](./secondHand-ERD.jpg)

## Cara Menjalankan Aplikasi
### Inisialisasi Database Pada Docker
```bash
# Aktifkan PostgreSQL pada Docker
docker run --name postgres15 -p 5433:5432 -e POSTGRES_USER={username} -e POSTGRES_PASSWORD={password} -d postgres:15-alpine

# Membuat Tabel SecondHand
docker exec -it postgres15 createdb --username={user} --owner={user} secondhand

# Melakukan Migrasi Database
migrate -path db/migration -database "postgresql://{username}:{password}@localhost:5432/secondhand?sslmode=disable" -verbose up
```

### Mengaktifkan RESTful Web API
```bash
go run .
```


## Endpoint API

API Second Hand berisi beberapa endpoint yang dibagi dalam beberapa kelompok yaitu :

### User 
- Register  
Untuk melakukan Registrasi User

- Register Admin
Untuk melakukan Registrasi Admin

- Login  
Untuk login user mengirimkan data berupa email dan password.  
Jika login berhasil user akan menerima JWT yang disimpan pada cookie.

- Logout  
Logout digunakan untuk menghapus JWT pada cookie.  

### Profile 
Terdapat tiga fitur utama dalam endpoint profile yaitu:  
- Get Profile  
Mengambil profile pengguna 

- Update Profile  
Mengubah data pengguna seperti menambahkan data kota, alamat, nomor telepon, dan nama pengguna.  
Agar pengguna dapat menjual barang mereka dan membuat penawaran, pengguna harus melengkapi profile terlebih dahulu.

- Update Image
Mengubah foto profile pengguna yang sedang login.

### Data 
- Create Kota dan Kategori  
Endpoint ini digunakan untuk menambahkan data kota dan kategori dan hanya bisa diakses oleh admin.  

- Get data Kota dan Kategori  
Endpoint ini bersifat publik untuk mengambil data kota dan kategori.

- Delete Kota dan Kategori  
Endpoint ini digunakan untuk menghapus data kota dan kategori dan hanya bisa diakses oleh admin.

### Product
- Create Product  
Menambahkan data produk, secara default data produk yang baru dibuat mempunyai status published = false. Dapat diakses ketika user sudah melengkapi profile. 

- Upload Product image
Menambahkan data gambar produk 

- Get All Product  
Menampilkan data produk yang dipublish dan belum terjual dan dapat diakes secara publik

- Get by Id Product  
Menampilkan data suatu produk secara detail

- Get by Id Product  
Menampilkan data suatu produk berdasarkan kategori dan dapat diakses secara publik

- Get Product By Account
Menampilkan data produk yang dipublish dan belum terjual dari suatu user

- Get my product
Menampilkan data produk user yang sedang login. Menerima query parameter seperti :  
    - published (true/false), untuk memfilter data produk yang dipublish dan belum dipublish
    - sold (true/false), memfilter data produk yang sudah terjual dan belum terjual

- Update Product  
Mengubah data product berdasarkan Id produk

- Update Thumbnail
Mengubah thumbnail produk dari gambar yang sudah di upload.

- Publish Product
Endpoint ini digunakan untuk mem-publish produk agar dapat diakses secara publik dan dapat diakses jika produk sudah mempunyai minimal 1 gambar 

- Set Status Product 
Mengubah status sold product menjadi true atau false

- Delete Product  
Menghapus data suatu produk berdasarkan id produk dan dapat diakses oleh pemilik product dan admin

- Delete Product Image
Menghapus gambar dari produk berdasarkan id produk dan id gambar dan dapat diakses oleh pemilik product dan admin

### Transaction
- Create Transacstion  
Membuat penawaran terhadap suatu produk, dapat diakses ketika user sudah melengkapi profile.  

- Get Transaction Detail By Id  
Mengambil data transaksi berdasarkan Id, hanya dapat diakses oleh user yang terlibat dalam transaksi dan admin.

- Get All Offer By Product Id  
Mengambil data transaki berdasarkan Id produk, hanya bisa diakses oleh pemilik produk dan admin. 

- Get Transaction By Buyer Id  
Mengambil data penawaran yang masuk kepada pengguna yang sedang login berdasarkan User Id yang melakukan penawaran.
  
- Get All offer To Current User  
Mengambil semua data penawaran yang masuk pada user yang sedang login. 
   
- Get My Transaction  
Mengambil semua data penawaran yang dibuat oleh user. 

- Update Offer Price  
Mengupdate data harga penawaran yang sudah dibuat. Hanya dapat diakses oleh penawar.

- Accept Transaction  
Mengubah status accepted menjadi true, hanya dapat diakses oleh pemilik produk.


- Delete Transaction  
Menghapus data transaksi. Hanya dapat diakses oleh pemilik produk dan admin.

## Dokumentasi Menggunakan Postman
Dokumentasi API dapat diakses pada :
https://documenter.getpostman.com/view/17275912/2s93eU1Z2f


# second-hand-golang

Repository ini berisi source code API Second Hand yang dibangun menggunakan bahasa pemrograman golang.

SecondHand Merupakan website e-commerce tempat berjual beli barang bekas, Dimana user dapat menjual barang bekas mereka dan membuat penawaran untuk barang bekas milik user lain. Jika penjual menerima tawaran dari pembeli, maka mereka dapat melanjutkan transaksi menggunakan WhatsApp. Project ini merupakan pengembangan dari project SecondHand Sebelumnya, dimana pada project ini terdapat fitur baru yaitu user role permission.

Pada API SecondHand terdapat dua role yaitu User dan Admin, dimana admin mempunyai privilege untuk mengubah dan menghapus semua data produk, serta mengubah dan menghapus data transaksi. selain itu juga Admin dapat menambahkan dan menghapus data kota dan data kategeori.

## Endpoint API

API Second Hand berisi beberapa endpoint yang dibagi dalam beberapa kelompok yaitu :

### User 
- Register  
Untuk melakukan Registrasi User mengirimkan nama, email, dan password.

- Login  
Untuk login user mengirimkan data berupa email dan password.  
Jika login berhasil user akan menerima JWT yang disimpan pada cookie.

- Logout  
logout digunakan untuk menghapus JWT pada cookie.  

### Profile 
Terdapat dua fitur utama dalam endpoint profile yaitu:  
- Get Profile  
mengambil informasi data pengguna yang sedang login, dan mengambil data pengguna berdasarkan id pengguna.  
- Update Profile  
Mengubah data pengguna seperti menambahkan data kota, alamat, nomor telepon, profile picture, dan nama pengguna.  
Agar pengguna dapat menjual barang mereka dan membuat penawaran, pengguna harus melengkapi profile terlebih dahulu.

### Data 
- Create Kota dan Kategori  
endpoint ini digunakan untuk menambahkan data kota dan kategori dan hanya bisa diakses oleh admin.  

- Get data Kota dan Kategori  
endpoint ini bersifat publik untuk mengambil data kota dan kategori.

- Delete Kota dan Kategori  
endpoint ini digunakan untuk menghapus data kota dan kategori dan hanya bisa diakses oleh admin.

### Product
- Get All Product  
Menampilkan data product yang dipublish dan belum terjual. Dapat diakes secara publik dan dapat menerima beberapa query parameter seperti :  
    - category, untuk memfilter data produk berdasarkan kategori
    - row, untuk membatasi jumlah data per halaman, jika tidak diisi secara default akan menampilkan 50 data.
    - page, untuk mengakses halaman data

- Get by Id Product  
Menampilkan data suatu product secara detail

- Get Product By User Id

- Get Product By Current Logged In User  
Menampilkan data product user yang sedang login. Menerima query parameter seperti :  
    - published (true/false), untuk memfilter data produk yang dipublish dan belum dipublish
    - sold (true/false), memfilter data produk yang sudah terjual dan belum terjual
    - row, untuk membatasi jumlah data per halaman, jika tidak diisi secara default akan menampilkan 50 data.
    - page, untuk mengakses halaman data

- Create Product  
Menambahkan data product, secara default data product yang baru dibuat mempunyai status published = false. Dapat diakses ketika user sudah melengkapi profile. 

- Update Product  
Dapat diakses oleh pemilik product dan admin.

- Delete Product  
Dapat diakses oleh pemilik product dan admin.

### Transaction
- Create Transacstion  
Membuat penawaran terhadap suatu produk, dapat diakses ketika user sudah melengkapi profile.  

- Get Transaction By Id  
Mengambil data transaksi berdasarkan Id, hanya dapat diakses oleh user yang terlibat dalam transaksi dan admin.

- Get Transaction By Product Id  
Mengambil data transaki berdasarkan Id Product, hanya bisa diakses oleh pemilik produk dan admin. Menerima query parameter seperti :  
    - accepted (true/false), untuk memfilter data transaksi yang sudah diterima atau belum diterima.
    - sold (true/false), memfilter transaksi yang sudah terjual dan belum terjual
    - row, untuk membatasi jumlah data per halaman, jika tidak diisi secara default akan menampilkan 50 data.
    - page, untuk mengakses halaman data.  

- Get Transaction By User Id  
Mengambil data transaki berdasarkan User Id terhadap semua produk yang dimiliki user yang sedang login. Menerima query parameter seperti :  
    - accepted (true/false), untuk memfilter data transaksi yang sudah diterima atau belum diterima.
    - sold (true/false), memfilter transaksi yang sudah terjual dan belum terjual
    - row, untuk membatasi jumlah data per halaman, jika tidak diisi secara default akan menampilkan 50 data.
    - page, untuk mengakses halaman data.  

- Get All offer To Current User  
Mengambil semua data penawaran yang masuk. Menerima query parameter seperti :  
    - accepted (true/false), untuk memfilter data transaksi yang sudah diterima atau belum diterima.
    - sold (true/false), memfilter transaksi yang sudah terjual dan belum terjual
    - row, untuk membatasi jumlah data per halaman, jika tidak diisi secara default akan menampilkan 50 data.
    - page, untuk mengakses halaman data.    

- Get Transaction  
Mengambil semua data transaksi yang dilakukan oleh user. Menerima query parameter seperti :  
    - accepted (true/false), untuk memfilter data transaksi yang sudah diterima atau belum diterima.
    - sold (true/false), memfilter transaksi yang sudah terjual dan belum terjual
    - row, untuk membatasi jumlah data per halaman, jika tidak diisi secara default akan menampilkan 50 data.
    - page, untuk mengakses halaman data.   

- Update Offer Price  
Mengupdate data harga penawaran yang sudah dibuat. Hanya dapat diakses oleh penawar dan admin.

- Accept Transaction  
Mengubah status accepted menjadi true, hanya dapat diakses oleh pemilik product dan admin.

- Set Sold Status  
Mengubah status Sold menjadi true, dan menghapus semua data penawaran terhadap produk yang sudah terjual. hanya dapat diakses oleh pemilik produk dan admin.

- Delete Transaction  
Menghapus data transaksi. Hanya dapat diakses oleh penjual dan admin.

### Image
- Upload Product Image  
Menambahkan gambar terhadap suatu product. Hanya dapat diakses pemilik produk dan admin.

- Delete Product Image
Menghapush salah satu gambar product. Hanya dapat diakses pemilik produk dan admin

## Dokumentasi Menggunakan Postman
Dokumentasi API dapat diakses pada :
https://documenter.getpostman.com/view/17275912/2s8ZDcxK4x

